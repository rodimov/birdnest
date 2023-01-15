package delivery

import (
	"birdnest/internal/pilots"
	"github.com/gin-gonic/gin"
)

func RegisterEndpoints(pilotsUC pilots.UseCase, router *gin.Engine) {
	handler := NewHandler(pilotsUC)
	api := router.Group("/api")
	{
		pilotsGroup := api.Group("/pilots")
		{
			pilotsGroup.GET("", handler.GetAll)
			pilotsGroup.GET("/:id", handler.GetById)
		}
	}
}
