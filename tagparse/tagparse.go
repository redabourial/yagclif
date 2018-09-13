package tagparse

import (
	"reflect"

	"../arguments"
)

const TagName = "tabona"

// Iterates over the type(tipe) fields passing
// (Name,type,tag) of each to callback
func forEachStructField(tipe reflect.Type, callBack func(field reflect.StructField, fieldIndex int)) {
	for i := 0; i < tipe.NumField(); i++ {
		field := tipe.Field(i)
		callBack(field, i)
	}
}

// Object to store your parameters
type CliConfiguartion struct {
	targetType reflect.Type
	args       []argumentTag
}

func (c *CliConfiguartion) getArgument(field reflect.StructField, fieldIndex int) (argumentTag, error) {
	// TODO add features
	fieldName, fieldTag, fieldType := field.Name, field.Tag.Get(TagName), field.Type
	if fieldTag != "" {
		return argumentTag{fieldTag, fieldIndex, fieldType}, nil
	}
	return argumentTag{fieldName, fieldIndex, fieldType}, nil
}

func (c *CliConfiguartion) Init() {
	args := make([]argumentTag, 0)
	if c.args == nil {
		forEachStructField(c.targetType, func(field reflect.StructField, fieldIndex int) {
			arg, err := c.getArgument(field, fieldIndex)
			if err == nil {
				args = append(args, arg)
			}
		})
	}
	c.args = args
}

func (c *CliConfiguartion) Parse(args arguments.Args, receiver interface{}) ([]error, []string) {
	errs := make([]error, 0)
	for _, arg := range c.args {
		errs = append(errs, arg.getFromArgumentTag(args, receiver))
	}
	return errs, args.GetUnused()
}

func (CliConfiguartion) GetHelp(helpIntro string, helpAfter string) string {
	return ""
}
