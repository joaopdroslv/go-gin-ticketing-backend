package main

import (
	"learning-gin/internal/config"
	"learning-gin/internal/database"
	"learning-gin/internal/user"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()

	db, err := database.NewMysql(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	api := r.Group("/api/v1")

	userRepository := user.NewMySQLRepository(db)
	userService := user.NewService(userRepository)
	user.RegisterRoutes(api, userService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})

	r.Run(":8080")
}
