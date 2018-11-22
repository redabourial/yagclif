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
	type TestStruct struct {
		A []int `yagclif:"delimiter:,;mandatory"`
	}
	app := NewCliApp("Hello", "simple hello worlds")
	var passedStruct *TestStruct
	var remainingArgs []string
	err := app.AddRoute("echo", "echoes the args", func(testStruct TestStruct, args []string) {
		passedStruct = &testStruct
		remainingArgs = args
	})
	assert.Nil(t, err)
	t.Run("works", func(t *testing.T) {
		os.Args = []string{"./main", "echo", "--a", "42,43", "world"}
		err = app.Run(false)
		assert.Nil(t, err)
		assert.NotNil(t, passedStruct)
		// assert.Equal(t, TestStruct{
		// 	A: []int{42, 43},
		// }, *passedStruct)
	})
	t.Run("missing route", func(t *testing.T) {
		os.Args = []string{"./main", "missingAction", "--a", "42,43", "world"}
		err = app.Run(false)
		assert.NotNil(t, err)
	})
	t.Run("no args", func(t *testing.T) {
		os.Args = []string{"./main"}
		err = app.Run(false)
		assert.NotNil(t, err)
		HelpErr := app.Run(true)
		assert.NotNil(t, HelpErr)
		assert.NotEqual(t, HelpErr, err)
	})
	t.Run("runtime error", func(t *testing.T) {
		os.Args = []string{"./main", "panic"}
		err := app.AddRoute("panic", "just panic", func(args []string) {
			panic("u..rge...to panic .... can't help it !")
		})
		assert.Nil(t, err)
		err = app.Run(false)
		assert.NotNil(t, err)
	})
}
func TestGetHelp(t *testing.T) {
	type Context struct {
		A int    `yagclif:"description:imA;mandatory"`
		B string `yagclif:"description:FOO;default:someDefaultValue"`
	}
	fmt.Println("------------SAMPLE APP HELP------------")
	app := NewCliApp("Hello", "simple hello worlds")
	err := app.AddRoute("echo", "echoes the args", func([]string) {})
	assert.Nil(t, err)
	err = app.AddRoute("someAction", "does stuff", func(Context, []string) {})
	assert.Nil(t, err)
	fmt.Print(app.GetHelp())
	fmt.Println("----------------------------------------------")
}
