package cliced

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
)

type parameters []parameter

// Returns the parameters from an object tags.
func newParameters(obj interface{}) (parameters, error) {
	params := parameters{}
	tipe := reflect.Indirect(reflect.ValueOf(obj)).Type()
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
					name, param.Name(), conflictingParam.Name(),
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
func (params parameters) getHelp() string {
	var buffer bytes.Buffer
	for _, param := range params {
		buffer.WriteString(param.GetHelp())
	}
	return buffer.String()
}

// Fills the object with the argument.
// This function only works if the obj
// value is not nil.
func (params parameters) ParseArguments(obj interface{}, args []string) ([]string, error) {
	// TODO add multiple usages removers
	remainingArgs := []string{}
	var callback func(string) error
	for i := 1; i < len(args); i++ {
		arg := args[i]
		param := params.find(arg)
		if callback == nil {
			if param != nil {
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

func Parse(obj interface{}) ([]string, error) {
	return ParseArguments(obj, os.Args)
}
