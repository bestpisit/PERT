package main

import (
    "fmt"
)

type PERT struct {
	Text string
	ES int
	Time int
	EF int
	LS int
	Slack int
	LF int
	Edges []Edge
	Dependencies []PERT
}

type Edge struct {
	Src PERT
	Dest PERT
}

func main() {
	fmt.Print("Hello")
}