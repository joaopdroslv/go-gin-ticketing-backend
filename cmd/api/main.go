package main

import (
	"log"
	"ticket-io/internal/config"
	"ticket-io/internal/database"
	userHandler "ticket-io/internal/user/handler"
	userRepository "ticket-io/internal/user/repository"
	userService "ticket-io/internal/user/service"

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

	repo := userRepository.NewMySQLUserRepository(db)
	srvc := userService.NewUserService(repo)
	userHandler.RegisterRoutes(api, srvc)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})

	r.Run(":8080")
}
