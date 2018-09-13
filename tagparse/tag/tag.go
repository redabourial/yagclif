package tag

import (
	"errors"
	"fmt"
	"reflect"

	"../../arguments"
)

type Tag struct {
	argumentName string
	fieldIndex   int
	tipe         reflect.Type
}

func (tag Tag) GetFromArgumentTag(args arguments.Args, receiver interface{}) error {
	argumentName, fieldIndex := tag.argumentName, tag.fieldIndex
	switch tag.tipe {
	case reflect.TypeOf(0):
		parameter, err := args.ParseIntParameter(argumentName)
		if err != nil {
			return err
		}
		reflect.ValueOf(receiver).Field(fieldIndex).SetInt(int64(parameter))
	case reflect.TypeOf(""):
		parameter, err := args.ParseStringParameter(argumentName)
		if err != nil {
			return err
		}
		reflect.ValueOf(receiver).Field(fieldIndex).SetString(parameter)
	case reflect.TypeOf(true):
		err := args.UseArgument(argumentName)
		if err != nil {
			return err
		}
		reflect.ValueOf(receiver).Field(fieldIndex).SetBool(true)
	}
	errorMsg := fmt.Sprint("incompatible type for ", tag)
	return errors.New(errorMsg)
}
