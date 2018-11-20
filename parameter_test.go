package yagclif

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	p := parameter{
		delimiter: "-",
	}
	splittedValues := p.Split("hello-world-!")
	assert.Equal(t, []string{"hello", "world", "!"}, splittedValues)
}
func TestHasShortName(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		p := parameter{
			shortName: "hello",
		}
		assert.True(t, p.hasShortName())
	})
	t.Run("negative", func(t *testing.T) {
		p := parameter{}
		assert.False(t, p.hasShortName())
	})
}
func TestUsed(t *testing.T) {
	p := parameter{}
	assert.False(t, p.Used())
	p.Use()
	assert.True(t, p.Used())
}

func TestFillParameter(t *testing.T) {
	t.Run("Works", func(t *testing.T) {
		param := &parameter{}
		assert.EqualValues(
			t,
			nil,
			param.fillParameter("description:42"),
			param.fillParameter("shortname:43"),
			param.fillParameter("mandatory"),
			param.fillParameter("delimiter:;"),
			param.fillParameter("default:44"),
		)
		assert.Equal(t, parameter{
			description:  "42",
			shortName:    "43",
			mandatory:    true,
			delimiter:    ";",
			defaultValue: "44",
		}, *param)
	})
	t.Run("splitError", func(t *testing.T) {
		param := &parameter{}
		assert.NotNil(t, param.fillParameter("description::"))
	})
	t.Run("error", func(t *testing.T) {
		param := &parameter{}
		assert.NotNil(t, param.fillParameter("something"))
	})
}

func TestSplitConstraint(t *testing.T) {
	t.Run("With value", func(t *testing.T) {
		kv, err := splitConstraint("hello:world")
		assert.Nil(t, err)
		assert.Equal(t, kv, keyValuePair{
			key:   "hello",
			value: "world",
		})
	})
	t.Run("Without value", func(t *testing.T) {
		kv, err := splitConstraint("hello")
		assert.Nil(t, err)
		assert.Equal(t, kv, keyValuePair{
			key:   "hello",
			value: "",
		})
	})
	t.Run("too many :::", func(t *testing.T) {
		_, err := splitConstraint("hello:::")
		assert.NotNil(t, err)
	})
}

func TestMatches(t *testing.T) {
	t.Run("With shortname", func(t *testing.T) {
		param := parameter{
			name:      "hello",
			shortName: "h",
		}
		t.Run("negative", func(t *testing.T) {
			assert.False(t, param.Matches("hell"))
		})
		t.Run("positives", func(t *testing.T) {
			assert.True(t, param.Matches("--hello"))
			assert.True(t, param.Matches("-h"))
		})
	})
	t.Run("Without shortname", func(t *testing.T) {
		param := parameter{
			name: "hello",
		}
		t.Run("negative", func(t *testing.T) {
			assert.False(t, param.Matches("hell"))
		})
		t.Run("positives", func(t *testing.T) {
			assert.True(t, param.Matches("--hello"))
		})
	})
}
func TestNewParameter(t *testing.T) {
	type foo struct {
		a int `yagclif:"something"`
		b int `yagclif:"mandatory;shortname:c"`
	}
	t.Run("Works", func(t *testing.T) {
		field := reflect.TypeOf(foo{}).Field(1)
		param, err := newParameter(field)
		assert.Nil(t, err)
		assert.True(t, param.mandatory)
		assert.True(t, param.Matches("--b"))
		assert.True(t, param.Matches("-c"))
		assert.False(t, param.Matches("-b"))
	})
	t.Run("delimtiter", func(t *testing.T) {
		t.Run("on empty delimiter", func(t *testing.T) {
			type foo struct {
				b []string `yagclif:""`
			}
			field := reflect.TypeOf(foo{}).Field(0)
			param, err := newParameter(field)
			assert.Nil(t, err)
			assert.Equal(t, ",", param.delimiter)
		})

		t.Run("on delimiter", func(t *testing.T) {
			type foo struct {
				a []int `yagclif:"delimiter:!"`
			}
			field := reflect.TypeOf(foo{}).Field(0)
			param, err := newParameter(field)
			assert.Nil(t, err)
			assert.Equal(t, param.delimiter, "!")
		})
	})
	t.Run("Returns error", func(t *testing.T) {
		field := reflect.TypeOf(foo{}).Field(0)
		_, err := newParameter(field)
		assert.NotNil(t, err)
	})
}

func TestGetValue(t *testing.T) {
	type foo struct {
		Bar int
	}
	field := reflect.TypeOf(foo{}).Field(0)
	param, err := newParameter(field)
	assert.Nil(t, err)
	fooVar := &foo{Bar: 0}
	barValue := param.getValue(fooVar)
	assert.Equal(t, 0, fooVar.Bar)
	barValue.SetInt(int64(42))
	assert.Equal(t, 42, fooVar.Bar)
}

