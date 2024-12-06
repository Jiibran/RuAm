package main

import (
	"RuAm/internal/pkg/config"
	"RuAm/internal/pkg/db"
	"RuAm/internal/pkg/handler"
	"RuAm/internal/pkg/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	db.InitDB(cfg)           // Pass the config to InitDB
	service.InitService(cfg) // Initialize the service with the config
	handler.InitHandler(cfg) // Initialize the handler with the config

	r := gin.Default()

	r.StaticFile("/", "./static/index.html")

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh_token", handler.RefreshToken)
	}

	r.GET("/home", handler.MiddlewareJWT(), handler.Home)

	r.Run(":" + cfg.Port)
}
