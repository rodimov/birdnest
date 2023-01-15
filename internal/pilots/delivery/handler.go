package delivery

import (
	"birdnest/internal/error"
	"birdnest/internal/pilots"
	"birdnest/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type pilotsHandler struct {
	usecase pilots.UseCase
}

func NewHandler(uc pilots.UseCase) *pilotsHandler {
	return &pilotsHandler{usecase: uc}
}

// GetAll godoc
// @Summary      GetAll pilots
// @Description  GetAll pilots
// @Tags         Pilots
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.PilotDto
// @Failure      400  {object}  error.AppError
// @Failure      404  {object}  error.AppError
// @Failure      405  {object}  error.AppError
// @Failure      500  {object}  error.AppError
// @Router       /pilots [get]
func (handler *pilotsHandler) GetAll(context *gin.Context) {
	result, err := handler.usecase.GetAll(context.Request.Context())

	if err != nil {
		logger.AppLogger.Info("something went wrong:", err)
		context.JSON(http.StatusBadRequest, error.BadRequest(err))
		return
	} else {
		context.JSON(http.StatusOK, SuccessResponse{Data: ToPilotsDto(result)})
	}
}

// GetById godoc
// @Summary      Get pilot by drone id
// @Description  Get the pilot by drone id
// @Tags         Pilots
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Drone ID"
// @Success      200  {object}  dto.PilotDto
// @Failure      400  {object}  error.AppError
// @Failure      404  {object}  error.AppError
// @Failure      500  {object}  error.AppError
// @Router       /pilots/{id} [get]
func (handler *pilotsHandler) GetById(context *gin.Context) {
	id := context.Param("id")
	ctx := context.Request.Context()
	pilot, err := handler.usecase.GetByDroneId(ctx, id)

	if err != nil {
		logger.AppLogger.Info("something wrong:", err)
		context.JSON(http.StatusNotFound, error.EntityNotFound(err))
	} else {
		context.JSON(http.StatusOK, SuccessResponse{Data: ToPilotDto(pilot)})
	}
}
