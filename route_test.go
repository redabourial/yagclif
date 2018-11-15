package cliced

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCustomCallBackType(t *testing.T) {
	type SomeStruct struct {
		a int
		b string
	}
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
	var passedArgs []string = nil
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

	})
	t.Run("simple callback", func(t *testing.T) {

	})
	t.Run("error", func(t *testing.T) {

	})
}
