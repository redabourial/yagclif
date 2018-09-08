package tagparse

import (
	"encoding/json"
	"testing"
)

type testStruct struct {
	string `json:",string"`
}

func testStructTags(t *testing.T) {
	json.Unmarshal
}
