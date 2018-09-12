package tagparse

import (
	"reflect"
)

const TagName = "tabona"

// Object to store your parameters
type CliConfiguartion struct {
	targetType reflect.Type
	args       []argumentTag
}

type argumentTag struct {
	text string
	tipe reflect.Type
}

// Iterates over the type(tipe) fields passing
// (Name,type,tag) of each to callback
func forEachStructField(tipe reflect.Type, callBack func(field reflect.StructField)) {
	for i := 0; i < tipe.NumField(); i++ {
		field := tipe.Field(i)
		callBack(field)
	}
}

func (c *CliConfiguartion) getArgument(field reflect.StructField) (argumentTag, error) {
	// TODO add features
	fieldName, fieldTag, fieldType := field.Name, field.Tag.Get(TagName), field.Type
	if fieldTag != "" {
		return argumentTag{fieldTag, fieldType}, nil
	}
	return argumentTag{fieldName, fieldType}, nil
}

func (c *CliConfiguartion) Init() {
	args := make([]argumentTag, 0)
	if c.args == nil {
		forEachStructField(c.targetType, func(field reflect.StructField) {
			arg, err := c.getArgument(field)
			if err == nil {
				args = append(args, arg)
			}
		})
	}
	c.args = args
}

// func (c *CliConfiguartion) getFromArg(arg arguments.Arg) {

// }
// func (c *CliConfiguartion) Parse(args arguments.Args, receiver interface{}) ([]error, []string) {
// 	errs := make([]error, 0)
// 	for i, arg := range c.args {
// 		switch arg.tipe {

// 		case reflect.TypeOf(0) :
// 		case reflect.TypeOf()
// 		}
// 	}

// }

func (CliConfiguartion) GetHelp(helpIntro string, helpAfter string) string {
	return ""
}
