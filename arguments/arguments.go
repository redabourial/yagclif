package arguments

import (
	"errors"
	"strconv"

	"./argument"
)

type Args []argument.Argument

func InitArguments(arguments []string) *Args {
	argObj := make(Args, 0)
	for _, argText := range arguments {
		newArg := argument.New(argText)
		argObj = append(argObj, newArg)
	}
	return &argObj
}

func (args Args) useArgument(argumentName string) error {
	argIndex, err := args.findArgumentIndex(argumentName)
	if err == nil {
		args[argIndex].Use()
		return nil
	}
	return err
}

func (args Args) findArgumentIndex(argumentName string) (int, error) {
	for index, argument := range args {
		if argument.Equals(argumentName) {
			if argument.IsNotUsed() {
				return index, nil
			}
			return -1, errors.New("error argument is used")
		}
	}
	return -1, errors.New("error argument not found")
}

func (args Args) ParseStringParameter(argumentName string) (string, error) {
	argIndex, err := args.findArgumentIndex(argumentName)
	if err != nil {
		return "", err
	}
	if len(args) > (argIndex + 1) {
		if args[argIndex].IsNotUsed() && args[argIndex+1].IsNotUsed() {
			args[argIndex].Use()
			args[argIndex+1].Use()
			return args[argIndex+1].Text(), nil
		}
		return "", errors.New("argument is used")
	}
	return "", errors.New("argument doesn't have parameter")
}

func (args Args) ParseIntParameter(argumentName string) (int, error) {
	parameter, err := args.ParseStringParameter(argumentName)
	if err != nil {
		return 0, err
	}
	parameterInt, err := strconv.Atoi(parameter)
	if err != nil {
		return 0, err
	}
	return parameterInt, nil
}
func (args Args) ParseStringArrayParameter(argumentName string) (int, error) {
	// TODO implement
	return 0, nil
}

func (args Args) ParseIntArrayParameter(argumentName string) (int, error) {
	// TODO implement
	return 0, nil
}

func (args Args) getUnused() []string {
	unuseds := make([]string, 0)
	for _, arg := range args {
		if arg.IsNotUsed() {
			unuseds = append(unuseds, arg.Text())
		}
	}
	return unuseds
}
