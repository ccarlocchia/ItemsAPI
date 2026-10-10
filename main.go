package main

import (
	"log"
	"net/http"
	"os"
	"project/internal/repositories"
	"project/internal/services"
	"strings"

	"project/cmd/api/handlers"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from GO!",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	repository := repositories.NewRepository()
	service := services.NewItemsService(&repository)
	handler := handlers.NewHandler(service)
	handler.SetupRoutes(r)

	return r
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	return port
}

func main() {
	r := setupRouter()

	if err := r.Run(getPort()); err != nil {
		log.Fatal(err)
	}
}
