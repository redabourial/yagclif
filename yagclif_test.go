package yagclif

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrependToArray(t *testing.T) {
	s := prependToArray([]string{"hello", "world"}, "|")
	assert.Equal(t, "|hello\r\n|world\r\n", s)
}
func TestAddRoute(t *testing.T) {
	app := NewCliApp("Hello", "simple hello worlds")
	err := app.AddRoute("echo", "", func(args []string) {})
	assert.Nil(t, err)
	assert.NotNil(t, app.routes["echo"])
	err = app.AddRoute("echo", "", func(args []string) {})
	assert.NotNil(t, err)
}

func TestRun(t *testing.T) {
	type EmbededStruct struct {
		Embeded int `yagclif:"default:42"`
	}
	type TestStruct struct {
		EmbededStruct
		A []int `yagclif:"delimiter:,;mandatory"`
	}
	app := NewCliApp("Hello", "simple hello worlds")
	var passedStruct *TestStruct
	err := app.AddRoute("echo", "echoes the args", func(testStruct TestStruct, args []string) {
		passedStruct = &testStruct
	})
	assert.Nil(t, err)
	t.Run("works", func(t *testing.T) {
		os.Args = []string{"./main", "echo", "--a", "42,43", "world"}
		err = app.RunNoPanic(false)
		assert.Nil(t, err)
		assert.NotNil(t, passedStruct)
		assert.Equal(t, TestStruct{
			EmbededStruct: EmbededStruct{Embeded: 42},
			A:             []int{42, 43},
		}, *passedStruct)
	})
	t.Run("missing route", func(t *testing.T) {
		os.Args = []string{"./main", "missingAction", "--a", "42,43", "world"}
		err = app.RunNoPanic(false)
		assert.NotNil(t, err)
	})
	t.Run("no args", func(t *testing.T) {
		os.Args = []string{"./main"}
		err = app.RunNoPanic(false)
		assert.NotNil(t, err)
		HelpErr := app.RunNoPanic(true)
		assert.NotNil(t, HelpErr)
		assert.NotEqual(t, HelpErr, err)
	})
	t.Run("runtime error", func(t *testing.T) {
		os.Args = []string{"./main", "panic"}
		err := app.AddRoute("panic", "just panic", func(args []string) {
			panic("u..rge...to panic .... can't help it !")
		})
		assert.Nil(t, err)
		err = app.RunNoPanic(false)
		assert.NotNil(t, err)
	})
}
func TestGetHelp(t *testing.T) {
	type Context struct {
		AT int    `yagclif:"shortname:a;description:imA;mandatory"`
		BT string `yagclif:"description:FOO;default:someDefaultValue;env:testgethelp"`
	}
	os.Setenv("testgethelp", "something")
	fmt.Println("------------SAMPLE APP HELP------------")
	app := NewCliApp("Hello", "simple hello worlds")
	err := app.AddRoute("echo", "echoes the args", func([]string) {})
	assert.Nil(t, err)
	err = app.AddRoute("someAction", "does stuff", func(Context, []string) {})
	assert.Nil(t, err)
	help := app.GetHelp()
	assert.Contains(t, help, "Hello")
	assert.Contains(t, help, "simple hello worlds")
	assert.Contains(t, help, "someAction : does stuff")
	assert.Contains(t, help, "--at -a int (mandatory): imA")
	assert.Contains(t, help, "--bt string (default=someDefaultValue;env={key:testgethelp,value:something}): FOO")
}
