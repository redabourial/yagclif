# YAGCLIF (Yet Another Go Command Line Framework)
Yagclif is a basic cli arguments parser that can be used as a framework.

## Usage :
### To Parse command line arguments :
#### Code 
```Go
package main

import (
	"fmt"
	"os"

	"github.com/potatomasterrace/yagclif"
)

// Define the struct that will be used as Context/
type MyContext struct {
	// Struct field name are converted to lower case
	// usage in cli is --{{fieldname}} and -{{shortname}}
	// if present.
	MyInteger int `yagclif:"shortname:mi;mandatory"`
	// delimiter defaults to ; if none is set;
	MyIntegerArray []int `yagclif:"delimiter:,`
	// default sets the value of the field in advance
	// description show in the help and errors texts
	MyString string `yagclif:"default:hello world !;description:short explaination"`
}

func main() {
	// Initiate an instance of the context
	context := MyContext{}
	// Pass it to be parsed
	remainingArgs, err := yagclif.Parse(&context)
	if err != nil {
		// Output error message followed by help
		fmt.Print(err)
		// Exit with an error code
		os.Exit(1)
	}
	// That's pretty much it. just outputing return values here.
	fmt.Printf(
		"Context %#v\r\nRemaining args : %#v\r\n",
		context, remainingArgs,
	)
}

```
#### CLI 
##### go run main.go 
will output an error as myinteger is mandatory and missing

    missing argument [--myinteger -mi] for MyInteger
    usage:
    --myinteger -mi int (mandatory)
    --myintegerarray []int delimiter ;
    --mystring string (default = hello world !): short explaination
    exit status 1
##### go run main.go -mi 42 anExtraArgument --mystring helloWorld anotherExtraArgument
    Context main.MyContext{MyInteger:42, MyIntegerArray:[]int(nil), MyString:"helloWorld"}
    Remaining args : []string{"anExtraArgument", "anotherExtraArgument"}
### To generate help text for context :
#### Code
```Go
    var helpText string = yagclif.GetHelp(&context)
```
#### Example output
    --myinteger -mi int (mandatory)
    --myintegerarray []int delimiter ;
    --mystring string (default = hello world !): short explaination
### As a Framework :
#### Code 
```Go
package main

import (
	"fmt"
	"os"

	"github.com/potatomasterrace/yagclif"
)

type MyContext struct {
	MyInteger      int    `yagclif:"shortname:mi;mandatory"`
	MyIntegerArray []int  `yagclif:"delimiter:,`
	MyString       string `yagclif:"default:hello world !;description:short explaination"`
}

func panicIfError(err error) {
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
func main() {
	app := yagclif.NewCliApp("My cool project name", "a cool description for my project")
	// supports two types of functions:
	// func([]string) for remaining arguments
	// first parameter struct that will be parsed then a []string parameter for remaining arguments
	err := app.AddRoute("actionA", "output the parsed context and the arguments",
		// Type *MyContext works too
		func(context MyContext, remainingArgs []string) {
			fmt.Println("you choose ActionA")
			fmt.Println("context:", context, "remaingArgs:", remainingArgs)
		})
	panicIfError(err)
	err = app.AddRoute("actionB", "output remaining arguments",
		func(remainingArgs []string) {
			fmt.Println("you choose ActionB")
			fmt.Println("remaingArgs:", remainingArgs)
		})
	panicIfError(err)
	// Output help on error
	err = app.Run(true)
	panicIfError(err)
}

```
#### CLI
##### go run main.go actionA someArguments...
will output an error as myinteger is mandatory and missing

panic: missing argument [--myinteger -mi] for MyInteger
My cool project name
a cool description for my project

         actionA : output the parsed context and the arguments
                 usage :
                        --myinteger -mi int (mandatory)
                        --myintegerarray []int delimiter ;
                        --mystring string (default = hello world !): short explaination

         actionB : output remaining arguments
##### go run main.go actionA -mi 42 foo bar
you choose ActionA
context: {42 [] hello world !} remaingArgs: [foo bar]
##### go run main.go actionB -mi 42 foo bar
you choose ActionB
[-mi 42 foo bar]
## Supported struct field types:
* boolean
* string 
* int 
* []int
* []string
## Tag options :
### ShortName
    Struct field can have a shortname for usage in the cli. 
    shortname will be preceeded by a hyphen (-). the name will be preceeded by two hyphens (--).
    both names and shortnames are converted to lower case.
```Go
    MyInteger int `yagclif:"shortname:somename"`
```
### Mandatory
    Any struct field marked as mandatory will cause an error if missing in arguments.
Example
```Go
    MyInteger int `yagclif:"mandatory"`
```
### Delimiter 
    a delimiter can be set for the fields with type []string []int.
    If none is set the delimiter is ;
```Go
    MyIntegerArray []int `yagclif:"delimiter:,"`
```
### Default
    a default value for the parameter if missing.
```Go
    MyIntegerArray []int `yagclif:"delimiter:,;default:1,2,3"`
```
### Description
    a description to be printed for the variable

## Known issues :
### Nested structs do NOT work
    Your parameter can not have nested struct. use inheritance instead
### Autocompletion
    Autocompletion is not available from the cli and is not planned to be added.
