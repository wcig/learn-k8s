package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println(">> app run")
	router := gin.Default()
	router.GET("/ready", readyHandler)
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

func readyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
