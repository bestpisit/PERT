package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Activity struct {
	Code         string
	Desc         string
	Edges        []Edge
	Dependencies []Activity
	Dependents   []string
	Duration     int
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
	activityNetwork["Start"] = &Activity{Code: "Start", Dependents: make([]string, 0)}
	activityNetwork["End"] = &Activity{Code: "End", Dependents: make([]string, 0)}
	leftActivities := make(map[string]bool)
	endActivities := make(map[string]bool)

	for i, activity := range activities {
		activityNetwork[activity.Code] = &activities[i]
		leftActivities[activity.Code] = false
		endActivities[activity.Code] = false
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
		for key, val := range leftActivities {
			if !val {
				for _, dependent := range activityNetwork[key].Dependents {
					newEdge := Edge{dependent, key}
					edges = append(edges, newEdge)
					activityNetwork[dependent].Edges = append(activityNetwork[dependent].Edges, newEdge)
					_, exists := endActivities[dependent]
					if exists {
						delete(endActivities, dependent)
					}
				}
				if len(activityNetwork[key].Dependents) == 0 {
					newEdge := Edge{"Start", key}
					edges = append(edges, newEdge)
					activityNetwork["Start"].Edges = append(activityNetwork["Start"].Edges, newEdge)
				}
				delete(leftActivities, key)
			}
		}
		if len(leftActivities) <= 0 {
			break
		}
	}
	for key, _ := range endActivities {
		newEdge := Edge{key, "End"}
		edges = append(edges, newEdge)
		activityNetwork[key].Edges = append(activityNetwork[key].Edges, newEdge)
	}
	fmt.Println(edges)
	PERT(activityNetwork)
}

type ActivityCase struct {
	activity *Activity
	ES       int
	Duration int
	EF       int
	LS       int
	Slack    int
	LF       int
}

func PERT(activityNetwork map[string]*Activity) {
	queueCase := list.New()
	startCase := &ActivityCase{activityNetwork["Start"], -1, activityNetwork["Start"].Duration, -1, -1, -1, -1}
	queueCase.PushBack(startCase)
	// startCase.activity.Dependents = append(startCase.activity.Dependents, "A")
	activityCases := make(map[string]*ActivityCase)
	for queueCase.Len() > 0 {
		element := queueCase.Front()
		currentCase := element.Value.(*ActivityCase)
		queueCase.Remove(element)
		_, exist := activityCases[currentCase.activity.Code]
		if !exist {
			activityCases[currentCase.activity.Code] = currentCase
		} else {
			currentCase = activityCases[currentCase.activity.Code]
		}
		//check parent
		pass := true
		if len(currentCase.activity.Dependents) > 0 {
			for _, dependent := range currentCase.activity.Dependents {
				parent, exist := activityCases[dependent]
				if !exist || parent.EF == -1{
					pass = false
					break
				}
			}
		} else { // means it start
			currentCase.ES = 0
		}

		if pass { //branching
			currentCase.EF = currentCase.ES + currentCase.Duration
		} else { //revert
			queueCase.PushBack(currentCase)
		}
		fmt.Println(currentCase.activity.Code)
	}
}
