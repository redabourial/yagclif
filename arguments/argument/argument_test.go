package argument

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	arg, assertEqual := New("hello"), func(expected interface{}, actual interface{}) {
		assert.Equal(t, expected, actual, "")
	}
	test := func(used bool) {
		assertEqual(&argument{"hello", used}, arg)
		assertEqual("hello", arg.Text())
		assertEqual(arg.Equals("hello"), true)
		assertEqual(arg.Equals("hell"), false)
		assertEqual(arg.IsUsed(), used)
		assertEqual(arg.IsNotUsed(), !used)
	}
	test(false)
	arg.Use()
	test(true)
}
