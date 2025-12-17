package main

import (
	"log"
	"ticket-io/internal/config"
	"ticket-io/internal/database"
	userFeatHandler "ticket-io/internal/user/handler"
	userFeatRepository "ticket-io/internal/user/repository"
	userFeatService "ticket-io/internal/user/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()

	db, err := database.NewMysql(cfg.DockerDatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	api := r.Group("/api/v1")

	userRepo := userFeatRepository.NewMySQLUserRepository(db)
	userStatusRepo := userFeatRepository.NewMySQLUserStatusRepository(db)
	userStatusSrvc := userFeatService.NewUserStatusService(userStatusRepo)
	userSrvc := userFeatService.NewUserService(userRepo, userStatusSrvc)
	userFeatHandler.RegisterRoutes(api, userSrvc)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})

	r.Run(":8080")
}
