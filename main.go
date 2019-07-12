package main

import (
	"fmt"
	"log"
	"os"

	"shorty/endpoints"

	"shorty/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	shortyEnv := os.Getenv("SHORTY_ENVIRONMENT")

	// Load from .env if development or travis
	if shortyEnv == "" || shortyEnv == "development" || shortyEnv == "travis" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.GET("/", endpoints.Index)

	// Create a DB Connection
	dbConn, err := db.NewDatabaser(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("NewDatabaser: ", err)
	}

	// Create an Endpointer
	e := endpoints.NewEndpointer(dbConn)

	authorized := router.Group("/api", gin.BasicAuth(gin.Accounts{
		os.Getenv("USERNAME"): os.Getenv("PASSWORD"),
	}))

	authorized.POST("/shorten", e.NewShortenedURL)

	// Redirect URL
	router.GET("/r/:slug", e.RedirectURL)

	// Generic 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"status": false, "message": "Endpoint not found"})
	})

	router.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