func TestSetterCallBacks(t *testing.T) {
	//TODO organise by complexity
	type foo struct {
		Bar bool
		Car string
		Dar int
		Ear []string `yagclif:"delimiter:,"`
		Far []int    `yagclif:"delimiter:,"`
		Gar interface{}
	}
	t.Run("Set Bool", func(t *testing.T) {
		field := reflect.TypeOf(foo{}).Field(0)
		param, err := newParameter(field)
		assert.Nil(t, err)
		fooVar := &foo{Bar: false}
		callBack, err := param.SetterCallback(fooVar)
		assert.Nil(t, err, callBack)
		assert.True(t, fooVar.Bar)
	})
	t.Run("Set String", func(t *testing.T) {
		field := reflect.TypeOf(foo{}).Field(1)
		param, err := newParameter(field)
		assert.Nil(t, err)
		fooVar := &foo{Car: "Hello"}
		callBack, err := param.SetterCallback(fooVar)
		assert.Nil(t, err)
		err = callBack("world")
		assert.Nil(t, err)
		assert.Equal(t, "world", fooVar.Car)
	})
	t.Run("Set Int", func(t *testing.T) {
		field := reflect.TypeOf(foo{}).Field(2)
		param, err := newParameter(field)
		assert.Nil(t, err)
		fooVar := &foo{Dar: 1}
		callBack, err := param.SetterCallback(fooVar)
		assert.Nil(t, err)
		t.Run("works", func(t *testing.T) {
			err := callBack("2")
			assert.Nil(t, err)
			assert.Equal(t, 2, fooVar.Dar)
		})
		t.Run("returns Error", func(t *testing.T) {
			err := callBack("q")
			assert.NotNil(t, err)
		})
	})
	t.Run("Set String Array", func(t *testing.T) {
		field := reflect.TypeOf(foo{}).Field(3)
		param, err := newParameter(field)
		assert.Nil(t, err)
		fooVar := &foo{}
		callBack, err := param.SetterCallback(fooVar)
		assert.Nil(t, err)
		err = callBack("hello,world")
		assert.Nil(t, err)
		assert.Equal(t, []string{"hello", "world"}, fooVar.Ear)
	})
	t.Run("Set Int Array", func(t *testing.T) {
		field := reflect.TypeOf(foo{}).Field(4)
		param, err := newParameter(field)
		assert.Nil(t, err)
		t.Run("works", func(t *testing.T) {
			fooVar := &foo{}
			callBack, err := param.SetterCallback(fooVar)
			assert.Nil(t, err)
			err = callBack("1,2,3")
			assert.Nil(t, err)
			assert.Equal(t, []int{1, 2, 3}, fooVar.Far)
		})
		t.Run("returns error", func(t *testing.T) {
			fooVar := &foo{}
			callBack, err := param.SetterCallback(fooVar)
			assert.Nil(t, err)
			err = callBack("hello world")
			assert.Nil(t, fooVar.Far)
			assert.NotNil(t, err)
		})
	})
	t.Run("returns error", func(t *testing.T) {
		field := reflect.TypeOf(foo{}).Field(5)
		param, err := newParameter(field)
		assert.Nil(t, err)
		fooVar := &foo{}
		callBack, err := param.SetterCallback(fooVar)
		assert.Nil(t, callBack)
		assert.NotNil(t, err)
	})
}

func TestParamHelp(t *testing.T) {
	stringContains := func(s string, substrs ...string) {
		for _, substr := range substrs {
			if !strings.Contains(s, substr) {
				t.Error("expectected help to contain : ", substr, "\r\n", "help :\r\n", s)
			}
		}
	}
	stringDoesnotContain := func(s string, substrs ...string) {
		for _, substr := range substrs {
			if strings.Contains(s, substr) {
				t.Error("expectected help to not contain : ", substr, "\r\n", "help :\r\n", s)
			}
		}
	}
	t.Run("bool type", func(t *testing.T) {
		param := parameter{
			name: "Bar",
			tipe: reflect.TypeOf(true),
		}
		help := param.GetHelp()
		stringContains(help, "--bar", "bool")
		stringDoesnotContain(help, ":", "mandatory")
	})
	t.Run("string type", func(t *testing.T) {
		param := parameter{
			name:        "Bar",
			tipe:        reflect.TypeOf(""),
			description: "some string",
		}
		help := param.GetHelp()
		stringContains(help, "--bar", "string", ":", "some string")
		stringDoesnotContain(help, "mandatory")
	})

	t.Run("int type with description", func(t *testing.T) {
		param := parameter{
			name:        "Bar",
			tipe:        reflect.TypeOf(""),
			mandatory:   true,
			description: "some int",
		}
		help := param.GetHelp()
		stringContains(help, "--bar", "string", ":", "some int", "mandatory")
	})
	t.Run("string array ", func(t *testing.T) {
		param := parameter{
			name:      "Bar",
			tipe:      reflect.TypeOf([]string{}),
			delimiter: ",",
			mandatory: true,
		}
		help := param.GetHelp()
		stringContains(help, "--bar", "[]string", "delimiter", ",", "mandatory")
	})
	t.Run("string array type", func(t *testing.T) {
		type foo struct {
			Bar []int `yagclif:"delimiter:,;mandatory;description:"`
		}
		param := parameter{
			name:        "Bar",
			tipe:        reflect.TypeOf([]int{}),
			delimiter:   " ",
			description: "some int array",
			mandatory:   true,
		}
		help := param.GetHelp()
		stringContains(help, "--bar", "[]int", "delimiter", "whitespace", "some int array", "mandatory")
	})
}
