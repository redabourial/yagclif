package tagparse

import (
	"fmt"
	"reflect"

	"../argumentparse"
)

func getArgument(field reflect.StructField) string {
	tag, name := field.Tag.Get("cli"), field.Name
	if tag != "" {
		return "--" + tag
	}
	return "--" + name
}

func tagParse(args []string, conf interface{}) ([]error, *[]string, interface{}) {
	if reflect.ValueOf(conf).Kind() != reflect.Ptr {
		panic("passed conf is not a pointer")
	}
	errors := make([]error, 0)
	confType := reflect.TypeOf(conf)
	newconf := reflect.New(confType)
	for i := 0; i < confType.NumField(); i++ {
		argumentName := getArgument(confType.Field(i))
		switch typeField := confType.Field(i).Type.String(); typeField {
		case "string":
			value, newargs, err := argumentparse.ParseArgumentStringParameter(&args, argumentName)
			args = *newargs
			if err != nil {
				errors = append(errors, err)
			} else {
				newconf.Field(i).SetString(*value)
			}
		case "int":
			value, newargs, err := argumentparse.ParseArgumentIntParameter(&args, argumentName)
			args = *newargs
			if err != nil {
				errors = append(errors, err)
			} else {
				newconf.Field(i).SetInt(int64(value))
			}
		case "bool":
			exists, newargs := argumentparse.Exists(&args, argumentName)
			args = *newargs
			newconf.Field(i).SetBool(exists)
		default:
			panic(fmt.Sprint("type ", confType.Field(i).Type.String(), " is not supported"))
		}
	}
	return errors, &args, newconf
}
