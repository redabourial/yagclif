package common

import (
	"fmt"
)

type testingInterface interface {
	Fatalf(format string, args ...interface{})
}

func AssertEqual(value interface{}, expected interface{}, t testingInterface, errorMsg ...interface{}) {
	usermsg := ""
	for _, msgElement := range errorMsg {
		usermsg += " " + fmt.Sprint(msgElement)
	}
	if fmt.Sprint(value) != fmt.Sprint(expected) {
		msg := fmt.Sprint(
			usermsg, "\r\n\t",
			"value : ", value, "\r\n\t",
			"expected : ", expected, "\r\n",
		)
		t.Fatalf(msg)
	}
}
