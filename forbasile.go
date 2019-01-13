package main

import (
	"fmt"
	"os"
	"plugin"
)

type Sum interface {
	Sum(int, int)
}

func main() {
        arg := os.Args[1]
	// module to load
	mod := fmt.Sprintf("%s%s%s%s%s", "./", arg, "/", arg, ".so")
        fmt.Printf(mod)

    f, err := os.Create("sumsquare.go")
    if err != nil {
        fmt.Println(err)
        return
    }
    program := "import \"fmt\"\nfunc Demo() { fmt.Println(\"HI\")}"
    l, err := f.WriteString(program)
    if err != nil {
        fmt.Println(err)
        f.Close()
        return
    }
    fmt.Println(l, "bytes written successfully")
    err = f.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
	

	// load module
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 2. look up a symbol (an exported function or variable)
	// in this case, variable Sum
	symSum, err := plug.Lookup("Sum")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 3. Assert that loaded symbol is of a desired type
	// in this case interface type Sum (defined above)
	var sum Sum
	sum, ok := symSum.(Sum)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}

	// 4. use the module
	sum.Sum(5, 3)

}
