package yagclif

import (
	"bytes"
	"fmt"
	"os"

	"github.com/potatomasterrace/catch"
)

// concatenates the string array by adding a return to line
// at the end and a string at the beginning of each line.
func prependToArray(strs []string, prepend string) string {
	var buffer bytes.Buffer
	for _, str := range strs {
		buffer.WriteString(prepend)
		buffer.WriteString(str)
		buffer.WriteString("\r\n")
	}
	return buffer.String()
}

// App is an implementation of the cli app.
// It handles routing.
type App struct {
	name        string
	description string
	routes      map[string]*route
}

// AddRoute is the methode for adding routes to the cli app.
func (app *App) AddRoute(name string, description string, callback interface{}) error {
	if app.routes[name] != nil {
		return fmt.Errorf(
			"route %s already used",
			name,
		)
	}
	route, err := newRoute(description, callback)
	if err == nil {
		app.routes[name] = route
	}
	return err
}

// RunNoPanic is the method to start running the cli app.
func (app *App) RunNoPanic(outputHelpOnError bool) error {
	return catch.Error(func() {
		app.Run(outputHelpOnError)
	})
}

// Run is the method to start running the cli app.
func (app *App) RunWithArgs(args []string, outputHelpOnError bool) {
	// formatError formats the error to output it.
	formatError := func(err interface{}) error {
		if outputHelpOnError {
			help := app.GetHelp()
			return fmt.Errorf("%s\r\n%s\r\n", err, help)
		}
		return fmt.Errorf("%s", err)
	}
	// if no argument was supplied.
	if len(args) < 2 {
		err := formatError("no action was selected")
		panic(err)
	}
	routeName := args[1]
	route := app.routes[routeName]
	if route == nil {
		errMsg := fmt.Sprintf("%s action not found", routeName)
		err := formatError(errMsg)
		panic(err)
	}
	err := route.run(args[2:])
	if err != nil {
		errMsg := formatError(err)
		panic(errMsg)
	}
}

// GetHelp return the help for the current cli app.
func (app *App) GetHelp() string {
	var buffer bytes.Buffer
	writeln := func(s string) {
		buffer.WriteString(s)
		buffer.WriteString("\r\n")
	}
	writeln(app.name)
	writeln(app.description)
	writeln("")
	for routeName, route := range app.routes {
		routeTitle := fmt.Sprintf("\t %s : %s", routeName, route.description)
		writeln(routeTitle)
		if route.parameterType != nil {
			writeln("\t\t usage :")
		}
		routeArgsHelp := route.getHelp()
		routeHelp := prependToArray(routeArgsHelp, "\t\t\t")
		writeln(routeHelp)
	}
	return buffer.String()
}

// Run is the method to start running the cli app.
func (app *App) Run(outputHelpOnError bool) {
	app.RunWithArgs(os.Args, outputHelpOnError)
}

// NewCliApp creates a new cli app.
func NewCliApp(name string, description string) *App {
	return &App{
		name:        name,
		description: description,
		routes:      map[string]*route{},
	}
}
