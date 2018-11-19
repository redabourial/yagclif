package cliced

import (
	"fmt"
)

func indent(s []string) []string {
	for i, line := range s {
		s[i] = fmt.Sprint("\t", line)
	}
	return s
}

type App struct {
	name        string
	description string
	routes      map[string]*route
}

func (app App) AddRoute(name string, description string, callback interface{}) error {
	if app.routes[name] != nil {
		return fmt.Errorf(
			"route %s already used",
			name,
		)
	}
	route, err := newRoute(description, callback)
	if err != nil {
		app.routes[name] = route
	}
	return err
}

func (app *App) Run() {
	// args := os.Args
}

func (app *App) GetHelp() string {
	return ""
}

func NewCliApp(name string, description string) *App {
	return nil
}
