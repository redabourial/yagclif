package tagparse

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct1 struct {
	stringField string `tabona:"foo"`
	intField    int    `tabona:"bar"`
}

func Test_New(t *testing.T) {
	actualParameters := make([]string, 0)
	callback := func(field reflect.StructField, index int) {
		fieldName, fieldTag, fieldType := field.Name, field.Tag.Get(TagName), field.Type
		sFieldName, sFieldType, sFieldTag := string(fieldName), fmt.Sprint(fieldType), string(fieldTag)
		actualParameters = append(actualParameters, string(index), sFieldName, sFieldType, sFieldTag)
	}
	forEachStructField(reflect.TypeOf(testStruct1{}), callback)
	assert.Equal(t, []string{
		"\x00", "stringField", "string", "foo",
		"\x01", "intField", "int", "bar",
	}, actualParameters)
}

func Test_getArgument(t *testing.T) {
	actualParameters := make([]string, 0)
	callback := func(field reflect.StructField, index int) {
		fieldName, fieldTag, fieldType := field.Name, field.Tag.Get(TagName), field.Type
		sFieldName, sFieldType, sFieldTag := string(fieldName), fmt.Sprint(fieldType), string(fieldTag)
		actualParameters = append(actualParameters, string(index), sFieldName, sFieldType, sFieldTag)
	}
	forEachStructField(reflect.TypeOf(testStruct1{}), callback)
	assert.Equal(t, []string{
		"\x00", "stringField", "string", "foo",
		"\x01", "intField", "int", "bar",
	}, actualParameters)
}

func Test_Init(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		conf := CliConfiguartion{
			reflect.TypeOf(testStruct1{}),
			nil,
		}
		conf.Init()
		assert.Equal(t, conf.args, []argumentTag{
			argumentTag{
				"foo",
				0,
				reflect.TypeOf(""),
			},
			argumentTag{
				"bar",
				1,
				reflect.TypeOf(0),
			},
		})
	})
}

func Test_getFromArgumentTag(t *testing.T) {
	t.Run("works", func(t *testing.T) {
	})
}
