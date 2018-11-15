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

type App map[string]*route

func (app App) AddRoute(name string, description string, callback interface{}) error {
	if app[name] != nil {
		return fmt.Errorf(
			"route %s already used",
			name,
		)
	}
	// TODO add route initialization
	return nil
}

func (app *App) Run() {

}

func (app *App) GetHelp() string {
	return ""
}

func NewCliApp(name string, description string) App {

	return nil
}
