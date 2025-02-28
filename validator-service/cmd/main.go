package main

import (
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
	"validator-service/internal/handlers"
	"validator-service/internal/middlewares"
	"validator-service/internal/models"
	"validator-service/internal/monitoring"
	"validator-service/internal/routers"
)

func main() {
	monitoring.InitPrometheus()

	var err error
	db, err := gorm.Open(sqlite.Open("validators.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(
		&models.ValidatorRequest{},
		&models.ValidatorKey{},
	)
	handler := handlers.CreateNewHandler(db)

	ginEngine := gin.Default()
	ginEngine.Use(gin.Logger())
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(middlewares.PrometheusMiddleware())

	routers.SetupRoutes(ginEngine, handler)

	err = ginEngine.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
