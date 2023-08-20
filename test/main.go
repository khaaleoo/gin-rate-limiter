package test

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CommonHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello world"})
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	// router.GET("/ping", PingHandler)
	// router.GET("/me", MeHandler)
	return router
}
