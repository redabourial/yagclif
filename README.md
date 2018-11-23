# YAGCLIF (Yet Another Go Command Line Framework)
Yagclif is a basic cli arguments parser that can be used as a framework.

## Usage :
### To Parse command line arguments :
#### Code
```Go
package main

import (
	"fmt"

	"github.com/potatomasterrace/yagclif"
)
// Define the struct that will be used as Context/
type MyContext struct {
    MyInteger      int    `yagclif:"shortname:mi;mandatory"`
    // delimiter defaults to ; if none is set;
    MyIntegerArray []int  `yagclif:"delimiter:,`
	MyString       string `yagclif:"default:hello world !;description:short explaination"`
}

func main() {
    // Initiate an instance of the context
    context := MyContext{}
    // Pass it to be parsed
	err := yagclif.Parse(&context)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", context)
	// That's pretty much it 
}
```
#### CLI 
##### ./main -mi 42 
    Hello


### As a Framework :
```Go
package main

import (
    "fmt"
    "github.com/potatomasterrace/yagclif"
)

type MyContext struct{
    MyInteger int `yagclif:"mandatory"`
    MyIntegerArray []int `yagclif:"delimiter:,;default:1,2,3,4,5"`
    MyString string `yagclif:"description:short explaination"`
}

func main(){
    context := MyContext{}
    err := yagclif.Parse(&context)
    if err != nil{
        panic(err)
    }
    fmt.Printf("%#v",context)
    // Context now has been loaded with arguments
}
```
## Supported struct field types:
* boolean
* string 
* int 
* int array
* string array
## Tag options :
### Mandatory
    any struct field marked as mandatory will cause an error if missing in arguments. any
Example
```Go
    MyInteger int `yagclif:"mandatory"`
```
### Mandatory
### Delimiter 
## Help generator
### for cli app :
TODO add sample help for app with example
TODO add sample help for route with example

## Known issues :
### Nested structs do NOT work
    Your parameter can not have nested struct.
### Pointer parameters do NOT work
    you can't have pointers as parameters for your callback functions.
### Autocompletion
    Autocompletion is not available from the cli and is not planned to be added.
