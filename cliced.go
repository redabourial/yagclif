package cliced

import (
	"fmt"
	"reflect"
)

type Parameters []Parameter

// Returns the parameters from an object tags.
func newParameters(obj interface{}) (Parameters, error) {
	params := Parameters{}
	tipe := reflect.TypeOf(obj)
	for i := 0; i < tipe.NumField(); i++ {
		field := tipe.Field(i)
		param, err := newParameter(field)
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return params, nil
}

// Validates that no conflict exists between parameter names.
// and that every array parameter has a delimiter
func (params Parameters) checkValidity() error {
	existingNames := *new(map[string]Parameter)
	for _, param := range params {
		for _, name := range param.CliNames() {
			conflictingParam := existingNames[name]
			if conflictingParam != nil {
				return fmt.Errorf(
					"conflict for cli name %s struct fields %s and %s",
					name, param.Name(), conflictingParam.Name(),
				)
			}
			existingNames[name] = param
		}
	}
	return nil
}

// Finds a parameter in the array by cli names :
// -name or --shortname.
func (params Parameters) find(s string) Parameter {
	for _, param := range params {
		if param.Matches(s) {
			return param
		}
	}
	return nil
}

// Fills the object with the argument.
// This function only works if the obj
// value is not nil.
func Parse(obj interface{}, args []string) ([]string, error) {
	remainingArgs := []string{}
	params, err := newParameters(obj)
	if err != nil {
		return nil, err
	}
	var callback func(string) error
	for i := 0; i < len(args); i++ {
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
	return []string{}, nil
}
