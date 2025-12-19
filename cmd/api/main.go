package main

import (
	"log"

	"ticket-io/internal/config"
	"ticket-io/internal/database"

	authhandler "ticket-io/internal/auth/handler"
	userhandler "ticket-io/internal/user/handler"

	authservice "ticket-io/internal/auth/service"
	statusservice "ticket-io/internal/user/service/status"
	userservice "ticket-io/internal/user/service/user"

	authrepository "ticket-io/internal/auth/repository"
	statusrepository "ticket-io/internal/user/repository/status"
	userrepository "ticket-io/internal/user/repository/user"

	authmiddleware "ticket-io/internal/auth/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()

	db, err := database.NewMysql(cfg.DockerDatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	apiV1Group := r.Group("/api/v1")

	// repositories
	userRepo := userrepository.New(db)
	statusRepo := statusrepository.New(db)
	authRepo := authrepository.New(db)

	// services
	statusService := statusservice.New(statusRepo)
	userService := userservice.New(userRepo, statusService)
	authService := authservice.New(authRepo, cfg.JWTSecret)

	// handlers
	authHandler := authhandler.New(authService)
	userHandler := userhandler.New(userService)

	// middlewares
	jwtMiddleware := authmiddleware.JWTAuthentication(cfg.JWTSecret)

	// routes

	// auth (public/no middleware)
	auth := apiV1Group.Group("/auth")
	{
		auth.POST("/register", authHandler.RegisterUser)
		auth.POST("/login", authHandler.LoginUser)
	}

	// users (protected)
	userGroup := apiV1Group.Group("/users")
	userGroup.Use(jwtMiddleware)
	userhandler.RegisterRoutes(userGroup, userHandler)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})

	r.Run(":8080")
}
