package delivery

import (
	"birdnest/internal/drones"
	"github.com/gin-gonic/gin"
)

func DronesRegisterEndPoints(drones drones.UseCase, router *gin.Engine) {
	handler := DroneHandler(drones)
	g := router.Group("/api/drones")
	g.GET("", handler.GetAll)
	g.GET("/:id", handler.GetById)

}
