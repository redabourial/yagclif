package yagclif

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SomeStruct struct {
	A int `yagclif:"mandatory"`
	B string
}

func TestGetCustomCallBackType(t *testing.T) {
	t.Run("args[] callback", func(t *testing.T) {
		argsCallback := func(args []string) {}
		tipe, err := getCustomCallBackType(argsCallback)
		assert.Nil(t, err)
		assert.Nil(t, tipe)
	})
	t.Run("args[], somestruct callback", func(t *testing.T) {
		customCallBack := func(data SomeStruct, args []string) {}
		tipe, err := getCustomCallBackType(customCallBack)
		assert.Nil(t, err)
		assert.Equal(t, tipe, reflect.TypeOf(SomeStruct{}))
	})
	t.Run("return error on nil interface", func(t *testing.T) {
		tipe, err := getCustomCallBackType(nil)
		assert.NotNil(t, err)
		assert.Nil(t, tipe)
	})
	t.Run("returns error on unsupported callback type", func(t *testing.T) {
		customCallBack := func(args []string, data SomeStruct) {}
		tipe, err := getCustomCallBackType(customCallBack)
		assert.NotNil(t, err)
		assert.Nil(t, tipe)
	})
}

func TestGetSimpleCallBack(t *testing.T) {
	var passedArgs []string
	t.Run("works", func(t *testing.T) {
		stub := func(args []string) {
			passedArgs = args
		}
		stubValue := reflect.ValueOf(stub)
		standardCallback := getSimpleCallBack(stubValue)
		mockArgs := []string{"hello", "world"}
		err := standardCallback(mockArgs)
		assert.Nil(t, err)
		assert.Equal(t, mockArgs, passedArgs)
	})
	t.Run("error", func(t *testing.T) {
		stubValue := reflect.ValueOf("hello")
		standardCallback := getSimpleCallBack(stubValue)
		mockArgs := []string{"hello", "world"}
		err := standardCallback(mockArgs)
		assert.NotNil(t, err)
	})
}

func TestGetCustomCallBack(t *testing.T) {
	t.Run("new parameters error", func(t *testing.T) {
		callbackFunc := reflect.ValueOf(func([]string) {})
		callback, err := getCustomCallBack(callbackFunc, reflect.TypeOf("hello"))
		assert.NotNil(t, err)
		assert.Nil(t, callback)
	})
	t.Run("callBack works", func(t *testing.T) {
		var passedValue *SomeStruct
		callbackFunc := reflect.ValueOf(func(ss SomeStruct, remainingArgs []string) {
			passedValue = &ss
		})
		callback, err := getCustomCallBack(callbackFunc, reflect.TypeOf(SomeStruct{}))
		assert.Nil(t, err)
		assert.Nil(t, passedValue)
		err = callback([]string{"--a", "1", "hello"})
		assert.Nil(t, err)
		assert.NotNil(t, passedValue)
		assert.Equal(t, &SomeStruct{
			A: 1,
		}, passedValue)
	})
	t.Run("callBack formating error", func(t *testing.T) {
		callbackFunc := reflect.ValueOf(func(i int, remainingArgs []string) {
		})
		callback, err := getCustomCallBack(callbackFunc, reflect.TypeOf(SomeStruct{}))
		assert.Nil(t, err)
		err = callback([]string{})
		assert.NotNil(t, err)
	})
	t.Run("panic inside callback", func(t *testing.T) {
		callbackFunc := reflect.ValueOf(func(SomeStruct, []string) {
			panic("hello")
		})
		callback, err := getCustomCallBack(callbackFunc, reflect.TypeOf(SomeStruct{}))
		assert.Nil(t, err)
		err = callback([]string{"--a", "42"})
		assert.NotNil(t, err)
	})
}

func TestFormatCallBack(t *testing.T) {
	t.Run("getSimpleCallBack", func(t *testing.T) {
		callbackFunc := reflect.ValueOf(func(remainingArgs []string) {})
		callback, err := formatCallBack(callbackFunc, nil)
		assert.Nil(t, err)
		assert.Equal(t, reflect.TypeOf(callback), reflect.TypeOf(func([]string) error { return nil }))
	})
	t.Run("getCustomCallBack", func(t *testing.T) {
		callbackFunc := reflect.ValueOf(func(ss SomeStruct, remainingArgs []string) {})
		callback, err := formatCallBack(callbackFunc, reflect.TypeOf(SomeStruct{}))
		assert.Nil(t, err)
		assert.Equal(t, reflect.TypeOf(callback), reflect.TypeOf(func([]string) error { return nil }))
	})
}

func TestNewRoute(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		route, err := newRoute("hello", func(ss SomeStruct, remainingArgs []string) {})
		assert.Nil(t, err)
		assert.NotNil(t, route)
	})
	t.Run("getCustomCallBackType error", func(t *testing.T) {
		route, err := newRoute("hello", nil)
		assert.Nil(t, route)
		assert.NotNil(t, err)
	})
	t.Run("getCustomCallBackType error", func(t *testing.T) {
		route, err := newRoute("hello", func(i int, remainingArgs []string) {})
		assert.Nil(t, route)
		assert.NotNil(t, err)
	})
}

func TestRouteRun(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		var passedArgs []string
		r := route{
			formatedCallback: func(args []string) error {
				passedArgs = args
				return fmt.Errorf("hello")
			},
		}
		expectedArgs := []string{"hello", "world"}
		err := r.run(expectedArgs)
		assert.NotNil(t, err)
		assert.Equal(t, expectedArgs, passedArgs)
	})
	t.Run("error", func(t *testing.T) {
		r := route{}
		err := r.run([]string{})
		assert.NotNil(t, err)
	})
}

func TestRouteGetHelp(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		r := route{
			parameterType: reflect.TypeOf(SomeStruct{}),
		}
		helpTexts := r.getHelp()
		helpContains := func(s string) bool {
			for _, helpText := range helpTexts {
				if strings.Contains(helpText, s) {
					return true
				}
			}
			return false
		}
		assertContains := func(s string) {
			assert.True(t, helpContains(s))
		}
		assertContains("--a")
		assertContains("(mandatory)")
		assertContains("int")
	})
	t.Run("returns empty on nil parameterType", func(t *testing.T) {
		r := route{
			parameterType: nil,
		}
		helpTexts := r.getHelp()
		length := len(helpTexts)
		assert.Equal(t, 0, length)
	})
	t.Run("returns error on unparseable parameterType", func(t *testing.T) {
		r := route{
			parameterType: reflect.TypeOf(true),
		}
		helpTexts := r.getHelp()
		assert.Equal(t, []string{"Could not parse parameter type"}, helpTexts)
	})
}
