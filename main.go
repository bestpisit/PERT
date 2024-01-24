package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Activity struct {
	Code         string
	Desc         string
	ES           int
	Duration     int
	EF           int
	LS           int
	Slack        int
	LF           int
	Edges        []Edge
	Dependencies []Activity
	Dependents   []string
}

type Edge struct {
	Src  string
	Dest string
}

func main() {
	fmt.Println("Loading project activities")
	jsonFile, err := os.Open("BuildingAHouse.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	fmt.Println("Project activities loaded")

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	var activities []Activity
	if err := json.Unmarshal(byteValue, &activities); err != nil {
		fmt.Println(err)
		return
	}

	activityNetwork := make(map[string]*Activity)
	leftActivities := make(map[string]bool)

	for i, activity := range activities {
		activityNetwork[activity.Code] = &activities[i]
		leftActivities[activity.Code] = false
	}

	for _, activity := range activities {
		for _, dependent := range activity.Dependents {
			dependentActivity, exists := activityNetwork[dependent]
			if exists {
				dependentActivity.Dependencies = append(dependentActivity.Dependencies, activity)
			}
		}
	}

	var edges []Edge

	for {
		for key := range leftActivities {
			leftActivities[key] = false
		}
		for key := range leftActivities {
			for _, dependent := range activityNetwork[key].Dependents {
				_, exists := leftActivities[dependent]
				if exists {
					leftActivities[key] = true
				}
			}
		}
		for key,val := range leftActivities {
			if !val {
				for _, dependent := range activityNetwork[key].Dependents {
					edges = append(edges,Edge{dependent,key})
					fmt.Println(edges)
				}
				delete(leftActivities,key)
			}
		}
		if len(leftActivities) <= 0 {
			break
		}
	}
}
