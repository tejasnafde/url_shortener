package main

import (
	// "context"
	"log"
	"time"

	"url_shortener/internal/database"
	"url_shortener/internal/models"
	"url_shortener/internal/repository"
	"url_shortener/pkg/logger"
)

// TODO: Implement main application entry point
//
// What to implement:
// 1. Application initialization
//    - Initialize logger
//    - Initialize database connection
//    - Initialize repository
//    - Initialize service
//    - Initialize handlers
//
// 2. Gin router setup
//    - Configure Gin mode (debug/release)
//    - Set up middleware (logging, recovery, CORS)
//    - Define routes and handlers
//    - Serve static files
//
// 3. Server startup
//    - Configure server settings
//    - Start HTTP server
//    - Graceful shutdown handling
//
// Routes to implement:
// - GET / - Serve frontend HTML
// - POST /shorten - Shorten URL
// - GET /:shortCode - Redirect to original URL
// - GET /urls - Get all URLs (admin)
// - DELETE /urls/:id - Delete URL (admin)
//
// Enterprise patterns to follow:
// - Dependency injection throughout the stack
// - Proper error handling and logging
// - Graceful shutdown
// - Configuration management
// - Health check endpoints
// - Request/response logging middleware

func main() {

	// Initialize logger: levels = debug|info|warn|error
	logger.Init("debug")
	logger.Logger.Info("First INFO log?")

	db, err := database.NewSQLite("urls.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	urlRepo := repository.NewURLRepositorySQL(db)

	u := &models.URL{ShortCode: "lmao", OriginalURL: "https://tejasnafde.github.io/"}
	if err := urlRepo.Create(u); err != nil {
		log.Fatal("Create failed:", err)
	}
	got, err := urlRepo.GetByShortCode("lmao")
	if err != nil {
		log.Fatal("GetByShortCode failed:", err)
	}
	log.Printf("fetched: shortcode=%s url=%s created=%s", got.ShortCode, got.OriginalURL, got.CreatedAt.Format(time.RFC3339))
	// TODO: Initialize repository
	// urlRepo := repository.NewURLRepository(db)

	// TODO: Initialize service
	// urlService := services.NewURLService(urlRepo)

	// TODO: Initialize handlers
	// urlHandler := handlers.NewURLHandler(urlService, appLogger)

	// TODO: Configure Gin router
	// gin.SetMode(gin.ReleaseMode) // Use gin.DebugMode for development
	// router := gin.New()

	// TODO: Add middleware
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())
	// router.Use(corsMiddleware())

	// TODO: Define routes
	// router.GET("/", urlHandler.ServeFrontend)
	// router.POST("/shorten", urlHandler.ShortenURL)
	// router.GET("/:shortCode", urlHandler.RedirectURL)
	// router.GET("/urls", urlHandler.GetAllURLs)
	// router.DELETE("/urls/:id", urlHandler.DeleteURL)

	// TODO: Serve static files
	// router.Static("/static", "./static")

	// TODO: Start server
	// server := &http.Server{
	//     Addr:    ":8080",
	//     Handler: router,
	// }

	// log.Println("Starting server on :8080")
	// if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//     log.Fatal("Failed to start server:", err)
	// }

	// Placeholder to prevent compilation errors
	log.Println("Main function - TODO: Implement application initialization")
}

// TODO: Implewment CORS middleware
// func corsMiddleware() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         c.Header("Access-Control-Allow-Origin", "*")
//         c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//         c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
//
//         if c.Request.Method == "OPTIONS" {
//             c.AbortWithStatus(204)
//             return
//         }
//
//         c.Next()
//     }
// }

// TODO: Implement graceful shutdown
// func setupGracefulShutdown(server *http.Server) {
//     c := make(chan os.Signal, 1)
//     signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//
//     go func() {
//         <-c
//         log.Println("Shutting down server...")
//
//         ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//         defer cancel()
//
//         if err := server.Shutdown(ctx); err != nil {
//             log.Fatal("Server forced to shutdown:", err)
//         }
//     }()
// }
