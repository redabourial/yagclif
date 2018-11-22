package yagclif

import (
	"fmt"
	"reflect"

	"github.com/potatomasterrace/catch"
)

type route struct {
	description string
	callback    func(args []string) error
	customType  reflect.Type
}

// Return the type of the custom argument
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

func getCustomCallBack(callBackFunctionValue reflect.Value, callBackCustomType reflect.Type) (func(args []string) error, error) {
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

func formatCallBack(callBackFunctionValue reflect.Value, callBackArgType reflect.Type) (startRoute func(args []string) error, err error) {
	if callBackArgType == nil {
		return getSimpleCallBack(callBackFunctionValue), nil
	}
	return getCustomCallBack(callBackFunctionValue, callBackArgType)
}

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
		description: description,
		callback:    formatedCallback,
		customType:  callBackArgType,
	}, nil
}

func (r *route) run(args []string) error {
	if r.callback != nil {
		return r.callback(args)
	}
	return fmt.Errorf("callback not defined")
}

func (r *route) getHelp() []string {
	if r.customType == nil {
		return []string{}
	}
	parameters, err := newParameters(r.customType)
	if err != nil {
		return []string{"Could not parse parameters"}
	}
	return parameters.getHelp()
}
