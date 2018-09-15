package arguments

import (
	"errors"
	"strconv"

	"./argument"
)

type arguments []argument.Argument

type Arguments interface {
	UseArgument(argumentName string) error
	findArgumentIndex(argumentName string) (int, error)
	ParseStringParameter(argumentName string) (string, error)
	ParseIntParameter(argumentName string) (int, error)
	ParseStringArrayParameter(argumentName string) (int, error)
	ParseIntArrayParameter(argumentName string) (int, error)
	GetUnused() []string
}

func New(args []string) Arguments {
	argObj := make(arguments, 0)
	for _, argText := range args {
		newArg := argument.New(argText)
		argObj = append(argObj, newArg)
	}
	return &argObj
}

func (args arguments) UseArgument(argumentName string) error {
	argIndex, err := args.findArgumentIndex(argumentName)
	if err == nil {
		args[argIndex].Use()
		return nil
	}
	return err
}

func (args arguments) findArgumentIndex(argumentName string) (int, error) {
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

func (args arguments) ParseStringParameter(argumentName string) (string, error) {
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

func (args arguments) ParseIntParameter(argumentName string) (int, error) {
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
func (args arguments) ParseStringArrayParameter(argumentName string) (int, error) {
	// TODO implement
	return 0, nil
}

func (args arguments) ParseIntArrayParameter(argumentName string) (int, error) {
	// TODO implement
	return 0, nil
}

func (args arguments) GetUnused() []string {
	unuseds := make([]string, 0)
	for _, arg := range args {
		if arg.IsNotUsed() {
			unuseds = append(unuseds, arg.Text())
		}
	}
	return unuseds
}
