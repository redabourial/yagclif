# YAGCLIF (Yet Another Go Command Line Framework)

## Usage :
### To Parse command line arguments :
```Go
import "github.com/potatomasterrace/yagclif"

// make sure
type MyContext struct{
    MyInteger `yaglif:"mandatory"`
}


func main(){
    context := MyContext{}
    err := yagclif.Parse(&context)
    if err != nil{
        panic(err)
    }
    // Context now has been loaded with arguments
}
```

### As a Framework :
```Go

```
## Supported Struct Tags:

## Tag options :
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
