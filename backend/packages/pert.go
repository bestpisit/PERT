package pert

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Activity struct {
	Code         string     `json:"Code"`
	Desc         string     `json:"Desc"`
	Edges        []Edge     `json:"Edges"`
	Dependencies []Activity `json:"Dependencies"`
	Dependents   []string   `json:"Dependents"`
	Duration     int        `json:"Duration"`
}

type Edge struct {
	Src  string `json:"Src"`
	Dest string `json:"Dest"`
}

func PertHandler(c *gin.Context) {
	var activities []Activity
	if err := c.ShouldBindJSON(&activities); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	for key := range endActivities {
		newEdge := Edge{key, "End"}
		edges = append(edges, newEdge)
		activityNetwork["End"].Dependents = append(activityNetwork["End"].Dependents, key)
		activityNetwork[key].Edges = append(activityNetwork[key].Edges, newEdge)
	}
	fmt.Println(edges)

	pertData := PERT(activityNetwork)
	jsonPERT, errr := json.Marshal(pertData)
	if errr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activities processed", "data": pertData, "json": jsonPERT})
}

type activityCase struct {
	activity *Activity
	ES       int
	Duration int
	EF       int
	LS       int
	Slack    int
	LF       int
}

func PERT(activityNetwork map[string]*Activity) map[string]*activityCase {
	queueCase := list.New()
	startCase := &activityCase{activityNetwork["Start"], -1, activityNetwork["Start"].Duration, -1, -1, -1, -1}
	queueCase.PushBack(startCase)
	activityCases := make(map[string]*activityCase)
	for queueCase.Len() > 0 {
		element := queueCase.Front()
		currentCase := element.Value.(*activityCase)
		queueCase.Remove(element)
		existCase, exist := activityCases[currentCase.activity.Code]
		if exist && currentCase != existCase {
			continue
		}
		if !exist {
			activityCases[currentCase.activity.Code] = currentCase
		} else {
			currentCase = activityCases[currentCase.activity.Code]
		}
		//check parent
		pass := true
		if len(currentCase.activity.Dependents) > 0 {
			maxEF := 0
			for _, dependent := range currentCase.activity.Dependents {
				parent, exist := activityCases[dependent]
				if !exist || parent.EF == -1 {
					pass = false
					break
				}
				if maxEF < parent.EF {
					maxEF = parent.EF
				}
			}
			currentCase.ES = maxEF
		} else { // means it start
			currentCase.ES = 0
		}
		if pass { //branching
			currentCase.EF = currentCase.ES + currentCase.Duration
			fmt.Println(currentCase.activity.Code, currentCase.activity.Dependents, currentCase.ES, currentCase.Duration, currentCase.EF)
			for _, edge := range currentCase.activity.Edges {
				newCase := &activityCase{activityNetwork[edge.Dest], -1, activityNetwork[edge.Dest].Duration, -1, -1, -1, -1}
				queueCase.PushBack(newCase)
			}
		} else { //revert
			queueCase.PushBack(currentCase)
		}
	}
	return activityCases
}
