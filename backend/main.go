package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"plocate-ui/config"
	"plocate-ui/handlers"
	"plocate-ui/indexer"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
	// Load configuration
	if err := config.Load(""); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize indexer
	if err := indexer.Initialize(); err != nil {
		log.Fatalf("Failed to initialize indexer: %v", err)
	}

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API routes
	api := r.Group("/api")
	{
		api.GET("/status", handlers.GetStatus)
		api.GET("/indices", handlers.GetIndices)
		api.GET("/search", handlers.Search)
		api.POST("/search", handlers.Search)
		api.POST("/control/start", handlers.StartIndexing)             // Start all enabled indices
		api.POST("/control/start/:indexName", handlers.StartIndexing)  // Start specific index
		api.POST("/control/stop", handlers.StopIndexing)               // Stop all indices
		api.POST("/control/stop/:indexName", handlers.StopIndexing)    // Stop specific index
		api.POST("/control/scheduler/enable", handlers.EnableScheduler)
		api.POST("/control/scheduler/disable", handlers.DisableScheduler)
	}

	// Serve frontend (embedded or from filesystem)
	serveFrontend(r)

	// Start server
	port := config.AppConfig.Server.Port
	log.Printf("Server starting on port %s", port)
	log.Printf("Database path: %s", config.AppConfig.Plocate.DatabasePath)
	log.Printf("Index paths: %v", config.AppConfig.Plocate.IndexPaths)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func serveFrontend(r *gin.Engine) {
	// Try to serve embedded frontend
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Println("Warning: embedded frontend not found, serving API only")
		return
	}

	// Serve index.html at root
	r.GET("/", func(c *gin.Context) {
		indexHTML, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "index.html not found"})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	// Serve static assets
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Remove leading slash
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}

		// Try to read the file
		data, err := fs.ReadFile(distFS, path)
		if err != nil {
			// If file not found, serve index.html for SPA routing
			indexHTML, err := fs.ReadFile(distFS, "index.html")
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
			return
		}

		// Determine content type based on file extension
		contentType := "application/octet-stream"
		if len(path) > 3 {
			switch {
			case len(path) > 3 && path[len(path)-3:] == ".js":
				contentType = "application/javascript"
			case len(path) > 4 && path[len(path)-4:] == ".css":
				contentType = "text/css"
			case len(path) > 5 && path[len(path)-5:] == ".html":
				contentType = "text/html; charset=utf-8"
			case len(path) > 4 && path[len(path)-4:] == ".png":
				contentType = "image/png"
			case len(path) > 4 && path[len(path)-4:] == ".jpg":
				contentType = "image/jpeg"
			case len(path) > 5 && path[len(path)-5:] == ".jpeg":
				contentType = "image/jpeg"
			case len(path) > 4 && path[len(path)-4:] == ".svg":
				contentType = "image/svg+xml"
			case len(path) > 5 && path[len(path)-5:] == ".woff":
				contentType = "font/woff"
			case len(path) > 6 && path[len(path)-6:] == ".woff2":
				contentType = "font/woff2"
			}
		}

		c.Data(http.StatusOK, contentType, data)
	})

	fmt.Println("Serving embedded frontend")
}
