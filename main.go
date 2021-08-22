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

	api.GET(fmt.Sprintf("/%s", handlers.Halifax), handlers.GetThermometerDataFor(handlers.CityIdMap[handlers.Halifax]))
	api.GET(fmt.Sprintf("/%s", handlers.Calgary), handlers.GetThermometerDataFor(handlers.CityIdMap[handlers.Calgary]))
	api.GET(fmt.Sprintf("/%s", handlers.Edmonton), handlers.GetThermometerDataFor(handlers.CityIdMap[handlers.Edmonton]))
	api.GET(fmt.Sprintf("/%s", handlers.Montreal), handlers.GetThermometerDataFor(handlers.CityIdMap[handlers.Montreal]))
	api.GET(fmt.Sprintf("/%s", handlers.StJohns), handlers.GetThermometerDataFor(handlers.CityIdMap[handlers.StJohns]))
	api.GET(fmt.Sprintf("/%s", handlers.Toronto), handlers.GetThermometerDataFor(handlers.CityIdMap[handlers.Toronto]))
	api.GET(fmt.Sprintf("/%s", handlers.Vancouver), handlers.GetThermometerDataFor(handlers.CityIdMap[handlers.Vancouver]))
	api.GET(fmt.Sprintf("/%s", handlers.London), handlers.GetThermometerDataFor(handlers.CityIdMap[handlers.London]))
	api.GET(fmt.Sprintf("/%s", handlers.Ottawa), handlers.GetThermometerDataFor(handlers.CityIdMap[handlers.Ottawa]))
	api.GET(fmt.Sprintf("/%s", handlers.Winnipeg), handlers.GetThermometerDataFor(handlers.CityIdMap[handlers.Winnipeg]))
	r.Use(static.Serve("/", static.LocalFile("dist/llsoup", false)))
	r.Run(":" + port)
}
