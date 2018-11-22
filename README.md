# YAGCLIF (Yet Another Go Command Line Framework)

## Usage :
### To Parse command line arguments :

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
    any struct field marked as mandatory will cause an error if missing in arguments 
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
