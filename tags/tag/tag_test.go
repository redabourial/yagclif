package tag

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func expectPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Fatalf("expected function to panic")
	}
}

func structTagFromString(s string) reflect.StructTag {
	return reflect.StructTag(s)
}

func newStructField(tag ...string) Tag {
	structTag := structTagFromString(strings.Join(tag, ""))
	return New(reflect.StructField{"name", "pkgPath", reflect.TypeOf(1), structTag, 0, []int{}, false}, 42)
}

func descriptionFromTag(tag ...string) string {
	structField := newStructField(strings.Join(tag, ""))
	return structField.GetArgumentHelp()
}
func Test_New(t *testing.T) {
	t.Run("empty tag", func(t *testing.T) {
		tag := newStructField("")
		assert.Equal(t, Tag{"name", 42, false, "", reflect.TypeOf(1)}, tag, "error")
	})
	t.Run("tag with description", func(t *testing.T) {
		tag := newStructField(TagName, ":", "\"desc\"")
		assert.Equal(t, Tag{"name", 42, false, "desc", reflect.TypeOf(1)}, tag, "error")
	})
	t.Run("mandatory tag", func(t *testing.T) {
		tag := newStructField(TagName, ":", "\"desc,mandatory\"")
		assert.Equal(t, Tag{"name", 42, true, "desc", reflect.TypeOf(1)}, tag, "error")
	})
	t.Run("panic on unregonized options in tag", func(t *testing.T) {
		defer expectPanic(t)
		tag := newStructField(TagName, ":", "\"desc,notMandatory\"")
		assert.Equal(t, Tag{"name", 42, true, "desc", reflect.TypeOf(1)}, tag, "error")
	})
	t.Run("panic on too many tag", func(t *testing.T) {
		defer expectPanic(t)
		newStructField(TagName, ":", "\"desc,mandatory,too,much\"")
	})
}

func Test_GetArgumentHelp(t *testing.T) {
	t.Run("empty tag", func(t *testing.T) {
		help := descriptionFromTag("")
		assert.Equal(t, "name int", help)
	})
	t.Run("tag with description", func(t *testing.T) {
		help := descriptionFromTag(TagName, ":", "\"desc\"")
		assert.Equal(t, "name int : desc", help)
	})
	t.Run("mandatory tag", func(t *testing.T) {
		help := descriptionFromTag(TagName, ":", "\"desc,mandatory\"")
		assert.Equal(t, "name int (mandatory) : desc", help)
	})
}

func Test_GetFromArgumentTag(t *testing.T) {
	t.Run("works", func(t *testing.T) {

	})
	t.Run("works", func(t *testing.T) {

	})
	t.Run("works", func(t *testing.T) {

	})
	t.Run("works", func(t *testing.T) {

	})
}
