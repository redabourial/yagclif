package yagclif

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/potatomasterrace/catch"
)

type parameters []*parameter

func isSupportedType(sf reflect.StructField) bool {
	supportedTypes := []reflect.Type{
		reflect.TypeOf(true),
		reflect.TypeOf(1), reflect.TypeOf(""),
		reflect.TypeOf([]string{}),
		reflect.TypeOf([]int{}),
	}
	for _, supportedType := range supportedTypes {
		if supportedType == sf.Type {
			return true
		}
	}
	return false
}

// Returns the parameters from an object tags.
func newParameters(tipe reflect.Type) (parameters, error) {
	params := parameters{}
	err := catch.Error(func() {
		tipe.NumField()
	})
	if err != nil {
		return nil, fmt.Errorf("%s\r\ncan not read fields of object \r\n hint: check that you're using a struct type ", err)
	}
	for i := 0; i < tipe.NumField(); i++ {
		field := tipe.Field(i)
		param, err := newParameter(field)
		if err != nil {
			return nil, err
		}
		if param != nil && isSupportedType(field) {
			params = append(params, param)
		} else if field.Tag.Get(tagName) != "omit" {
			inheritedParams, err := newParameters(field.Type)
			if err != nil {
				return nil, fmt.Errorf("%s\r\n error parsing recursively field %s  ", err, field.Name)
			}
			params = append(params, inheritedParams...)
		}
	}
	if err = params.checkValidity(); err != nil {
		return nil, err
	}
	return params, nil
}

// Validates that no conflict exists between parameter names.
// and that every array parameter has a delimiter
func (params *parameters) checkValidity() error {
	existingNames := make(map[string]*parameter, 0)
	for _, param := range *params {
		for _, name := range param.CliNames() {
			conflictingParam := existingNames[name]
			if conflictingParam != nil {
				return fmt.Errorf(
					"conflict for cli name %s struct fields %s and %s",
					name, param.name, conflictingParam.name,
				)
			}
			existingNames[name] = param
		}
	}
	return nil
}

// Finds a parameter in the array by cli names :
// -name or --shortname.
func (params *parameters) find(s string) *parameter {
	for _, param := range *params {
		if param.Matches(s) {
			return param
		}
	}
	return nil
}

// Returns an array describing the parameters.
func (params *parameters) getHelp() []string {
	var buffer []string
	for _, param := range *params {
		buffer = append(buffer, param.GetHelp())
	}
	return buffer
}

func (params *parameters) assignDefaults(obj interface{}) error {
	for _, param := range *params {
		_, err := param.setDefault(obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func (params *parameters) checkForMissingMandatory() error {
	for _, param := range *params {
		if param.mandatory && !param.used {
			if param.description != "" {
				return fmt.Errorf("missing argument %s for %s %s", param.CliNames(), param.name, param.description)
			}
			return fmt.Errorf("missing argument %s for %s", param.CliNames(), param.name)
		}
	}
	return nil
}

// Fills the object with the argument.
// This function only works if the obj
// value is not nil.
func (params *parameters) ParseArguments(obj interface{}, args []string) ([]string, error) {
	params.assignDefaults(obj)
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
	if err := params.checkForMissingMandatory(); err != nil {
		return nil, err
	}
	return remainingArgs, nil
}

func Parse(obj interface{}) (remainingArgs []string, err error) {
	tipe := reflect.TypeOf(obj).Elem()
	params, err := newParameters(tipe)
	if err != nil {
		return nil, err
	}
	remainingArgs, err = params.ParseArguments(obj, os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf(
			"%s\r\nusage:\r\n%s\r\n",
			err, strings.Join(
				params.getHelp(),
				"\r\n",
			),
		)
	}
	return remainingArgs, nil
}

func GetHelp(obj interface{}) string {
	tipe := reflect.TypeOf(obj).Elem()
	params, err := newParameters(tipe)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return strings.Join(
		params.getHelp(),
		"\r\n",
	)

}
