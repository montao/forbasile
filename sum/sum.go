package main

import "fmt"

type sum string

func Add(a int, b int) {
        sum := a + b
	fmt.Println("Sum: ", sum)
}

// exported
var Sum sum
