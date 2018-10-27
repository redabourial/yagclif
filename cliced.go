package cliced

func indent(s string) string {
	// TODO implement
	return s
}

type App []action

func (app *App) AddRoute(name string, description string, callback interface{}) {
	return nil
}

func (app *App) Run() {

}

func (app *App) getHelp() string {
	return ""
}

func NewCliApp(name string, description string) App {
	var app App = make([]action, 0)
	return app
}

func Resume([]string) App {
	return nil
}
