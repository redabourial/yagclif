package cliced

import (
	"fmt"
	"reflect"

	"github.com/potatomasterrace/catch"
)

type route struct {
	description string
	parameters  parameters
	// Panicky callback
	formatedCallBack func([]string) error
}

// Return the type of the custom argument
func getCustomCallBackType(callBack interface{}) (reflect.Type, error) {
	if callBack == nil {
		return nil, fmt.Errorf("callback value cannot be nil")
	}
	callBackTipe := reflect.TypeOf(callBack)
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
	return nil, fmt.Errorf("expected type func([]string) or func(SomeStruct,[]string) but instead found %s",
		callBackTipe)
}

func getSimpleCallBack(callBackFunctionValue reflect.Value) func(args []string) error {
	return func(args []string) error {
		_, err := catch.Panic(func() {
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

func getCustomCallBack(firstParamInstance reflect.Value, callBackFunctionValue reflect.Value) (func(args []string) error, error) {
	params, err := newParameters(firstParamInstance)
	if err != nil {
		return nil, err
	}
	return func(args []string) error {
		remainingArgs, err := params.ParseArguments(firstParamInstance, args)
		if err != nil {
			return err
		}
		arguments := make([]reflect.Value, 2)
		arguments[0] = reflect.ValueOf(params)
		arguments[1] = reflect.ValueOf(remainingArgs)
		callBackFunctionValue.Call(arguments)
		return nil
	}, nil
}

func formatCallBack(callBack interface{}) (startRoute func(args []string) error, err error) {
	callBackFunctionValue := reflect.ValueOf(callBack)
	callBackArgType, err := getCustomCallBackType(callBack)
	if err != nil {
		return nil, err
	}
	if callBackArgType == nil {
		return getSimpleCallBack(callBackFunctionValue), nil
	}
	firstParamInstance := reflect.New(callBackArgType)
	return getCustomCallBack(firstParamInstance, callBackFunctionValue)
}

func newRoute(description string, callBack interface{}) (*route, error) {
	formatedCallBack, err := formatCallBack(callBack)
	if err != nil {
		return nil, err
	}
	return &route{
		description:      description,
		formatedCallBack: formatedCallBack,
	}, nil
}

func (r route) getHelp() string {
	return ""
}

func (r route) execute(args []string) {

}
