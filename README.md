# YAGCLIF (Yet Another Go Command Line Framework)

## Usage :
### To Parse command line arguments :
```Go

import "github.com/potatomasterrace/yagclif"

type MyContext struct{

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