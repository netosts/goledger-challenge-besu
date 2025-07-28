package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/netosts/goledger-challenge-besu/internal/handlers"
)

func SetupRoutes(handler *handlers.Handler) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		SetupContractRoutes(v1, handler)
	}

	return router
}
