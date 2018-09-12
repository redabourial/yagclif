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

func Test_forEachStructField(t *testing.T) {
	actualParameters := make([]string, 0)
	callback := func(field reflect.StructField) {
		fieldName, fieldTag, fieldType := field.Name, field.Tag.Get(TagName), field.Type
		sFieldName, sFieldType, sFieldTag := string(fieldName), fmt.Sprint(fieldType), string(fieldTag)
		actualParameters = append(actualParameters, sFieldName, sFieldType, sFieldTag)
	}
	forEachStructField(reflect.TypeOf(testStruct1{}), callback)
	assert.Equal(t, []string{
		"stringField", "string", "foo",
		"intField", "int", "bar",
	}, actualParameters)
}

func Test_getArgument(t *testing.T) {
	actualParameters := make([]string, 0)
	callback := func(field reflect.StructField) {
		fieldName, fieldTag, fieldType := field.Name, field.Tag.Get(TagName), field.Type
		sFieldName, sFieldType, sFieldTag := string(fieldName), fmt.Sprint(fieldType), string(fieldTag)
		actualParameters = append(actualParameters, sFieldName, sFieldType, sFieldTag)
	}
	forEachStructField(reflect.TypeOf(testStruct1{}), callback)
	assert.Equal(t, []string{
		"stringField", "string", "foo",
		"intField", "int", "bar",
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
				reflect.TypeOf(""),
			},
			argumentTag{
				"bar",
				reflect.TypeOf(0),
			},
		})
	})
}

func Test_Parse(t *testing.T) {
	// t.Run("works", func(t *testing.T) {
	// 	conf := CliConfiguartion{
	// 		reflect.TypeOf(testStruct1{}),
	// 		nil,
	// 	}
	// 	conf.Init()
	// 	conf.Parse(ar)
	// })
}
