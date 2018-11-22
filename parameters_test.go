package yagclif

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type validStruct struct {
	A int
	B string `yagclif:"shortname:sb;description:foo;default:3"`
	C bool   `yagclif:"shortname:sc"`
}

var validStructType = reflect.TypeOf(validStruct{})

func TestNewParameters(t *testing.T) {
	type faultyStruct struct {
		a int
		b string `yagclif:"something"`
	}
	var faultyStructType = reflect.TypeOf(faultyStruct{})
	t.Run("returns value", func(t *testing.T) {
		params, err := newParameters(validStructType)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(params))
		assert.Equal(t, "A", params[0].name)
		assert.Equal(t, 0, params[0].index)
		assert.Equal(t, "B", params[1].name)
		assert.Equal(t, "foo", params[1].description)
		assert.Equal(t, "3", params[1].defaultValue)
		assert.False(t, params[1].mandatory)
		assert.Equal(t, 1, params[1].index)
		assert.Equal(t, "C", params[2].name)
		assert.False(t, params[2].mandatory)
		assert.Equal(t, 2, params[2].index)
	})
	t.Run("returns error", func(t *testing.T) {
		params, err := newParameters(faultyStructType)
		assert.NotNil(t, err)
		assert.Nil(t, params)
	})
}
func TestCheckValidty(t *testing.T) {
	t.Run("no duplicates", func(t *testing.T) {
		params, err := newParameters(validStructType)
		assert.Nil(t, err)
		assert.Nil(t, params.checkValidity())
	})
	t.Run("duplicates", func(t *testing.T) {
		params := parameters{
			&parameter{
				name:      "A",
				shortName: "B",
			},
			&parameter{
				name:      "B",
				shortName: "B",
			},
		}
		assert.NotNil(t, params.checkValidity())
	})
}

func TestFind(t *testing.T) {
	params, err := newParameters(validStructType)
	assert.Nil(t, err)
	t.Run("finds", func(t *testing.T) {
		param := params.find("--a")
		assert.NotNil(t, param)
	})
	t.Run("does not find", func(t *testing.T) {
		param := params.find("-a")
		assert.Nil(t, param)
	})
}

func TestParseArguments(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		params, err := newParameters(validStructType)
		assert.Nil(t, err)
		testStruct := reflect.New(validStructType).Interface()
		remaining, err := params.ParseArguments(testStruct, []string{"hello", "--b", "world", "-sc", "!"})
		assert.Nil(t, err)
		assert.Equal(t, []string{"hello", "!"}, remaining)
		assert.Equal(t, &validStruct{
			B: "world",
			C: true,
		}, testStruct)
	})
	t.Run("error at setter callback generating", func(t *testing.T) {
		type faultyStruct struct {
			a int
			b interface{}
		}
		faultyStructType := reflect.TypeOf(faultyStruct{})
		testStruct := &faultyStruct{}
		params, err := newParameters(faultyStructType)
		assert.Nil(t, err)
		remaining, err := params.ParseArguments(testStruct, []string{"--b", "hello"})
		assert.Nil(t, remaining)
		assert.NotNil(t, err)
	})
	t.Run("error at executing callback", func(t *testing.T) {
		type foo struct {
			bar int
		}
		fooType := reflect.TypeOf(foo{})
		testStruct := &foo{}
		params, err := newParameters(fooType)
		assert.Nil(t, err)
		remaining, err := params.ParseArguments(testStruct, []string{"--bar", "notanumber"})
		assert.Nil(t, remaining)
		assert.NotNil(t, err)
	})
}

func TestParamsGetHelp(t *testing.T) {
	params, err := newParameters(validStructType)
	assert.Nil(t, err)
	help := params.getHelp()
	assert.Len(t, help, 3)
}

func TestParse(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		testStruct := &validStruct{}
		os.Args = []string{"hello", "--b", "world", "!"}
		remaining, err := Parse(testStruct)
		assert.Nil(t, err)
		assert.Equal(t, []string{"hello", "!"}, remaining)
	})
	t.Run("return err", func(t *testing.T) {
		type faultyStruct struct {
			a int
			b string `yagclif:"something"`
		}
		testStruct := &faultyStruct{}
		os.Args = []string{"hello", "-b", "world", "!"}
		remaining, err := Parse(testStruct)
		assert.Nil(t, remaining)
		assert.NotNil(t, err)
	})
}
