package main

import (
	"fmt"
	"llsoup/handlers"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")

	r := gin.Default()
	r.Use(cors.Default())
	api := r.Group("/api")
	api.GET("/all", handlers.GetAllData)
	api.GET("/test", handlers.AttemptAPIAccess())
	r.Use(static.Serve("/", static.LocalFile("dist/llsoup", false)))
	r.Run(":" + port)
}
