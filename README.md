# Go Plugin

The code in this repository uses the `plugin` package in Go 1.10 (see https://tip.golang.org/pkg/plugin/).  A Go plugin is package compiled with the `-buildmode=plugin` which creates a shared object (`.so`) library file instead of the standar archive (`.a`) library file.  As you will see here, using the standard library's `plugin` package, Go can dynamically load the shared object file at runtime to access exported elements such as functions an variables.

You can read the related article [on Medium](https://medium.com/learning-the-go-programming-language/writing-modular-go-programs-with-plugins-ec46381ee1a9).

## Requirements
The plugin system requires Go version 1.10.  At this time, it is only supports plugin on Linux.  Attempt 
to compile plugins on OSX, for instance, will result in  `-buildmode=plugin not supported on darwin/amd64` error.

the body (that is, some Go code) of that sum function should be a
program argument. Maybe


    go run forbasile.go SUM 'return x+y' 3 5


and of course


    go run forbasile.go SUMSQUARE 'return x*x + y*y' 3 4


is expected to print 25
and the body could even be some more complex Go statement.

## A Pluggable System
The demo in this repository implements a simple plugin generator.  The plugin package (directory `/tmp/<NAME>`) implements code that prints a sum.  File `./sum.go` uses the new Go `plugin` package to load the pluggable modules and displays the proper message using passed command-line parameters.

For instance, when the program is executed it prints a greeting in English or Chinese 
using the passed parameter to select the plugin to load for the appropriate language.
```
> go run forbasile.go sum
.. calling it with xx=3 yy=5
start of SUM: x=3, y=5
result of SUM (3, 5) is 8
```

As you can see, the capability of the driver program is dynamically expanded by the plugins allowing it to display a greeting message in different language without the need to recompile the program.

Let us see how this is done.


## The Plugin
To create a pluggable package is simple.  Simply create a regular Go package designated as `main`. Use the capitalization rule to indicate functions and variables that are exported as part of the plugin.  This is shown below in file  `./eng/greeter.go`.  This plugin is responsible for displaying a message in `English`.  

File [./sum/sum.go](./sum/sum.go)

```go
package main

import "fmt"

type sum int

func (g sum) Sum() {
	fmt.Println(a+b)
}

// this is exported
var Sum sum
```
Notice a few things about the pluggable module:

- Pluggable packages are basically regular Go packages
- The package must be marked `main`
- The exported variables and functions can be of any type (I found no documented restrictions)

The previous code exports variable `Greeter` of type `greeting`.  As we will see later, the code that will consume this exported value must have a compatible type for assertion.  One way this can be handled is to have an interface with the same method set. (The plugin package in directory `./chi` is exactly the same code except the message is in Chinese.)

### Compiling the Plugins
The plugin package is compiled using the normal Go toolchain.  The only requirement is to use the `buildmode=plugin` compilation flag as shown below:

```
go build -buildmode=plugin -o sum/sum.so sum/sum.go
```
The compilation step will create `./sum/sum.so` plugin file.

### Using the Plugin
Once the plugin module is available, it can be loaded dynamically using the Go standard library's `plugin` package.  Let us examine file [./forbasile.go](./forbasile.go), the driver program that loads and uses the plugin at runtime. Loading and using a shared object library is done in several steps as outlined below:

#### 1. Import package plugin
```
import (
	...
	"plugin"
)
```
#### 2. Define/select type for imported elements (optional)
The exported elements, from the pluggable package, can be of any type.  The consumer code, loading the plugin, must have a compatible type defined (or pre-defined in case of built-in types) for assertion. In this example we define interface type `Greeter` as a type that will be asserted against the exported variable from the plugin module. 
```
type Greeter interface {
	Greet()
}
```
#### 3. Determine the .so file to load
The `.so` file must be in a location accessible from you program in order to open it.  In this example, the file .so files are located in directories `./eng` and `./chi`.  and are selected based on the value of a command-line argument.  The selected name is then assigned to variable `mod`.
```
func main() {
	//  module to load
    mod = "./sum/sum.so"
	}
...
```
#### 4. Open the plugin package
Using the Go standard library's `plugin` package, we can now `open` the plugin module.  That step creates a value of type `*plugin.Plugin`.  It is used later to manage access to the plugin's exported elements.

```
func main(){
...
	// load module
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
...
```
#### 6. Lookup a Symbol
Next, we use the `*plugin.Plugin` value to search for symbols that matches the name of the exported elements from the plugin module.  In our example plugin ([./eng/greeter.go](./eng/greeter.go), seen earlier), we exported a variable called `Greeter`.  Therefore, we use `plug.Lookup("Greeter")` to locate that symbol.  The loaded symbol is then assigned to variable `symGreeter` (of type `package.Symbol`).
```
func main(){
...
	// 2. look up a symbol (an exported function or variable)
	// in this case, variable Greeter
	symGreeter, err := plug.Lookup("Greeter")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
...
```

#### 7. Assert the symbol's type and use it
Once we have the symbol loaded, we still have one additional step before we can use it.  We must use type assertion to validate that the symbol is of an expected type and (optionally) assign its value to a variable of that type.  In this example, we assert symbol `symGreeter` to be of interface type `Greeter` with `symGreeter.(Greeter)`.  Since the exported symbol from the plugin module `./eng/eng.so` is a variable with method `Greet` attached, the assertion is true and the value is assigned to variable `greeter`.  Lastly, we invoke the method from the plugin module with `greeter.Greet()`.
```
func main(){
...
	// 3. Assert that loaded symbol is of a desired type
	// in this case interface type Greeter (defined above)
	var sum Sum
	greeter, ok := symGreeter.(Greeter)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}

	// 4. use the module
	sum.Add(3,5)

}
```
