package cliced

import (
	"fmt"
	"reflect"
)

// route is used to manage route routing.
type Route struct {
	name        string
	description string
	parameters  parameters
	// Panicky callback
	callBack    func([]string)
}

func formatRouteCallBack(callBack interface{}) (func([]string), parameters,error) {
	callBackTipe := reflect.TypeOf(callBack)
	callBackValue := reflect.ValueOf(callBack)
	switch callBackTipe.NumIn() {
	case 1:
		// instance of expectedClass
		expectTypeInstance := func([]string) {}
		if callBackTipe.AssignableTo(reflect.TypeOf(expectTypeInstance)) {
			return func(args []string) {
				in := []reflect.Value{reflect.ValueOf(args)}
				callBackValue.Call(in)
			}, nil
		}
	case 2:
		if callBackTipe.In(0) == reflect.TypeOf([]string{}) {
			firstParamType := callBackTipe.In(0)
			firstParamParamInstance := reflect.New(firstParamType)
			// test if second parameter is []string
			if callBackTipe.In(1) == reflect.TypeOf([]string{}) {
				return func(args []string) {
					remainingArgs, err := ParseArguments(&firstParamParamInstance, args)
					if err != nil {
						panic(err)
					}
					in := []reflect.Value{
						reflect.ValueOf(firstParamParamInstance),
						reflect.ValueOf(remainingArgs),
					}
					callBackValue.Call(in)
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("expected type func([]string) or func(SomeStruct,[]string) but instead found %s",
		callBackTipe)
}

func newAction(name string, description string, callBack interface{}) (*action, error) {
	formatedActionCallBack, err := formatActionCallBack(callBack)
	if err != nil {
		return nil, err
	}
	return &Route{
		name:        name,
		description: description,
		callBack:    formatedActionCallBack,
		
	}, nil
}
func (r Route) getHelp() string {
	return fmt.Printf(
		""
	)
}
func (r Route) execute(args []string) {

}
