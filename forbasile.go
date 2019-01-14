package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"
	"reflect"
	"strconv"
)

type SUM interface {
	SUMIT(x int, y int) int
}

func main() {
	arg := os.Args[1]
	// module to load
	mod := fmt.Sprintf("%s%s%s%s%s", "./", arg, "/", arg, ".so")
	fmt.Printf(mod)
	os.Mkdir("."+string(filepath.Separator)+os.Args[1], 0777)
	filename := fmt.Sprintf("%s/%s.go", os.Args[1], os.Args[1])
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	strprg := fmt.Sprintf("package main\ntype %s string\nfunc(s %s) SUMIT (x int, y int) int { %s}\nvar %s %s", strings.ToLower(os.Args[1]), strings.ToLower(os.Args[1]), os.Args[2], strings.Title(os.Args[1]), strings.ToLower(os.Args[1]))

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

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", "/home/developer/proj/gitlab.com/forbasile/SUM/SUM.so", "/home/developer/proj/gitlab.com/forbasile/SUM/SUM.go")

	out, err2 := cmd.Output()
	fmt.Println(out)

	if err2 != nil {
		fmt.Println(err2)
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
	// in this case, variable os.Args[1]
	symX, err := plug.Lookup(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 3. Assert that loaded symbol is of a desired type
	// in this case interface type X (defined above)
	var sum SUM
	sum, ok := symX.(SUM)
	if !ok {
		fmt.Println(fmt.Sprintf("unexpected type from module symbol %s", reflect.TypeOf(symX.(SUM))))
		os.Exit(1)
	}

	// 4. use the module
	x1, err := strconv.Atoi(os.Args[3])
	y1, err := strconv.Atoi(os.Args[4])

	fmt.Println(sum.SUMIT(x1, y1))

}
