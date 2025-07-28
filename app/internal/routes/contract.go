package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netosts/goledger-challenge-besu/internal/handlers"
)

func SetupContractRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	router.POST("/set", handler.SetValue)
	router.GET("/get", handler.GetValue)
	router.POST("/sync", handler.SyncValue)
	router.GET("/check", handler.CheckValue)
}
