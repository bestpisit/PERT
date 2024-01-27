package main

import (
	"github.com/gin-gonic/gin"
	"pert/packages"
)

func main() {
	r := gin.Default()
	r.POST("/pert", pert.PertHandler)
	r.Run(":8080")
}