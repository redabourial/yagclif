package cliced

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type validStruct struct {
	A int
	B string `cliced:"mandatory;shortname:sb;default:3"`
	C bool   `cliced:"mandatory;shortname:sc"`
}

func TestNewParameters(t *testing.T) {
	type faultyStruct struct {
		a int
		b string `cliced:"something"`
	}
	t.Run("returns value", func(t *testing.T) {
		params, err := newParameters(&validStruct{})
		assert.Nil(t, err)
		assert.Equal(t, 3, len(params))
		assert.Equal(t, "A", params[0].Name())
		assert.Equal(t, 0, params[0].Index())
		assert.Equal(t, "B", params[1].Name())
		assert.True(t, params[1].Mandatory())
		assert.Equal(t, 1, params[1].Index())
		assert.Equal(t, "C", params[2].Name())
		assert.True(t, params[2].Mandatory())
		assert.Equal(t, 2, params[2].Index())
	})
	t.Run("returns error", func(t *testing.T) {
		params, err := newParameters(faultyStruct{})
		assert.NotNil(t, err)
		assert.Nil(t, params)
	})
}
func TestCheckValidty(t *testing.T) {
	t.Run("no duplicates", func(t *testing.T) {
		params, err := newParameters(&validStruct{})
		assert.Nil(t, err)
		assert.Nil(t, params.checkValidity())
	})
	t.Run("duplicates", func(t *testing.T) {
		assert.NotNil(t, parameters{
			parameter{
				name:      "A",
				shortName: "B",
			},
			parameter{
				name:      "B",
				shortName: "B",
			},
		}.checkValidity())
	})
}

func TestFind(t *testing.T) {
	params, err := newParameters(&validStruct{})
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
		testStruct := &validStruct{}
		params, err := newParameters(testStruct)
		assert.Nil(t, err)
		remaining, err := params.ParseArguments(testStruct, []string{"main", "hello", "--b", "world", "!"})
		assert.Nil(t, err)
		assert.Equal(t, []string{"hello", "!"}, remaining)
	})
	t.Run("error at setter callback generating", func(t *testing.T) {
		type faultyStruct struct {
			a int
			b interface{}
		}
		testStruct := &faultyStruct{}
		params, err := newParameters(testStruct)
		assert.Nil(t, err)
		remaining, err := params.ParseArguments(testStruct, []string{"main", "--b", "hello"})
		assert.Nil(t, remaining)
		assert.NotNil(t, err)
	})
	t.Run("error at executing callback", func(t *testing.T) {
		type foo struct {
			bar int
		}
		testStruct := &foo{}
		params, err := newParameters(testStruct)
		assert.Nil(t, err)
		remaining, err := params.ParseArguments(testStruct, []string{"main", "--bar", "notanumber"})
		assert.Nil(t, remaining)
		assert.NotNil(t, err)
	})
}

func TestGetHelp(t *testing.T) {
	testStruct := &validStruct{}
	params, err := newParameters(testStruct)
	assert.Nil(t, err)
	help := params.getHelp()
	assert.Len(t, help, 3)
	fmt.Println(help)
}

func TestParse(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		testStruct := &validStruct{}
		os.Args = []string{"main", "hello", "--b", "world", "!"}
		remaining, err := Parse(testStruct)
		assert.Equal(t, []string{"hello", "!"}, remaining)
		assert.Nil(t, err)
	})
	t.Run("return err", func(t *testing.T) {
		type faultyStruct struct {
			a int
			b string `cliced:"something"`
		}
		testStruct := &faultyStruct{}
		os.Args = []string{"main", "hello", "-b", "world", "!"}
		remaining, err := Parse(testStruct)
		assert.Nil(t, remaining)
		assert.NotNil(t, err)
	})
}
