package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
)

func main() {
	//Declare good variable names and create the file for the code
	funame := os.Args[1]
	funamel := strings.ToLower(funame)
	funamet := strings.Title(funame)
	fubody := os.Args[2]
	x1, err := strconv.Atoi(os.Args[3])
	y1, err := strconv.Atoi(os.Args[4])
	filename := fmt.Sprintf("/tmp/code_%s.go", funame)
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
func %s(x int, y int) int { fmt.Println("")
%s} 
`, funamet, fubody)
	fmt.Printf("func %s(x int, y int) int { \n", funamet)
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
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", fmt.Sprintf("%s%s%s", "/tmp/plugin_", funame, ".so"), fmt.Sprintf("%s%s%s", "/tmp/code_", funame, ".go"))

	out, err2 := cmd.Output()
	fmt.Println(out)

	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println("loading module")
	// load module
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(fmt.Sprintf("%s%s%s", "/tmp/plugin_", funame, ".so"))
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

	plugFunc, ok := symX.(func(int, int) int)
	if !ok {
		panic("Plugin has no such function")
	}

	output := plugFunc(x1, y1)
	fmt.Println(output)

	fmt.Println(fmt.Sprintf("Generated code: %s", fmt.Sprintf("/tmp/code_%s%s", funamet , ".go") ))
	fmt.Println(fmt.Sprintf("Generated object file: %s", fmt.Sprintf("/tmp/plugin_%s%s", funamet , ".so") ))

	cmd2 := exec.Command("pmap", strconv.Itoa(os.Getpid()))
	out2, err3 := cmd2.Output()
	fmt.Println(string(out2))

	if err3 != nil {
		fmt.Println(err3)
		return
	}

}
