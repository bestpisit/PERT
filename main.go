package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Activitie struct {
	Code string
	Desc string
	ES int
	Duration int
	EF int
	LS int
	Slack int
	LF int
	Edges []Edge
	Dependencies []Activitie
	Dependents []string
}

type Edge struct {
	Src Activitie
	Dest Activitie
}

func main() {
	fmt.Println("Loading project activities")
	jsonFile, err := os.Open("BuildingAHouse.json")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("Project activities Loaded")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var activities []Activitie
	json.Unmarshal(byteValue,&activities)
	
}