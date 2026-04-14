package main

import (
	"log"
	"os"
	"resume-analyzer-api/handlers"
	"resume-analyzer-api/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	godotenv.Load()

	// Initialize services
	groqAPIKey := os.Getenv("GROQ_API_KEY")
	if groqAPIKey == "" {
		log.Fatal("GROQ_API_KEY environment variable is required")
	}
	
        groqService := services.NewGroqService(groqAPIKey)
	pdfService := services.NewPDFService()
	
        // Initialize handlers
	analyzeHandler := handlers.NewAnalyzeHandler(groqService, pdfService)
	
        // Setup Gin router
	r := gin.Default()
	
        // Configure CORS
	allowedOrigins := []string{"http://localhost:3000"}
	if corsOrigin := os.Getenv("CORS_ORIGIN"); corsOrigin != "" {
		allowedOrigins = append(allowedOrigins, corsOrigin)
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Routes
	r.GET("/api/health", handlers.HealthCheck)
	r.POST("/api/analyze", analyzeHandler.Analyze)
	
        // Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
        
	log.Printf("Server starting on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
        }
}
