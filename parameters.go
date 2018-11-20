package yagclif

import (
	"fmt"
	"os"
	"reflect"

	"github.com/potatomasterrace/catch"
)

type parameters []parameter

// Returns the parameters from an object tags.
func newParameters(tipe reflect.Type) (parameters, error) {
	params := parameters{}
	_, err := catch.Panic(func() {
		tipe.NumField()
	})
	if err != nil {
		return nil, fmt.Errorf("can not read fields of object \r\n hint: check that you're using a struct type ")
	}
	for i := 0; i < tipe.NumField(); i++ {
		field := tipe.Field(i)
		param, err := newParameter(field)
		if err != nil {
			return nil, err
		}
		params = append(params, *param)
	}
	return params, nil
}

// Validates that no conflict exists between parameter names.
// and that every array parameter has a delimiter
func (params parameters) checkValidity() error {
	existingNames := make(map[string]*parameter, 0)
	for _, param := range params {
		for _, name := range param.CliNames() {
			conflictingParam := existingNames[name]
			if conflictingParam != nil {
				return fmt.Errorf(
					"conflict for cli name %s struct fields %s and %s",
					name, param.name, conflictingParam.name,
				)
			}
			existingNames[name] = &param
		}
	}

	return nil
}

// Finds a parameter in the array by cli names :
// -name or --shortname.
func (params parameters) find(s string) *parameter {
	for _, param := range params {
		if param.Matches(s) {
			return &param
		}
	}
	return nil
}

// Returns an array describing the parameters.
func (params parameters) getHelp() []string {
	var buffer []string
	for _, param := range params {
		buffer = append(buffer, param.GetHelp())
	}
	return buffer
}

// Fills the object with the argument.
// This function only works if the obj
// value is not nil.
func (params parameters) ParseArguments(obj interface{}, args []string) ([]string, error) {
	// TODO add multiple usages
	// TODO assign defaults
	remainingArgs := []string{}
	var callback func(string) error
	for _, arg := range args {
		param := params.find(arg)
		if callback == nil {
			if param != nil {
				var err error
				callback, err = param.SetterCallback(obj)
				if err != nil {
					return nil, err
				}
			} else {
				remainingArgs = append(remainingArgs, arg)
			}
		} else {
			err := callback(arg)
			if err != nil {
				return nil, err
			}
			callback = nil
		}
	}
	return remainingArgs, nil
}

func Parse(obj interface{}) (remainingArgs []string, err error) {
	tipe := reflect.TypeOf(obj).Elem()
	params, err := newParameters(tipe)
	if err != nil {
		return nil, err
	}
	return params.ParseArguments(obj, os.Args)
}
