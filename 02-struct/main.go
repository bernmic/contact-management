package main

import (
	"fmt"
)

type address struct {
	street string
	zip    string
	city   string
	state  string
}

func main() {
	a := address{"Kennedyallee 62-76", "53175", "Bonn", "NRW"}
	fmt.Printf("%v\n", a)
}
