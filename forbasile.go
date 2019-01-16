package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"reflect"
	"strconv"
	"strings"
)


//TODO: Investigate how to relax the name FUNCTION into a variable
type Xinterface interface {
	FUNCTION(x int, y int) int
}

func main() {
	//Declare good variable names and create the file for the code
	funame := os.Args[1]
	funamel := strings.ToLower(funame)
	funamet := strings.Title(funame)
	fubody := os.Args[2]
	x1, err := strconv.Atoi(os.Args[3])
	y1, err := strconv.Atoi(os.Args[4])
	filename := fmt.Sprintf("/tmp/%s.go", funame)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	//Here comes the program
	strprg := fmt.Sprintf(`package main 
import (
	"fmt"
)
type %s string 
func(s %s) FUNCTION (x int, y int) int { fmt.Println("")
%s} 
var %s %s`, funamel, funamel, fubody, funamet, funamel)
	fmt.Printf("func(s %s) FUNCTION (x int, y int) int { \n", funamel)
	fmt.Printf("start of %s: x=%d, y=%d\n", funamel, x1, y1)
	l, err := f.WriteString(strprg)
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

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	fmt.Println("compiling plugin")
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", fmt.Sprintf("%s%s%s", "/tmp/", funame, ".so"), fmt.Sprintf("%s%s%s", "/tmp/", funame, ".go"))

	out, err2 := cmd.Output()
	fmt.Println(out)

	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println("loading module")
	// load module
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(fmt.Sprintf("%s%s%s", "/tmp/", funame, ".so"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("looking up symbol")
	// 2. look up a symbol (an exported function or variable)
	// in this case, variable funame
	symX, err := plug.Lookup(funame)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("checking module")
	// 3. Assert that loaded symbol is of a desired type
	// in this case interface type X (defined above)
	var myvar Xinterface
	myvar, ok := symX.(Xinterface)
	if !ok {
		fmt.Println(fmt.Sprintf("unexpected type from module symbol %s", reflect.TypeOf(symX.(Xinterface))))
		os.Exit(1)
	}

	// 4. use the module

	fmt.Println(myvar.FUNCTION(x1, y1))
	fmt.Println(fmt.Sprintf("Generated code: %s", fmt.Sprintf("/tmp/%s%s", funamet , ".go") ))
	fmt.Println(fmt.Sprintf("Generated object file: %s", fmt.Sprintf("/tmp/%s%s", funamet , ".so") ))

}
