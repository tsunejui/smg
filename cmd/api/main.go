package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"smg/pkg/config"
	"smg/pkg/handlers"
	"smg/pkg/middleware"
	"smg/pkg/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Initialize configuration
	cfg := config.New()

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: "",
		DB:       0,
	})

	// Initialize services
	userService := services.NewUserService(db)
	topicService := services.NewTopicService(db)
	mediaService := services.NewMediaService(db)
	articleService := services.NewArticleService(db)
	systemService := services.NewSystemService(db)
	authService := services.NewAuthService(db, redisClient)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	topicHandler := handlers.NewTopicHandler(topicService)
	mediaHandler := handlers.NewMediaHandler(mediaService)
	articleHandler := handlers.NewArticleHandler(articleService)
	systemHandler := handlers.NewSystemHandler(systemService)

	// Setup Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/qr-generate", authHandler.GenerateQRCode)
		auth.POST("/qr-verify", authHandler.VerifyQRCode)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(authService))
	{
		// User routes
		users := api.Group("/users")
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.GET("/", userHandler.GetUsers) // Admin only
			users.GET("/:id", userHandler.GetUserByID)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Topic routes
		topics := api.Group("/topics")
		{
			topics.GET("/", topicHandler.GetTopics)
			topics.POST("/", topicHandler.CreateTopic)
			topics.GET("/:id", topicHandler.GetTopic)
			topics.PUT("/:id", topicHandler.UpdateTopic)
			topics.DELETE("/:id", topicHandler.DeleteTopic)
		}

		// Media account routes
		media := api.Group("/media")
		{
			media.GET("/accounts", mediaHandler.GetAccounts)
			media.POST("/accounts", mediaHandler.CreateAccount)
			media.GET("/accounts/:id", mediaHandler.GetAccount)
			media.PUT("/accounts/:id", mediaHandler.UpdateAccount)
			media.DELETE("/accounts/:id", mediaHandler.DeleteAccount)
			media.POST("/connect/:platform", mediaHandler.ConnectPlatform)
			media.POST("/disconnect/:id", mediaHandler.DisconnectAccount)
		}

		// Article routes
		articles := api.Group("/articles")
		{
			articles.GET("/", articleHandler.GetArticles)
			articles.POST("/", articleHandler.CreateArticle)
			articles.GET("/:id", articleHandler.GetArticle)
			articles.PUT("/:id", articleHandler.UpdateArticle)
			articles.DELETE("/:id", articleHandler.DeleteArticle)
			articles.POST("/:id/repost", articleHandler.RepostArticle)
			articles.GET("/reposts", articleHandler.GetReposts)
		}

		// System settings routes (Admin only)
		system := api.Group("/system")
		system.Use(middleware.AdminMiddleware())
		{
			system.GET("/settings", systemHandler.GetSettings)
			system.PUT("/settings", systemHandler.UpdateSettings)
			system.GET("/stats", systemHandler.GetStats)
			system.GET("/platforms", systemHandler.GetPlatforms)
			system.POST("/platforms", systemHandler.CreatePlatform)
			system.PUT("/platforms/:id", systemHandler.UpdatePlatform)
			system.DELETE("/platforms/:id", systemHandler.DeletePlatform)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}