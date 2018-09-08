package argumentparse

import (
	"fmt"
	"testing"

	"../common"
)

const argumentName, parameterValue = "--argument", "parameter"

var exampleArgs = [][]string{
	{argumentName},
	{argumentName, parameterValue},
	{argumentName, parameterValue, "after"},
	{argumentName, "42"},
	{"somethingBefore", argumentName},
	{"somethingBefore", argumentName, "42"},
	{"somethingBefore", argumentName, parameterValue, "after"},
}

var exampleIndexes = []int{
	0,
	0,
	0,
	0,
	1,
	1,
	1,
}

func TestFindArgumentIndex(t *testing.T) {
	errValue, err := findArgumentIndex(&(exampleArgs[3]), "foo")
	if err == nil {
		t.Fatalf("doesn't return error on exception")
	}
	common.AssertEqual(errValue, -1, t, "doesn't return -1 on error")
	for index, example := range exampleArgs {
		argIndex, err := findArgumentIndex(&example, argumentName)
		common.AssertEqual(err, nil, t, "failed with example ", index, " : unexpected error")
		common.AssertEqual(argIndex, exampleIndexes[index], t, "failed with example ", index, " : wrong value")
	}
}

func compareArrays(t *testing.T, array []string, expected []string, index int) {
	common.AssertEqual(
		len(array), len(expected), t,
		"failed with example ", index, " : size mismatch ",
	)
	for elementNumber, expectedValue := range array {
		if expected[elementNumber] != expectedValue {
			common.AssertEqual(
				expectedValue, expected[elementNumber], t,
				"failed with example ", index,
				" : wrong value on element ", elementNumber,
			)
		}
	}
}

var exampleParsesString = [][]interface{}{
	[]interface{}{nil, nil, true},
	[]interface{}{parameterValue, []string{}, false},
	[]interface{}{parameterValue, []string{"after"}, false},
	[]interface{}{"42", []string{}, false},
	[]interface{}{nil, nil, true},
	[]interface{}{"42", []string{"somethingBefore"}, false},
	[]interface{}{parameterValue, []string{"somethingBefore", "after"}, false},
}

func TestParseArgumentStringParameter(t *testing.T) {
	stringParameter, arrayReturned, err := ParseArgumentStringParameter(&exampleArgs[0], "foo")
	compareArrays(t, *arrayReturned, exampleArgs[0], 0)
	if err == nil {
		t.Fatalf("doesn't return error on exception")
	}
	common.AssertEqual(stringParameter, nil, t, "doesn't return nil on error")
	for index, args := range exampleArgs {
		stringParameter, arrayReturned, err := ParseArgumentStringParameter(&args, argumentName)
		if err == nil && exampleParsesString[index][2] == true {
			msg := fmt.Sprint("failed with example  ", index, " : error was returned")
			t.Fatalf(msg)
		} else if exampleParsesString[index][0] != nil && stringParameter == nil {
			t.Fatalf(fmt.Sprint("failed with example ", index, " : returned parameter value is not right "))
		} else if exampleParsesString[index][1] != nil && arrayReturned == nil {
			t.Fatalf(fmt.Sprint("failed with example ", index, " : returned array is nil "))
		} else if arrayReturned != nil {
			// msg := fmt.Sprint("original array", args, "data : ", *arrayReturned, " array : ", exampleParsesString[index][1])
			// fmt.Println(msg)
			expectedArray, ok := exampleParsesString[index][1].([]string)
			if ok {
				compareArrays(t, *arrayReturned, expectedArray, index)
			} else {
				compareArrays(t, *arrayReturned, args, index)
			}
		} else if arrayReturned == nil {
			t.Fatalf(fmt.Sprint("failed with example ", index, " : returned array is nils "))
		}
		if err != nil && exampleParsesString[index][2] != true {
			t.Fatalf(fmt.Sprint("failed with example ", index, " : error wasn't returned"))
		}
	}
}

var exampleParsesInt = [][]interface{}{
	[]interface{}{nil, nil, true},
	[]interface{}{nil, nil, true},
	[]interface{}{nil, nil, true},
	[]interface{}{42, []string{}, false},
	[]interface{}{nil, nil, true},
	[]interface{}{42, []string{"somethingBefore"}, false},
	[]interface{}{nil, nil, true},
}

func TestParseArgumentIntParameter(t *testing.T) {
	for index, args := range exampleArgs {
		intParameter, arrayReturned, err := ParseArgumentIntParameter(&args, argumentName)
		if exampleParsesInt[index][2] == true && err == nil {
			t.Fatalf(fmt.Sprint("failed with example  ", index, " : expected error was not returned"))
		} else if exampleParsesInt[index][2] == false && err != nil {
			t.Fatalf(fmt.Sprint("failed with example  ", index, " : unexpected error was returned"))
		} else if exampleParsesInt[index][2] == false && intParameter == -1 {
			t.Fatalf(fmt.Sprint("failed with example ", index, " : returned parameter value is not right "))
		} else if exampleParsesInt[index][1] != nil && arrayReturned == nil {
			t.Fatalf(fmt.Sprint("failed with example ", index, " : returned array is nil "))
		} else if exampleParsesInt[index][1] != nil && arrayReturned != nil {
			if expectedArray, ok := exampleParsesInt[index][1].([]string); ok {
				// msg := fmt.Sprint("index ", index, " original array", args, "actual : ", *arrayReturned, " expected : ", expectedArray)
				// fmt.Println(msg)
				compareArrays(t, *arrayReturned, expectedArray, index)
			} else {
				compareArrays(t, *arrayReturned, args, index)
			}
		} else if arrayReturned == nil {
			t.Fatalf(fmt.Sprint("failed with example ", index, " : returned array is nil "))
		}
		if err != nil && exampleParsesInt[index][2] != true {
			t.Fatalf(fmt.Sprint("failed with example ", index, " : error wasn't returned"))
		}
	}
}

func TestExists(t *testing.T) {
	expectedFalse, arrayReturned := Exists(&exampleArgs[1], "foo")
	compareArrays(t, *arrayReturned, exampleArgs[1], 0)
	common.AssertEqual(expectedFalse, false, t, "false positive")
	expectedTrue, arrayReturned := Exists(&exampleArgs[1], exampleArgs[1][1])
	fmt.Println(arrayReturned)
	common.AssertEqual(*arrayReturned, []string{"--argument"}, t, "false negative")
	common.AssertEqual(expectedTrue, true, t, "false negative")
}
