package yagclif

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Name of the tag to parse.
const tagName = "yagclif"

// Value to prefix to name value.
const namePrefix = "--"

// Value to prefix to shortName value.
const shortNamePrefix = "-"

// Value of the delimiter between constraint
// key-value pairs.
const constraintValueDelimiter = ":"

// Value of the delimiter between constraints.
const constraintsDelimiter = ";"

// Struct for stroring key-value string pair
type keyValuePair struct {
	key   string
	value string
}

// Split a constraint as key-value constraint
func splitConstraint(constraint string) (keyValuePair, error) {
	parts := strings.Split(constraint, constraintValueDelimiter)
	switch len(parts) {
	case 1:
		return keyValuePair{
			parts[0], "",
		}, nil
	case 2:
		return keyValuePair{
			parts[0], parts[1],
		}, nil
	}
	return keyValuePair{}, fmt.Errorf("syntax error too many characters %s ", constraintValueDelimiter)
}

// Struct defining a parameter from a structField.
type parameter struct {
	// Name of the parameter arguments are tested
	// by appending an underscore to this value.
	name string
	// ShortName of the parameter argument.
	// ShortName matches are evaluated after
	// appending two underscore to this value.
	shortName string
	// Index of the structField this parameter
	// was created from.
	index int
	// Short description of the parameter
	// used for help and error messages.
	description string
	// If true not finding this parameter
	// will result is an error.
	mandatory bool
	// If true not finding this parameter
	// will result is an error.
	used bool
	// Value used to parse array types.
	delimiter string
	// Type of the parameter only types
	// bool,int,string,[]int,[]string are supported.
	tipe reflect.Type
	// Default Value
	defaultValue string
	// Default ENV key
	envKey string
}

// Returns Cli names (text before the parameter)
// as lowercase strings.
func (p *parameter) CliNames() []string {
	if p.hasShortName() {
		return []string{
			fmt.Sprint(namePrefix, strings.ToLower(p.name)),
			fmt.Sprint(shortNamePrefix, strings.ToLower(p.shortName)),
		}
	}
	return []string{
		fmt.Sprint(namePrefix, strings.ToLower(p.name)),
	}
}

// Splits a string by the delimiter.
func (p *parameter) Split(s string) []string {
	return strings.Split(s, p.delimiter)
}

// Returns the help of a parameter.
func (p *parameter) GetHelp() string {
	var buffer bytes.Buffer
	buffer.WriteString(strings.Join(p.CliNames(), " "))
	buffer.WriteString(" ")
	buffer.WriteString(p.tipe.String())
	buffer.WriteString(" ")
	if p.IsArrayType() {
		buffer.WriteString("delimiter ")
		if p.delimiter == " " {
			buffer.WriteString("whitespace ")

		} else {
			buffer.WriteString(p.delimiter)
			buffer.WriteString(" ")
		}
	}

	parenthesis := p.mandatory || p.defaultValue != "" || p.envKey != ""
	if parenthesis {
		buffer.WriteString("(")
	}
	infos := make([]string, 0)
	if p.mandatory {
		infos = append(infos, "mandatory")
	}
	if p.defaultValue != "" {
		v := fmt.Sprint("default=", p.defaultValue)
		infos = append(infos, v)
	}
	if p.envKey != "" {
		envValue := os.Getenv(p.envKey)
		v := fmt.Sprint("env={key:", p.envKey, ",value:", envValue, "}")
		infos = append(infos, v)
	}
	if parenthesis {
		buffer.WriteString(strings.Join(infos, ";"))
		buffer.WriteString(")")
	}
	if p.description != "" {
		buffer.WriteString(": ")
		buffer.WriteString(p.description)
	}
	return buffer.String()
}

// Returns if a shortName has been defined.
func (p *parameter) hasShortName() bool {
	return p.shortName != ""
}

// Returns if the parameter matches the string.
func (p *parameter) Matches(s string) bool {
	for _, name := range p.CliNames() {
		if name == s {
			return true
		}
	}
	return false
}

func (p *parameter) IsArrayType() bool {
	stringArrayType, intArrayType := reflect.TypeOf([]string{}), reflect.TypeOf([]int{})
	t := p.tipe
	return t == stringArrayType || t == intArrayType
}

// Gets value of the object by reflect
func (p *parameter) getValue(obj interface{}) reflect.Value {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}
	fieldValue := objValue.FieldByName(p.name)
	return fieldValue
}

// Sets
func (p *parameter) setBool(target reflect.Value) func(value string) error {
	target.SetBool(true)
	return nil
}

