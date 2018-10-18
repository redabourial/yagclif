package cliced

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type validStruct struct {
	a int
	b string `cliced:"mandatory;shortname:sb"`
	c bool   `cliced:"mandatory;shortname:sc"`
}
type faultyStruct struct {
	a int
	b string `cliced:"something"`
}

func TestNewParameters(t *testing.T) {

	t.Run("returns value", func(t *testing.T) {
		params, err := newParameters(validStruct{})
		assert.Nil(t, err)
		assert.Equal(t, 3, len(params))
		assert.Equal(t, "a", params[0].Name())
		assert.Equal(t, 0, params[0].Index())
		assert.Equal(t, "b", params[1].Name())
		assert.True(t, params[1].Mandatory())
		assert.Equal(t, 1, params[1].Index())
		assert.Equal(t, "c", params[2].Name())
		assert.True(t, params[2].Mandatory())
		assert.Equal(t, 2, params[2].Index())
	})
	t.Run("returns error", func(t *testing.T) {
		params, err := newParameters(faultyStruct{})
		assert.NotNil(t, err)
		assert.Nil(t, params)
	})
}
