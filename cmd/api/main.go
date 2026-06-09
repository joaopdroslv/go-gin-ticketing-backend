package main

import (
	"context"
	"log"
	"time"

	"go-gin-ticketing-backend/internal/access_control/http"
	"go-gin-ticketing-backend/internal/access_control/repository"
	"go-gin-ticketing-backend/internal/access_control/service"
	"go-gin-ticketing-backend/internal/api"
	"go-gin-ticketing-backend/internal/auth"
	"go-gin-ticketing-backend/internal/config"
	"go-gin-ticketing-backend/internal/infra"
	"go-gin-ticketing-backend/internal/middlewares"
	"go-gin-ticketing-backend/internal/user"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	env := config.NewEnv()
	logger := config.NewLogger()
	_ = logger // Ignore it for now

	db, err := infra.NewMysqlDatabase(env.DockerDatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	// r.Use(gin.Logger(), gin.Recovery())

	jwtMiddleware := middlewares.JWTAuthenticationMiddleware(env.JWTSecret)
	rateLimitMiddleware := middlewares.RateLimitMiddleware(env.RequestsPerMinute, time.Minute)

	// All routes covered by the rate limit middleware
	r.Use(rateLimitMiddleware)

	userRepository := user.NewUserMysqlRepository(db)
	authRepository := auth.NewAuthMysqlRepository(db)
	permissionRepository := acrepository.NewPermissionRepositoryMysql(db)

	ctx := context.Background()
	userService, err := user.NewUserService(ctx, userRepository)
	if err != nil {
		log.Fatal("failed to create the user service")
	}
	authService := auth.NewAuthService(authRepository, env.JWTSecret, env.JWTTTL)
	permissionService := acservice.NewPermissionService(permissionRepository)

	authHandler := auth.NewAuthHandler(authService)
	userHandler := user.NewUserHandler(userService)
	permissionHandler := achttp.NewPermissionHandler(permissionService)

	dependencies := api.Dependencies{
		AuthHandler:       authHandler,
		UserHandler:       userHandler,
		PermissionHandler: permissionHandler,
		JWTMiddleware:     &jwtMiddleware,
		PermissionService: permissionService,
	}

	api.Register(r, dependencies)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})

	r.Run(":" + env.HTTPPort)
}
