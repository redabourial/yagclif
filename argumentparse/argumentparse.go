package argumentparse

import (
	"errors"
	"strconv"
)

func findArgumentIndex(args *[]string, argumentName string) (int, error) {
	for index, argument := range *args {
		if argument == argumentName {
			return index, nil
		}
	}
	return -1, errors.New("error argument not found")
}

func copyArgs(args *[]string, removeIndexes ...int) *[]string {
	isRemoveIndex := func(index int) bool {
		for _, removeIndex := range removeIndexes {
			if index == removeIndex {
				return true
			}
		}
		return false
	}
	newArgs := make([]string, 0)
	for index, arg := range *args {
		if !isRemoveIndex(index) {
			newArgs = append(newArgs, arg)
		}
	}
	return &newArgs
}

func ParseArgumentStringParameter(args *[]string, argumentName string) (*string, *[]string, error) {
	argIndex, err := findArgumentIndex(args, argumentName)
	if err != nil {
		return nil, args, err
	}
	if len(*args) > (argIndex + 1) {
		parameter := (*args)[argIndex+1]
		newArgs := copyArgs(args, argIndex, argIndex+1)
		return &parameter, newArgs, nil
	}
	return nil, args, errors.New("argument doesn't have parameter")
}

func ParseArgumentIntParameter(args *[]string, argumentName string) (int, *[]string, error) {
	parameter, remainingArgs, err := ParseArgumentStringParameter(args, argumentName)
	if err != nil {
		return -1, args, err
	}
	parameterInt, err := strconv.Atoi(*parameter)
	if err != nil {
		return -1, args, err
	}
	return parameterInt, remainingArgs, nil
}

func Exists(args *[]string, argumentName string) (bool, *[]string) {
	argIndex, err := findArgumentIndex(args, argumentName)
	if err != nil {
		return false, args
	}
	return true, copyArgs(args, argIndex)
}
