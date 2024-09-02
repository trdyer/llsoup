package handlers

import (
	"encoding/json"
	"fmt"
	"llsoup/models"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	Calgary   = "calgary"
	Halifax   = "halifax"
	Montreal  = "montreal"
	Ottawa    = "ottawa"
	Toronto   = "toronto"
	Vancouver = "vancouver"
	// QuebecCity    = "QuebecCity"
	// Charlottetown = "charlottetown"
	// Fredericton   = "fredericton"
	Canada = "canada"
	// Blueprint     = "blueprint"
	StJohns    = "stjohns"    //1416
	Regina     = "regina"     // 1480
	Winnipeg   = "winnipeg"   //1419
	Edmonton   = "edmonton"   //1411
	Saskatoon  = "saskatoon"  // 1490
	London     = "london"     //1413
	Laval      = "laval"      //1450
	Belleville = "belleville" //1470
	Kingston   = "kingston"   //1600
)

var CityIdMap = map[string]int{
	Calgary:    1590,
	Edmonton:   1591,
	Halifax:    1592,
	London:     1593,
	Montreal:   1594,
	Ottawa:     1595,
	StJohns:    1596,
	Toronto:    1597,
	Vancouver:  1598,
	Winnipeg:   1599,
	Laval:      1630,
	Belleville: 1610,
	Saskatoon:  1601,
	Regina:     1620,
	Kingston:   1600,
}

var CityGoalMap = map[string]uint64{
	Calgary:   300000,
	Halifax:   1000000,
	Montreal:  1200000,
	Ottawa:    377000,
	Toronto:   1400000,
	Vancouver: 1000000,
	// London:        275000,
	// StJohns:       200000,
	// Winnipeg:      225000,
	// QuebecCity:    100000,
	// Regina:        40000,
	// Saskatoon:     143000,
	// Charlottetown: 27300,
	// Blueprint:     1000000,
	// Fredericton:   31500,
	Canada:     5800000,
	Laval:      20000,
	Belleville: 100000,
	London:     130000,
	Saskatoon:  40000,
	Edmonton:   65000,
	Winnipeg:   100000,
	Regina:     47000,
	StJohns:    175000,
}

var includeCity = map[string]bool{
	Calgary:    true,
	Halifax:    true,
	Montreal:   true,
	Ottawa:     true,
	Toronto:    true,
	Vancouver:  true,
	Edmonton:   true,
	London:     true,
	StJohns:    true,
	Winnipeg:   true,
	Laval:      true,
	Belleville: true,
	Saskatoon:  true,
	Regina:     true,
	Kingston:   true,
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
	// communityRaised := 0.0
	combinedData := models.AllThermometers{}
	for city := range CityIdMap {
		val, err := api.GetDataForCity(city)
		if err != nil {
			fmt.Printf("Error parsing LLS %s amount raised: %v", city, err)
			val = &models.Thermometer{Raised: 0, Goal: 1}
		}
		if ok := includeCity[city]; ok {
			combinedData[city] = val
		}
		//  else {
		// 	communityRaised += val.Raised
		// }
	}
	// combinedData["Community"] = &models.Thermometer{Raised: communityRaised, Goal: 677000}

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
