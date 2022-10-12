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
	Halifax   = "halifax"
	Calgary   = "calgary"
	Edmonton  = "edmonton"
	Montreal  = "montreal"
	StJohns   = "stjohns"
	Toronto   = "toronto"
	Vancouver = "vancouver"
	Ottawa    = "ottawa"
	London    = "london"
	Winnipeg  = "winnipeg"
)

var CityIdMap = map[string]int{
	Halifax:   1152,
	Calgary:   1350,
	Edmonton:  1156,
	Montreal:  1148,
	StJohns:   1151,
	Toronto:   1141,
	Vancouver: 1157,
	London:    1143,
	Ottawa:    1142,
	Winnipeg:  1153,
}

func GetAllData(c *gin.Context) {
	hfx, err := getThermometerDataFromID(CityIdMap[Halifax])
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, fmt.Errorf("Error parsing LLS Halifax amount raised"))
		return
	}
	cal, err := getThermometerDataFromID(CityIdMap[Calgary])
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, fmt.Errorf("Error parsing LLS Halifax amount raised"))
		return
	}
	yeg, err := getThermometerDataFromID(CityIdMap[Edmonton])
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, fmt.Errorf("Error parsing LLS Halifax amount raised"))
		return
	}
	mtl, err := getThermometerDataFromID(CityIdMap[Montreal])
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, fmt.Errorf("Error parsing LLS Halifax amount raised"))
		return
	}
	stj, err := getThermometerDataFromID(CityIdMap[StJohns])
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, fmt.Errorf("Error parsing LLS Halifax amount raised"))
		return
	}
	tor, err := getThermometerDataFromID(CityIdMap[Toronto])
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, fmt.Errorf("Error parsing LLS Halifax amount raised"))
		return
	}
	yvr, err := getThermometerDataFromID(CityIdMap[Vancouver])
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, fmt.Errorf("Error parsing LLS Halifax amount raised"))
		return
	}
	lon, err := getThermometerDataFromID(CityIdMap[London])
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, fmt.Errorf("Error parsing LLS Halifax amount raised"))
		return
	}
	ott, err := getThermometerDataFromID(CityIdMap[Ottawa])
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, fmt.Errorf("Error parsing LLS Halifax amount raised"))
		return
	}
	peg, err := getThermometerDataFromID(CityIdMap[Winnipeg])
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, fmt.Errorf("Error parsing LLS Halifax amount raised"))
		return
	}
	c.JSON(http.StatusOK, models.AllThermometers{
		Halifax:   hfx,
		Calgary:   cal,
		Edmonton:  yeg,
		Montreal:  mtl,
		StJohns:   stj,
		Toronto:   tor,
		Vancouver: yvr,
		London:    lon,
		Ottawa:    ott,
		Winnipeg:  peg,
	})
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
