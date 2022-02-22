package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	docs "github.com/miles990/kahaapi/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

// type SheetData struct {
// 	ID                   string `json:"ID"`
// 	Sp                   string `json:"SP,omitempty"`
// 	NameContextID        string `json:"NameContextID,omitempty"`
// 	DescriptionContextID string `json:"DescriptionContextID,omitempty"`
// 	StoryContextID       string `json:"StoryContextID, omitempty"`
// 	AnimationInfo        string `json:"AnimationInfo"`
// 	References           string `json:"References"`
// 	Command              string `json:"Command"`
// 	IsDrop               string `json:"IsDrop"`
// 	Tag                  string `json:"Tag"`
// }

// type SheetDataResponse struct {
// 	ID                   json.Number  `json:"ID"`
// 	Sp                   *json.Number `json:"SP,omitempty"`
// 	NameContextID        *json.Number `json:"NameContextID,omitempty"`
// 	DescriptionContextID *json.Number `json:"DescriptionContextID,omitempty"`
// 	StoryContextID       *json.Number `json:"StoryContextID,omitempty"`
// 	AnimationInfo        string       `json:"AnimationInfo"`
// 	References           string       `json:"References"`
// 	Command              string       `json:"Command"`
// 	IsDrop               *json.Number `json:"IsDrop,omitempty"`
// 	Tag                  string       `json:"Tag"`
// 	// NOEX                 string `json:"-, NOEX_技能故事(不輸出),omitempty"`
// }

type SheetData map[string]interface{}

type SheetDataResponse map[string]interface{}

func getNumberVaule(data string) *json.Number {
	_, err := strconv.Atoi(data)
	if err != nil {
		return nil
	}
	val := json.Number(data)
	return &val
}

func getSheetData(sheetId string, sheetName string) ([]byte, interface{}, error) {
	var apiUrl = fmt.Sprintf("https://opensheet.elk.sh/%s/%s", sheetId, sheetName)
	resp, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println("No response from request")
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	//fmt.Println(string(body))

	var results []SheetData = make([]SheetData, 0)
	if err := json.Unmarshal(body, &results); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return nil, nil, err
	}
	// fmt.Println(string(body))

	var dataResponse []SheetDataResponse = make([]SheetDataResponse, 0)
	for _, data := range results {
		value, exist := data["ID"]
		if !exist || value == "" {
			continue
		}

		var responseData = SheetDataResponse{}

		for k, v := range data {
			// fmt.Printf("%+v %+v\n", k, v)

			if strings.Contains(k, "NOEX") {
				continue
			}

			var valStr = fmt.Sprintf("%v", v)
			var valNumber = getNumberVaule(valStr)
			if valNumber != nil {
				responseData[k], _ = valNumber.Int64()
			} else {
				if valStr != "" {
					responseData[k] = v
				}
			}
		}
		dataResponse = append(dataResponse, responseData)
	}
	return body, dataResponse, err
}

// SheetData godoc
// @Summary get sheetdata example
// @Schemes
// @Description do sheetData
// @Tags SheetData
// @Param id query string true "sheet id"
// @Param name query string true "sheet name"
// @Param pretty query string false "pretty json"
// @Produce json
// @Success 200 {string} Helloworld
// @Router /sheetData [get]
func GetSheetData(c *gin.Context) {
	sheetId := c.Query("id")
	sheetName := c.Query("name")
	pretty := c.Query("pretty")

	_, skillDatas, err := getSheetData(sheetId, sheetName)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if pretty == "true" || pretty == "1" {
		c.IndentedJSON(200, skillDatas)
	} else {
		c.JSON(200, skillDatas)
	}

}

func main() {

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{

		eg := v1.Group("/sheetData")
		{
			eg.GET("/", GetSheetData)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
