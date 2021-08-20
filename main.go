package main

import (
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

	api.GET("/halifax", handlers.GetThermometerDataFor(1152))
	api.GET("/calgary", handlers.GetThermometerDataFor(1155))
	api.GET("/edmonton", handlers.GetThermometerDataFor(1156))
	api.GET("/montreal", handlers.GetThermometerDataFor(1148))
	api.GET("/stjohns", handlers.GetThermometerDataFor(1151))
	api.GET("/toronto", handlers.GetThermometerDataFor(1141))
	api.GET("/vancouver", handlers.GetThermometerDataFor(1157))
	r.Use(static.Serve("/", static.LocalFile("dist/llsoup", false)))
	r.Run(":" + port)
}
