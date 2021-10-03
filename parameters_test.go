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
		t.Run("new parameters error", func(t *testing.T) {
			params, err := newParameters(faultyStructType)
			assert.NotNil(t, err)
			assert.Nil(t, params)
		})
		t.Run("non struct type", func(t *testing.T) {
			params, err := newParameters(reflect.TypeOf(""))
			assert.NotNil(t, err)
			assert.Nil(t, params)
		})
		t.Run("non valid tags", func(t *testing.T) {
			type foo struct {
				field1 bool `yagclif:"shortname:sb"`
				field2 bool `yagclif:"shortname:sb"`
			}
			params, err := newParameters(reflect.TypeOf(foo{}))
			assert.NotNil(t, err)
			assert.Nil(t, params)
		})
		t.Run("non valid recursion", func(t *testing.T) {
			type foo2 struct {
				faultyStruct
			}
			params, err := newParameters(reflect.TypeOf(foo2{}))
			assert.NotNil(t, err)
			assert.Nil(t, params)
		})
	})
}

type inheritanceTestStruct struct {
	validStruct
	D string `yagclif:"shortname:sd;description:foo;default:3"`
	E bool   `yagclif:"shortname:se"`
}

var inheritanceTestStructType = reflect.TypeOf(inheritanceTestStruct{})

func TestNewParametersInheritance(t *testing.T) {
	t.Run("returns value", func(t *testing.T) {
		params, err := newParameters(inheritanceTestStructType)
		assert.Nil(t, err)
		assert.Equal(t, 5, len(params))
		assert.Equal(t, "A", params[0].name)
		assert.Equal(t, 0, params[0].index)
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

func TestParamsGetHelp(t *testing.T) {
	params, err := newParameters(validStructType)
	assert.Nil(t, err)
	help := params.getHelp()
	assert.Len(t, help, 3)
}
func TestAssignDefault(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		type Foo struct {
			Solution      int `yagclif:"default:42"`
			OtherSolution int
			Problem       int `yagclif:"default:52;env:TestAssignDefault_Problem"`
			NoProblem     int `yagclif:"omit"`
			OtherProblem  int `yagclif:"default:51;env:TestAssignDefault_Problem_empty_key"`
		}
		os.Setenv("TestAssignDefault_Problem", "53")
		fooInstance := Foo{}
		params, err := newParameters(reflect.TypeOf(fooInstance))
		assert.Nil(t, err)
		assert.NotNil(t, params)
		err = params.assignDefaults(&fooInstance)
		assert.Nil(t, err)
		assert.Equal(t, Foo{
			Solution:      42,
			Problem:       53,
			NoProblem:     0,
			OtherSolution: 0,
			OtherProblem:  51,
		}, fooInstance)
		os.Setenv("TestAssignDefault_Problem_empty_key", "toto")
		fooInstance = Foo{}
		params, err = newParameters(reflect.TypeOf(fooInstance))
		assert.NotNil(t, err)
		assert.Equal(t, Foo{
			Solution:      0,
			Problem:       0,
			NoProblem:     0,
			OtherSolution: 0,
			OtherProblem:  0,
		}, fooInstance)

		os.Setenv("TestAssignDefault_Problem_empty_key", "toto")

		type BadFoo struct {
			Solution os.File `yagclif:"default:toto"`
		}
		badfooInstance := BadFoo{}
		params, err = newParameters(reflect.TypeOf(badfooInstance))
		assert.NotNil(t, err)
		assert.Equal(t, BadFoo{}, badfooInstance)
	})
	t.Run("returns errors", func(t *testing.T) {
		type Foo struct {
			Solution int
		}
		fooInstance := Foo{}
		params, err := newParameters(reflect.TypeOf(fooInstance))
		assert.Nil(t, err)
		assert.NotNil(t, params)
		// stub an unparselable value
		params[0].defaultValue = "hello"
		err = params.assignDefaults(&fooInstance)
		assert.NotNil(t, err)
	})
}

func TestCheckForMissingMandatory(t *testing.T) {
	params := parameters{
		&parameter{
			used:      true,
			mandatory: true,
		},
		&parameter{
			used:      true,
			mandatory: true,
		},
	}
	t.Run("false negative", func(t *testing.T) {
		err := params.checkForMissingMandatory()
		assert.Nil(t, err)
	})
	t.Run("false positive", func(t *testing.T) {
		params[1].used = false
		err := params.checkForMissingMandatory()
		assert.NotNil(t, err)
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
		// type faultyStruct struct {
		// 	a int
		// 	b interface{}
		// }
		// faultyStructType := reflect.TypeOf(faultyStruct{})
		// testStruct := &faultyStruct{}
		// params, err := newParameters(faultyStructType)
		// assert.Nil(t, err)
		// remaining, err := params.ParseArguments(testStruct, []string{"--b", "hello"})
		// assert.Nil(t, remaining)
		// assert.NotNil(t, err)
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
	t.Run("error on missing mandatory", func(t *testing.T) {
		type Foo struct {
			Foo int `yagclif:"mandatory"`
		}
		fooType := reflect.TypeOf(Foo{})
		testStruct := &Foo{}
		params, err := newParameters(fooType)
		assert.Nil(t, err)
		remaining, err := params.ParseArguments(testStruct, []string{})
		assert.Nil(t, remaining)
		assert.NotNil(t, err)
	})
}
func TestParse(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		testStruct := &validStruct{}
		os.Args = []string{"main", "hello", "--b", "world", "!"}
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
		os.Args = []string{"main", "hello", "-b", "world", "!"}
		remaining, err := Parse(testStruct)
		assert.Nil(t, remaining)
		assert.NotNil(t, err)
	})
}
