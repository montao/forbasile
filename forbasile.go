package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	os.Mkdir("."+string(filepath.Separator)+"sumsquare", 0777)
	f, err := os.Create("sumsquare/sumsquare.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	program := "package main\nimport \"fmt\"\ntype sumsquare string\nfunc (s sumsquare) Sumsquare(a int, b int) { fmt.Println(\"HI\")}\nvar Sumsquare sumsquare"
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

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", "/home/developer/proj/gitlab.com/forbasile/sumsquare/sumsquare.so", "/home/developer/proj/gitlab.com/forbasile/sumsquare/sumsquare.go")
	//cmd.Dir = "/home/developer/proj/gitlab.com/forbasile"
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
