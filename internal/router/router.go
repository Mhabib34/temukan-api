package router

import (
	"temukan-api/internal/handler"
	"temukan-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(handler handler.UserHandler) *gin.Engine {
	r := gin.New()

	// middleware
	r.Use(gin.Logger())
	r.Use(middleware.ErrorRecovery())
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api/v1/auth")
	{
		api.POST("/register", handler.Create)
		api.POST("/login", handler.Login)
		api.POST("/refresh", handler.RefreshToken)

		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware())
		{
			authorized.POST("/me", handler.Me)
			authorized.GET("/logout", handler.Logout)
		}
	}

	return r
}
