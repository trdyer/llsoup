package handlers

import (
	"encoding/json"
	"fmt"
	"llsoup/models"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/gin-gonic/gin"
)

const (
	Calgary       = "calgary"
	Halifax       = "halifax"
	Montreal      = "montreal"
	Ottawa        = "ottawa"
	Toronto       = "toronto"
	Vancouver     = "vancouver"
	// Edmonton      = "edmonton"
	// StJohns       = "stjohns"
	// London        = "london"
	// Winnipeg      = "winnipeg"
	// QuebecCity    = "QuebecCity"
	// Regina        = "regina"
	// Saskatoon     = "saskatoon"
	// Charlottetown = "charlottetown"
	// Fredericton   = "fredericton"
	Canada        = "canada"
	// Blueprint     = "blueprint"
)

var CityIdMap = map[string]int{
	Calgary:       1410,
	Halifax:       1412,
	Montreal:      1414,
	Ottawa:        1415,
	Toronto:       1417,
	Vancouver:     1418,


	// Edmonton:      1351,
	// London:        1353,
	// StJohns:       1356,
	// Winnipeg:      1359,
	// QuebecCity:    1360,
	// Regina:        1361,
	// Saskatoon:     1362,
	// Charlottetown: 1363,
	// Fredericton:   1364,
	// Blueprint:     1370,
}

var CityGoalMap = map[string]uint64{
	Calgary:       550000,
	Halifax:       824000,
	Montreal:      1200000,
	Ottawa:        352500,
	Toronto:       1341500,
	Vancouver:     1075000,
	// London:        275000,
	// StJohns:       200000,
	// Winnipeg:      225000,
	// QuebecCity:    100000,
	// Regina:        40000,
	// Saskatoon:     143000,
	// Charlottetown: 27300,
	// Fredericton:   31500,
	// Blueprint:     1000000,
	Canada:        5800000,
}

func GetAllData(c *gin.Context) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	api := &APIAccessor{
		c: client,
	}
	if err := api.Login(); err != nil {
		fmt.Println(err)
		c.AbortWithError(500, err)
		return
	}

	combinedData := models.AllThermometers{}
	for city := range CityIdMap {
		val, err := api.GetDataForCity(city)
		if err != nil {
			fmt.Printf("Error parsing LLS %s amount raised: %v", city, err)
			val = &models.Thermometer{Raised: 0, Goal: 1}
		}
		combinedData[city] = val
	}
	totalRaised := 0.0
	for city, j := range combinedData {
		totalRaised += j.Raised
		fmt.Printf("New! city: %s, cityTotal: %.2f, subtotal: %.2f\n", city, j.Raised, totalRaised)
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

func AttemptAPIAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		jar, _ := cookiejar.New(nil)
		client := &http.Client{
			Jar: jar,
		}
		api := &APIAccessor{
			c: client,
		}
		if err := api.Login(); err != nil {
			fmt.Println(err)
			c.AbortWithError(500, err)
			return
		}

		combinedData := models.AllThermometers{}
		for city := range CityIdMap {
			val, err := api.GetDataForCity(city)
			if err != nil {
				fmt.Printf("Error parsing LLS %s amount raised: %v", city, err)
				val = &models.Thermometer{Raised: 0, Goal: 1}
			}
			combinedData[city] = val
		}
		totalRaised := 0.0
		for city, j := range combinedData {
			totalRaised += j.Raised
			fmt.Printf("New! city: %s, cityTotal: %.2f, subtotal: %.2f\n", city, j.Raised, totalRaised)
		}
		combinedData[Canada] = &models.Thermometer{
			Raised: totalRaised,
			Goal:   5800000,
		}
		c.JSON(http.StatusOK, combinedData)
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

type APIAccessor struct {
	c          *http.Client
	JSESSIONID string
	Token      string
	RoutingID  string
}

func (a *APIAccessor) Login() error {
	_, err := a.c.Get("https://secure.llscanada.org/site/SPageServer/?pagename=LTN_2022_national")
	if err != nil {
		return fmt.Errorf("Error getting JSESSIONID: %w", err)
	}
	resp, err := a.c.Get("https://secure.llscanada.org/site/CRConsAPI?luminateExtend=1.8.2&api_key=qx8ztp18oatitUCr&method=getLoginUrl&response_format=json&v=1.0")
	if err != nil {
		return fmt.Errorf("Error calling lls Login: %w", err)
	}
	responseBody := &responseBody{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(responseBody)
	if err != nil {
		return fmt.Errorf("Error parsing LLS Login: %w")
	}
	a.JSESSIONID = responseBody.GetLoginResponse.JSESSIONID
	a.Token = responseBody.GetLoginResponse.Token
	return nil
}

func (a *APIAccessor) GetDataForCity(city string) (*models.Thermometer, error) {
	cityID := CityIdMap[city]
	cityGoal := CityGoalMap[city]
	urlValues := url.Values{
		"method":                  []string{"render"},
		"content":                 []string{fmt.Sprintf("[[S42:%d:dollars]]", cityID)},
		"api_key":                 []string{"qx8ztp18oatitUCr"},
		"response_format":         []string{"json"},
		"suppress_response_codes": []string{"true"},
		"v":                       []string{"1.0"},
		"auth":                    []string{a.Token},
		"JSESSIONID":              []string{a.JSESSIONID},
		"ts":                      []string{fmt.Sprintf("%13d", time.Now().Unix())},
	}
	resp2, err := a.c.PostForm(fmt.Sprintf("https://secure.llscanada.org/site/CRContentAPI;jsessionid=%s", a.RoutingID), urlValues)
	if err != nil {
		return nil, fmt.Errorf("Error calling LLS API: %w", err)
	}
	apiResponse := &APIResponse{}
	defer resp2.Body.Close()
	err = json.NewDecoder(resp2.Body).Decode(apiResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing LLS API Data: %w", err)
	}
	newAmount := strings.ReplaceAll(apiResponse.RenderResponse.Content, "$", "")
	newAmount = strings.ReplaceAll(newAmount, ",", "")
	raised, err := strconv.ParseFloat(newAmount, 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing LLS API Data into number: %w", err)
	}
	return &models.Thermometer{
		Goal:   cityGoal,
		Raised: raised,
	}, nil
}

type responseBody struct {
	GetLoginResponse LoginResponse `json:"getLoginUrlResponse"`
}

type LoginResponse struct {
	JSESSIONID string `json:"JSESSIONID"`
	RoutingID  string `json:"routing_id"`
	Url        string `json:"url"`
	Token      string `json:"token"`
}

type APIResponse struct {
	RenderResponse ContentBody `json:"renderResponse"`
}

type ContentBody struct {
	Content string `json:"content"`
}
