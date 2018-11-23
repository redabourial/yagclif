package yagclif

import (
	"fmt"
	"os"
	"strings"
)

// concatenates the string array by adding a return to line
// at the end and a string at the beginning of each line.
func prependToArray(strs []string, prepend string) string {
	var buffer strings.Builder
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

// Run is the method to start running the cli app.
func (app *App) Run(outputHelpOnError bool) error {
	args := os.Args
	// getError formats the error to output it.
	getError := func(err interface{}) error {
		if outputHelpOnError {
			help := app.GetHelp()
			return fmt.Errorf("%s\r\n%s", err, help)
		}
		return fmt.Errorf("%s", err)
	}
	// if no argument was supplied.
	if len(args) < 2 {
		return getError("no action was selected")
	}
	routeName := args[1]
	route := app.routes[routeName]
	if route == nil {
		errMsg := fmt.Sprintf("%s action not found", routeName)
		return getError(errMsg)
	}
	err := route.run(args[2:])
	if err != nil {
		return getError(err)
	}
	return nil
}

// GetHelp return the help for the current cli app.
func (app *App) GetHelp() string {
	var buffer strings.Builder
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
		routeArgsHelp := route.getHelp()
		routeHelp := prependToArray(routeArgsHelp, "\t\t")
		writeln(routeHelp)
	}
	return buffer.String()
}

// NewCliApp creates a new cli app.
func NewCliApp(name string, description string) *App {
	return &App{
		name:        name,
		description: description,
		routes:      map[string]*route{},
	}
}
