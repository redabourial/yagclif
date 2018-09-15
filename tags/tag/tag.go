package tag

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"../../arguments"
)

const TagName = "tabona"

type Tag struct {
	argumentName string
	fieldIndex   int
	mandatory    bool
	description  string
	tipe         reflect.Type
}

func New(field reflect.StructField, fieldIndex int) Tag {
	// TODO add features
	fieldName, fieldType := field.Name, field.Type
	fieldTag := field.Tag.Get(TagName)
	tagDetails := strings.Split(fieldTag, ",")
	lenTagDetails := len(tagDetails)
	if lenTagDetails > 2 {
		errMsg := fmt.Sprint("Error parsing field", field.Name, "too many arguments")
		panic(errMsg)
	} else if lenTagDetails == 2 {
		fieldDescription, mandatory := tagDetails[0], tagDetails[1]
		if mandatory == "mandatory" {
			return Tag{fieldName, fieldIndex, true, fieldDescription, fieldType}
		}
		errMsg := fmt.Sprint("Error parsing field", field.Name, "unrecognized option", mandatory)
		panic(errMsg)
	} else if lenTagDetails == 1 {
		fieldDescription := tagDetails[0]
		return Tag{fieldName, fieldIndex, false, fieldDescription, fieldType}
	}
	return Tag{fieldName, fieldIndex, false, "", fieldType}
}

func (tag *Tag) GetArgumentHelp() string {
	if tag.mandatory {
		return fmt.Sprint(tag.argumentName, " ", tag.tipe, " (mandatory) : ", tag.description)
	} else if tag.description != "" {
		return fmt.Sprint(tag.argumentName, " ", tag.tipe, " : ", tag.description)
	}
	return fmt.Sprint(tag.argumentName, " ", tag.tipe)
}

func (tag *Tag) GetFromArgumentTag(args arguments.Arguments, receiver interface{}) error {
	argumentName, fieldIndex := tag.argumentName, tag.fieldIndex
	switch tag.tipe {
	// for type int
	case reflect.TypeOf(0):
		parameter, err := args.ParseIntParameter(argumentName)
		if err != nil {
			return err
		}
		reflect.ValueOf(receiver).Field(fieldIndex).SetInt(int64(parameter))
	// for type boolean
	case reflect.TypeOf(true):
		err := args.UseArgument(argumentName)
		if err != nil {
			return err
		}
		reflect.ValueOf(receiver).Field(fieldIndex).SetBool(true)
	// for type string
	case reflect.TypeOf(""):
		parameter, err := args.ParseStringParameter(argumentName)
		if err != nil {
			return err
		}
		reflect.ValueOf(receiver).Field(fieldIndex).SetString(parameter)
	}
	errorMsg := fmt.Sprint("incompatible type for ", tag)
	return errors.New(errorMsg)
}
