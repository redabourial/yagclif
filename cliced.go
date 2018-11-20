package yagclif

import (
	"fmt"
	"os"
	"strings"
)

func prependToArray(strs []string, prepend string) string {
	var buffer strings.Builder
	for _, str := range strs {
		buffer.WriteString(prepend)
		buffer.WriteString(str)
		buffer.WriteString("\r\n")
	}
	return buffer.String()
}

type App struct {
	name        string
	description string
	routes      map[string]*route
}

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

func (app *App) Run() error {
	args := os.Args
	getError := func(err interface{}) error {
		help := app.GetHelp()
		return fmt.Errorf("%s\r\n%s", err, help)
	}
	if len(args) < 2 {
		return getError("no route was found")
	}
	routeName := args[1]
	route := app.routes[routeName]
	if route == nil {
		errMsg := fmt.Sprintf("%s route not found", routeName)
		return getError(errMsg)
	}
	err := route.run(args[2:])
	if err != nil {
		return getError(err)
	}
	return nil
}

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

func NewCliApp(name string, description string) *App {
	return &App{
		name:        name,
		description: description,
		routes:      map[string]*route{},
	}
}
