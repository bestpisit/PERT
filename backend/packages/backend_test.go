package pert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	msg, err := HelloTest()
	if err != nil {
		t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err,"Nice")
	}
}
func HelloTest() (string, error) {
	fmt.Println("Loading project activities")
	jsonFile, err := os.Open("Exercise.json")
	if err != nil {
		fmt.Println(err)
		return "Error", err
	}
	defer jsonFile.Close()
	fmt.Println("Project activities loaded")
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return "Error", err
	}
	var activities []Activity
	if err := json.Unmarshal(byteValue, &activities); err != nil {
		fmt.Println(err)
		return "Error", err
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
		for key := range endActivities {
			newEdge := Edge{key, "End"}
			edges = append(edges, newEdge)
			activityNetwork["End"].Dependents = append(activityNetwork["End"].Dependents, key)
			activityNetwork[key].Edges = append(activityNetwork[key].Edges, newEdge)
		}
		fmt.Println(edges)
	}
	PERT(activityNetwork)
	return "Success", nil
}