func (p *parameter) setInt(target reflect.Value) func(value string) error {
	return func(value string) error {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		target.SetInt(int64(intValue))
		return nil
	}
}
func (p *parameter) setString(target reflect.Value) func(value string) error {
	return func(value string) error {
		target.SetString(value)
		return nil
	}
}
func (p *parameter) setStringArray(target reflect.Value) func(value string) error {
	return func(value string) error {
		parts := p.Split(value)
		target.Set(reflect.ValueOf(parts))
		return nil
	}
}
func (p *parameter) setIntArray(target reflect.Value) func(value string) error {
	return func(value string) error {
		parts := p.Split(value)
		intParts := []int{}
		for _, i := range parts {
			j, err := strconv.Atoi(i)
			if err != nil {
				return err
			}
			intParts = append(intParts, j)
		}
		target.Set(reflect.ValueOf(intParts))
		return nil
	}
}

func (p *parameter) setterOnValue(target reflect.Value) func(value string) error {
	switch p.tipe {
	case reflect.TypeOf(true):
		return p.setBool(target)
	case reflect.TypeOf(1):
		return p.setInt(target)
	case reflect.TypeOf(""):
		return p.setString(target)
	case reflect.TypeOf([]string{}):
		return p.setStringArray(target)
	case reflect.TypeOf([]int{}):
		return p.setIntArray(target)
	}
	return nil
}

// fills an object with the desired value
func (p *parameter) SetterCallback(obj interface{}) (func(value string) error, error) {
	if p.used {
		return nil, fmt.Errorf("%s used multiple times", p.name)
	}
	p.used = true
	target := p.getValue(obj)
	setter := p.setterOnValue(target)
	// no setter callback for bool type
	if setter == nil && p.tipe != reflect.TypeOf(true) {
		return nil, fmt.Errorf("Incompatible type")
	}
	return setter, nil
}

func (p *parameter) setDefault(value reflect.Value) error {
	exists, err := p.setDefaultFromEnv(value)
	if exists && err != nil {
		return err
	} else if exists {
		return nil
	}
	if p.defaultValue != "" {
		setter := p.setterOnValue(value)
		return setter(p.defaultValue)
	}
	return nil
}

func (p *parameter) setDefaultFromEnv(value reflect.Value) (exists bool, err error) {
	setter := p.setterOnValue(value)
	envValue := os.Getenv(p.envKey)
	if envValue != "" {
		err := setter(envValue)
		return true, err
	}
	return false, nil
}

func (p *parameter) testDefaultValue() error {
	setMockValue := func(value interface{}) error {
		valueType := reflect.TypeOf(value)
		mockValue := reflect.New(valueType).Elem()
		return p.setDefault(mockValue)
	}
	switch p.tipe {
	case reflect.TypeOf(false):
		return setMockValue(false)
	case reflect.TypeOf(1):
		return setMockValue(1)
	case reflect.TypeOf(""):
		return setMockValue("")
	case reflect.TypeOf([]string{}):
		return setMockValue([]string{})
	case reflect.TypeOf([]int{}):
		return setMockValue([]int{})
	}
	return fmt.Errorf("Incompatible type")
}
func (p *parameter) validate() error {
	getError := func(s string) error {
		return fmt.Errorf("parameter %s : %s",
			p.name, s,
		)
	}
	if (p.mandatory || p.tipe == reflect.TypeOf(true)) && (p.defaultValue != "" || p.envKey != "") {
		return getError("can not be mandatory or have a default value")
	} else if !p.IsArrayType() && strings.Trim(p.delimiter, " ") != "" {
		return getError("delimiter on non array type")
	} else if p.mandatory && p.tipe == reflect.TypeOf(true) {
		return getError("boolean type can not be mandatory")
	}
	return p.testDefaultValue()
}

// Changes the parameter by the value of the constraint.
func (p *parameter) fillParameter(constraint string) error {
	splittedConstraint, err := splitConstraint(constraint)
	key, value := splittedConstraint.key, splittedConstraint.value
	if err != nil {
		return err
	}
	switch key {
	case "description":
		p.description = value
		return nil
	case "shortname":
		p.shortName = value
		return nil
	case "mandatory":
		p.mandatory = true
		return nil
	case "default":
		p.defaultValue = value
		return nil
	case "env":
		p.envKey = value
		return nil
	case "delimiter":
		p.delimiter = value
		return nil
	}
	return fmt.Errorf("unknown key %s", splittedConstraint.value)
}

// Returns a new Parameter from the structField
func newParameter(sf reflect.StructField) (*parameter, error) {
	tag, newParam := sf.Tag.Get(tagName), parameter{
		name:  sf.Name,
		index: sf.Index[0],
		tipe:  sf.Type,
	}
	if tag == "omit" {
		return nil, nil
	}
	if newParam.IsArrayType() && newParam.delimiter == "" {
		newParam.delimiter = constraintsDelimiter
	}
	if tag == "" {
		return &newParam, nil
	}
	constraints := strings.Split(tag, constraintsDelimiter)
	for _, constraint := range constraints {
		err := newParam.fillParameter(constraint)
		if err != nil {
			return nil, fmt.Errorf(
				"error parsing constraint %s at field %s : %e",
				constraint, newParam.name, err)
		}
	}
	if err := newParam.validate(); err != nil {
		return nil, err
	}
	return &newParam, nil
}
