package router

import (
	"os"
	"strings"
	"temukan-api/internal/handler"
	"temukan-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler handler.UserHandler,
	reportHandler handler.ReportHandler,
	matchHandler handler.MatchHandler,
	notificationHandler handler.NotificationHandler,
) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(middleware.ErrorRecovery())
	r.Use(corsMiddleware())

	api := r.Group("/api/v1")

	// ── Health ──────────────────────────────────────────────────────────────
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "OK",
			"data": gin.H{
				"version": "1.0.0",
				"env":     "development",
			},
		})
	})

	// ── Auth ─────────────────────────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/register", userHandler.Create)
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.RefreshToken)

		authPrivate := auth.Group("/")
		authPrivate.Use(middleware.AuthMiddleware())
		{
			authPrivate.GET("/me", userHandler.Me)
			authPrivate.POST("/logout", userHandler.Logout)
		}
	}

	// ── Reports ───────────────────────────────────────────────────────────────
	reports := api.Group("/reports")
	{
		// Public
		reports.GET("", reportHandler.GetAll)
		reports.GET("/:id", reportHandler.GetByID)

		// Private
		reportsPrivate := reports.Group("")
		reportsPrivate.Use(middleware.AuthMiddleware())
		{
			reportsPrivate.GET("/my", reportHandler.GetMyReports)
			reportsPrivate.POST("", reportHandler.Create)
			reportsPrivate.PUT("/:id", reportHandler.Update)
			reportsPrivate.DELETE("/:id", reportHandler.Delete)
			reportsPrivate.POST("/:id/photo", reportHandler.UploadPhoto)
		}
	}

	// ── Map ───────────────────────────────────────────────────────────────────
	mapGroup := api.Group("/map")
	{
		mapGroup.GET("/pins", reportHandler.GetMapPins)
	}

	// ── Matches ───────────────────────────────────────────────────────────────
	matches := api.Group("/matches")
	matches.Use(middleware.AuthMiddleware())
	{
		matches.GET("", matchHandler.GetAll)
		matches.GET("/:id", matchHandler.GetByID)
	}

	// ── Notifications ─────────────────────────────────────────────────────────
	notifications := api.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware())
	{
		notifications.GET("", notificationHandler.GetAll)
		notifications.PATCH("/read-all", notificationHandler.GetAll)
		notifications.PATCH("/:id/read", notificationHandler.MarkAsRead)
	}

	return r
}

// corsMiddleware mengganti wildcard "*" dengan origin spesifik agar
// kompatibel dengan withCredentials: true (httpOnly cookie).
func corsMiddleware() gin.HandlerFunc {
	// Baca dari env: ALLOWED_ORIGINS=http://localhost:4000,https://temukan.id
	raw := os.Getenv("ALLOWED_ORIGINS")
	allowed := map[string]struct{}{}
	if raw != "" {
		for _, o := range strings.Split(raw, ",") {
			allowed[strings.TrimSpace(o)] = struct{}{}
		}
	} else {
		// Default development
		allowed["http://localhost:3000"] = struct{}{}
		allowed["http://localhost:4000"] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if _, ok := allowed[origin]; ok {
			c.Header("Access-Control-Allow-Origin", origin)      // ← spesifik, bukan *
			c.Header("Access-Control-Allow-Credentials", "true") // ← wajib untuk cookie
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Client-Type")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
