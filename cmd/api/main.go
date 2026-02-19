package main

import (
	"log"

	"go-gin-ticketing-backend/internal/config"
	"go-gin-ticketing-backend/internal/database"

	authmiddleware "go-gin-ticketing-backend/internal/auth/middleware"

	authhandler "go-gin-ticketing-backend/internal/auth/handler"
	userhandler "go-gin-ticketing-backend/internal/user/handler"

	authservice "go-gin-ticketing-backend/internal/auth/service"
	userservice "go-gin-ticketing-backend/internal/user/service/user"
	userstatusservice "go-gin-ticketing-backend/internal/user/service/user_status"

	authrepository "go-gin-ticketing-backend/internal/auth/repository/auth"
	permissionrepository "go-gin-ticketing-backend/internal/auth/repository/permission"
	userrepository "go-gin-ticketing-backend/internal/user/repository/user"
	userstatusrepository "go-gin-ticketing-backend/internal/user/repository/user_status"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	env := config.NewEnv()
	logger := config.NewLogger()
	_ = logger // Discarding it for now

	db, err := database.NewMysql(env.DockerDatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	apiV1Group := r.Group("/api/v1")

	// repositories
	userRepo := userrepository.New(db)
	userStatusRepo := userstatusrepository.New(db)
	authRepo := authrepository.New(db)
	permissionRepo := permissionrepository.New(db)

	// services
	userStatusService := userstatusservice.New(userStatusRepo)
	userService := userservice.New(userRepo, userStatusService)
	authService := authservice.New(authRepo, permissionRepo, env.JWTSecret, env.JWTTTL)

	// handlers
	authHandler := authhandler.New(authService)
	userHandler := userhandler.New(userService)

	// middlewares
	jwtMiddleware := authmiddleware.JWTAuthentication(env.JWTSecret)

	// routes

	// auth (public)
	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.POST("/register", authHandler.RegisterUser)
		authGroup.POST("/login", authHandler.LoginUser)
	}

	// users (protected)
	userGroup := apiV1Group.Group("/users")
	userGroup.Use(jwtMiddleware)
	userhandler.RegisterRoutes(userGroup, userHandler, authService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})

	r.Run(":8080")
}
