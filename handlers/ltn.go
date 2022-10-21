package handlers

import (
	"fmt"
	"llsoup/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/gin-gonic/gin"
)

const (
	Halifax       = "halifax"
	Calgary       = "calgary"
	Edmonton      = "edmonton"
	Montreal      = "montreal"
	StJohns       = "stjohns"
	Toronto       = "toronto"
	Vancouver     = "vancouver"
	Ottawa        = "ottawa"
	London        = "london"
	Winnipeg      = "winnipeg"
	QuebecCity    = "QuebecCity"
	Regina        = "regina"
	Saskatoon     = "saskatoon"
	Charlottetown = "charlottetown"
	Fredericton   = "fredericton"
	Canada        = "canada"
	Blueprint     = "blueprint"
)

var CityIdMap = map[string]int{
	Edmonton:      1351,
	Halifax:       1352,
	London:        1353,
	Montreal:      1354,
	Ottawa:        1355,
	StJohns:       1356,
	Toronto:       1357,
	Vancouver:     1358,
	Winnipeg:      1359,
	QuebecCity:    1360,
	Regina:        1361,
	Saskatoon:     1362,
	Calgary:       1350,
	Charlottetown: 1363,
	Fredericton:   1364,
	Blueprint:     1370,
}

func GetAllData(c *gin.Context) {
	combinedData := models.AllThermometers{}
	for city, cityID := range CityIdMap {
		val, err := getThermometerDataFromID(cityID)
		if err != nil {
			fmt.Printf("Error parsing LLS %s amount raised", city)
			val = &models.Thermometer{Raised: 0, Goal: 1}
		}
		combinedData[city] = val
	}
	totalRaised := 0.0
	for _, j := range combinedData {
		totalRaised += j.Raised
	}
	combinedData[Canada] = &models.Thermometer{
		Raised: totalRaised,
		Goal:   5800000,
	}
	c.JSON(http.StatusOK, combinedData)
}

func GetThermometerDataFor(id int) gin.HandlerFunc {
	return func(c *gin.Context) {
		thermDATA, err := getThermometerDataFromID(id)
		if err != nil {
			c.AbortWithError(http.StatusBadGateway, fmt.Errorf("Error parsing LLS Halifax amount raised"))
			return
		}
		c.JSON(http.StatusOK, thermDATA)
	}
}

func getThermometerDataFromID(id int) (*models.Thermometer, error) {
	resp, err := soup.Get(fmt.Sprintf("https://secure.llscanada.org/site/TR?fr_id=%d&pg=entry&s_locale=en_CA", id))
	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)
	thermometer := doc.Find("div", "id", "thermometer")
	goalh5 := thermometer.FindPrevElementSibling()
	goal, err := parseGoalData(goalh5.Text())
	if err != nil {
		return nil, err
	}
	raisedH3 := thermometer.FindNextElementSibling().Find("h3")
	dollarAmount, err := parseRaised(raisedH3.Text())
	if err != nil {
		return nil, err
	}
	return &models.Thermometer{
		Raised: dollarAmount,
		Goal:   goal,
	}, nil
}

func parseGoalData(goal string) (uint64, error) {
	goal = strings.TrimPrefix(goal, " of $")
	goal = strings.TrimSuffix(goal, " raised")
	goal = strings.ReplaceAll(goal, ",", "")
	return strconv.ParseUint(goal, 10, 64)
}

func parseRaised(raised string) (float64, error) {
	raised = strings.ReplaceAll(raised, "$", "")
	raised = strings.ReplaceAll(raised, ",", "")
	return strconv.ParseFloat(raised, 64)
}
