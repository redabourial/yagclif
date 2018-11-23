package yagclif

import (
	"fmt"
	"reflect"

	"github.com/potatomasterrace/catch"
)

// route is an implementation of a cli route.
type route struct {
	description      string
	formatedCallback func(args []string) error
	parameterType    reflect.Type
}

// Return the type of the custom argument.
// returns nil,nil if type of callback is func(somestruct,[]string).
func getCustomCallBackType(callBack interface{}) (reflect.Type, error) {
	if callBack == nil {
		return nil, fmt.Errorf("callback value cannot be nil")
	}
	callBackTipe := reflect.TypeOf(callBack)
	if callBackTipe.Kind().String() == "func" {
		switch callBackTipe.NumIn() {
		case 1:
			// instance of the expected type []string
			if callBackTipe.AssignableTo(reflect.TypeOf(func([]string) {})) {
				return nil, nil
			}
		case 2:
			if callBackTipe.In(1) == reflect.TypeOf([]string{}) {
				return callBackTipe.In(0), nil
			}
		}
	}
	return nil, fmt.Errorf(
		"expected type func([]string) or func(SomeStruct,[]string) but instead found %s",
		callBackTipe,
	)
}

// getSimpleCallBack returns a function that calls the callbackFunction with remaining arguments.
func getSimpleCallBack(callBackFunctionValue reflect.Value) func(args []string) error {
	return func(args []string) error {
		err := catch.CatchError(func() {
			arguments := make([]reflect.Value, 1)
			arguments[0] = reflect.ValueOf(args)
			callBackFunctionValue.Call(arguments)
		})
		if err != nil {
			return fmt.Errorf("%s", err)
		}
		return nil
	}
}

// getSimpleCallBack returns a function that calls the callbackFunction with an instance
// of its custom parameter and remaining arguments.
func getCustomCallBack(callBackFunctionValue reflect.Value, callBackCustomType reflect.Type) (callback func(args []string) error, err error) {
	params, err := newParameters(callBackCustomType)
	if err != nil {
		return nil, err
	}
	firstParamInstance := reflect.New(callBackCustomType)
	return func(args []string) error {
		remainingArgs, err := params.ParseArguments(firstParamInstance.Interface(), args)
		if err != nil {
			return err
		}
		arguments := make([]reflect.Value, 2)
		arguments[0] = firstParamInstance.Elem()
		arguments[1] = reflect.ValueOf(remainingArgs)
		_, callError := catch.Panic(func() {
			callBackFunctionValue.Call(arguments)
		})
		if callError != nil {
			return fmt.Errorf("%s", callError)
		}
		return nil
	}, nil
}

// formatCallBack formats the callback function into a func(args []string)error that executes the callback with arguments.
func formatCallBack(callBackFunctionValue reflect.Value, callBackArgType reflect.Type) (executeCallback func(args []string) error, err error) {
	if callBackArgType == nil {
		return getSimpleCallBack(callBackFunctionValue), nil
	}
	return getCustomCallBack(callBackFunctionValue, callBackArgType)
}

// newRoute creates a new route.
func newRoute(description string, callBack interface{}) (*route, error) {
	callBackFunctionValue := reflect.ValueOf(callBack)
	callBackArgType, err := getCustomCallBackType(callBack)
	if err != nil {
		return nil, err
	}
	formatedCallback, err := formatCallBack(callBackFunctionValue, callBackArgType)
	if err != nil {
		return nil, err
	}
	return &route{
		description:      description,
		formatedCallback: formatedCallback,
		parameterType:    callBackArgType,
	}, nil
}

// run executes the formated callback with the arguments.
func (r *route) run(args []string) error {
	formatedCallback := r.formatedCallback
	if formatedCallback != nil {
		return formatedCallback(args)
	}
	return fmt.Errorf("callback not defined")
}

// getHelp returns an array string.
// Each element is a line of the help text.
func (r *route) getHelp() []string {
	if r.parameterType == nil {
		return []string{}
	}
	parameters, err := newParameters(r.parameterType)
	if err != nil {
		return []string{"Could not parse parameter type"}
	}
	return parameters.getHelp()
}
