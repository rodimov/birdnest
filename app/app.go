package app

import (
	dronesDelivery "birdnest/internal/drones/delivery"
	dronesRepo "birdnest/internal/drones/repository"
	dronesUC "birdnest/internal/drones/usecase"
	eventManager "birdnest/internal/events/manager"
	pilotsDelivery "birdnest/internal/pilots/delivery"
	pilotsRepo "birdnest/internal/pilots/repository"
	pilotsUC "birdnest/internal/pilots/usecase"
	schedulerUC "birdnest/internal/scheduler/usecase"
	"birdnest/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	migrate "github.com/rubenv/sql-migrate"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"net/http"
)

type App struct {
	DbStore *gorm.DB
	Config  *Config
	Router  *gin.Engine
}

func NewApp(cfg *Config) (*App, error) {
	router := gin.Default()
	healthChecker(router)
	logger.NewLogger(cfg.LogStorage)
	dbInstance, err := Connect(cfg)

	if err != nil {
		logger.AppLogger.Fatal("error while connection to the DB:", err)
	}

	logger.AppLogger.Infof("Successfully connected to database.")

	sqlDb, err := dbInstance.DB()
	n, err := migrate.Exec(sqlDb, cfg.Dialect, migrate.FileMigrationSource{Dir: "migrations"}, migrate.Up)
	if err != nil {
		logger.AppLogger.Fatal("error while DB migration: ", err)
	}
	logger.AppLogger.Infof("Applied %d migrations!\n", n)

	dronesRepository := dronesRepo.DronesRepository(dbInstance)
	dronesUseCase := dronesUC.DronesUseCase(dronesRepository)
	pilotsRepository := pilotsRepo.NewPilotsRepository(dbInstance)
	pilotsUseCase := pilotsUC.NewUseCase(pilotsRepository, dronesUseCase)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	dronesDelivery.DronesRegisterEndPoints(dronesUseCase, router)
	pilotsDelivery.RegisterEndpoints(pilotsUseCase, router)

	eventsChan := eventManager.StartEventServer(router)

	schedulerUseCase := schedulerUC.NewUseCase(dronesUseCase, pilotsUseCase,
		10, "https://assignments.reaktor.com/birdnest/drones",
		"https://assignments.reaktor.com/birdnest/pilots/", eventsChan)

	_ = schedulerUseCase.StartScheduler()

	return &App{
		DbStore: dbInstance,
		Config:  cfg,
		Router:  router,
	}, nil
}

func StartRouter(app *App) {
	app.Router.Run(fmt.Sprintf(":%d", app.Config.HostPort))
}

// healthChecker godoc
// @Summary      Check app status
// @Description  Check app status
// @Tags         Health
// @Accept       json
// @Produce      json
// @Router       /health [get]
func healthChecker(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
}
