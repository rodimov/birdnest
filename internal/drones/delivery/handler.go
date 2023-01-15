package delivery

import (
	"birdnest/internal/drones"
	"birdnest/internal/error"
	"birdnest/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type dronesHandler struct {
	usecase drones.UseCase
}

func DroneHandler(uc drones.UseCase) *dronesHandler {
	return &dronesHandler{usecase: uc}
}

// GetAll godoc
// @Summary      GetAll drones
// @Description  GetAll drones
// @Tags         Drones
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.DroneDto
// @Failure      400  {object}  error.AppError
// @Failure      404  {object}  error.AppError
// @Failure      405  {object}  error.AppError
// @Failure      500  {object}  error.AppError
// @Router       /drones [get]
func (handler *dronesHandler) GetAll(context *gin.Context) {
	result, err := handler.usecase.GetAll(context.Request.Context())

	if err != nil {
		logger.AppLogger.Info("something went wrong:", err)
		context.JSON(http.StatusBadRequest, error.BadRequest(err))
		return
	} else {
		context.JSON(http.StatusOK, SuccessResponse{Data: ToDronesDto(result)})
	}
}

// GetById godoc
// @Summary      Get drone by id
// @Description  Get drone by id
// @Tags         Drones
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Drone ID"
// @Success      200  {object}  dto.DroneDto
// @Failure      400  {object}  error.AppError
// @Failure      404  {object}  error.AppError
// @Failure      500  {object}  error.AppError
// @Router       /drones/{id} [get]
func (handler *dronesHandler) GetById(context *gin.Context) {
	id := context.Param("id")
	ctx := context.Request.Context()
	result, err := handler.usecase.GetById(ctx, id)

	if err != nil {
		logger.AppLogger.Info("something went wrong:", err)
		context.JSON(http.StatusBadRequest, error.BadRequest(err))
		return
	} else {
		context.JSON(http.StatusOK, SuccessResponse{Data: ToDroneDto(result)})
	}
}
