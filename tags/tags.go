package tagparse

import (
	"reflect"

	"../arguments"
	"./tag"
)

type Tags []tag.Tag

func New(targetType reflect.Type) Tags {
	tags := make(Tags, 0)
	for fieldIndex := 0; fieldIndex < targetType.NumField(); fieldIndex++ {
		field := targetType.Field(fieldIndex)
		tag := tag.New(field, fieldIndex)
		tags = append(tags, tag)
	}
	return tags
}

func (tags Tags) GetHelp(helpIntro string, helpAfter string) string {
	return ""
}

func (tags Tags) Parse(args arguments.Arguments, receiver interface{}) ([]error, []string) {
	errs := make([]error, 0)
	for _, tag := range tags {
		tag.GetFromArgumentTag(args, receiver)
	}
	return errs, args.GetUnused()
}
