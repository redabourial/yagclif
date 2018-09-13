package arguments

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var args = []string{"alpha", "beta", "gamma", "delta"}

func Test_FindArgumentIndex(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		arguments := InitArguments(args)
		index, err := arguments.findArgumentIndex("beta")
		assert.Equal(t, index, 1)
		assert.Nil(t, err, "Unexpected error")
	})
	t.Run("argument not found error", func(t *testing.T) {
		arguments := InitArguments(args)
		_, err := arguments.findArgumentIndex("epsilone")
		assert.NotNil(t, err, "expecting error found nil")
	})
	t.Run("used argument error", func(t *testing.T) {
		arguments := InitArguments(args)
		arguments.UseArgument("beta")
		_, err := arguments.findArgumentIndex("beta")
		assert.NotNil(t, err, "expecting error found nil")
	})
}

func Test_UseArgument(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		arguments := InitArguments(args)
		err := arguments.UseArgument("beta")
		assert.Nil(t, err, "unexpected error")
		unuseds := arguments.GetUnused()
		assert.Equal(t, []string{"alpha", "gamma", "delta"}, unuseds, "arguments not marked as used")
	})
	t.Run("returns error", func(t *testing.T) {
		arguments := InitArguments(args)
		err := arguments.UseArgument("epsilone")
		assert.NotNil(t, err, "expecting error found nil")
	})
}

func Test_ParseStringParameter(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		arguments := InitArguments(args)
		param, err := arguments.ParseStringParameter("beta")
		assert.Nil(t, err, "expected error")
		assert.Equal(t, "gamma", param, "expected error")
		unuseds := arguments.GetUnused()
		assert.Equal(t, []string{"alpha", "delta"}, unuseds, "arguments not marked as used")
	})
	t.Run("error from findArgumentIndex ", func(t *testing.T) {
		arguments := InitArguments(args)
		_, err := arguments.ParseStringParameter("epsilone")
		assert.NotNil(t, err, "expected error")
	})
	t.Run("used argument error", func(t *testing.T) {
		arguments := InitArguments(args)
		err := arguments.UseArgument("gamma")
		_, err = arguments.ParseStringParameter("beta")
		assert.NotNil(t, err, "expected error")
	})
	t.Run("empty parameter error", func(t *testing.T) {
		arguments := InitArguments(args)
		_, err := arguments.ParseStringParameter("delta")
		assert.NotNil(t, err, "expected error")
	})
}

func Test_ParseIntParameter(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		arguments := InitArguments([]string{"alpha", "beta", "2"})
		param, err := arguments.ParseIntParameter("beta")
		assert.Nil(t, err, "expected error")
		assert.Equal(t, 2, param, "expected error")
		unuseds := arguments.GetUnused()
		assert.Equal(t, []string{"alpha"}, unuseds, "arguments not marked as used when they are")
	})
	t.Run("error from ParseStringParameter ", func(t *testing.T) {
		arguments := InitArguments(args)
		_, err := arguments.ParseIntParameter("epsilone")
		assert.NotNil(t, err, "expected error")
	})
	t.Run("error from strconv.Atoi", func(t *testing.T) {
		arguments := InitArguments(args)
		_, err := arguments.ParseIntParameter("alpha")
		assert.NotNil(t, err, "expected error")
	})
}

func Test_getUnused(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		arguments := InitArguments(args)
		param := arguments.GetUnused()
		assert.Equal(t, args, param, "unused argument is broken")
	})
	t.Run("error from ParseStringParameter ", func(t *testing.T) {
		arguments := InitArguments(args)
		_, err := arguments.ParseIntParameter("epsilone")
		assert.NotNil(t, err, "expected error")
	})
	t.Run("error from strconv.Atoi", func(t *testing.T) {
		arguments := InitArguments(args)
		_, err := arguments.ParseStringParameter("alpha")
		assert.Nil(t, err, "expected error")
	})
}
