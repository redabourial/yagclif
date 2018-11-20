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
	os.Args = []string{"./main", "echo", "--a", "42,43", "world"}
	err = app.Run()
	assert.Nil(t, err)
	assert.NotNil(t, passedStruct)
	assert.Equal(t, TestStruct{
		A: []int{42, 43},
	}, *passedStruct)
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
