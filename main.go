package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	docs "github.com/miles990/kahaapi/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

type SheetData struct {
	ID                   string `json:"ID"`
	Sp                   string `json:"SP,omitempty"`
	NameContextID        string `json:"NameContextID,omitempty"`
	DescriptionContextID string `json:"DescriptionContextID,omitempty"`
	StoryContextID       string `json:"StoryContextID, omitempty"`
	AnimationInfo        string `json:"AnimationInfo"`
	References           string `json:"References"`
	Command              string `json:"Command"`
	IsDrop               string `json:"IsDrop"`
	Tag                  string `json:"Tag"`
}

type SheetDataResponse struct {
	ID                   json.Number  `json:"ID"`
	Sp                   *json.Number `json:"SP"`
	NameContextID        *json.Number `json:"NameContextID"`
	DescriptionContextID *json.Number `json:"DescriptionContextID"`
	StoryContextID       *json.Number `json:"StoryContextID"`
	AnimationInfo        string       `json:"AnimationInfo"`
	References           string       `json:"References"`
	Command              string       `json:"Command"`
	IsDrop               *json.Number `json:"IsDrop"`
	Tag                  string       `json:"Tag"`
	// NOEX                 string `json:"-, NOEX_技能故事(不輸出),omitempty"`
}

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
		if data.ID == "" {
			continue
			// responseData = append(responseData, data)
			// skillDataMap[data.ID] = data
			// fmt.Printf("%+v\n", data)
			// fmt.Println(data.ID, data.StoryContextID)
		}
		var responseData = SheetDataResponse{
			ID:                   json.Number(data.ID),
			Sp:                   getNumberVaule(data.Sp),
			NameContextID:        getNumberVaule(data.NameContextID),
			DescriptionContextID: getNumberVaule(data.DescriptionContextID),
			StoryContextID:       getNumberVaule(data.StoryContextID),
			AnimationInfo:        data.AnimationInfo,
			References:           data.References,
			Command:              data.Command,
			IsDrop:               getNumberVaule(data.IsDrop),
			Tag:                  data.Tag,
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
// @Produce json
// @Success 200 {string} Helloworld
// @Router /sheetData [get]
func GetSheetData(c *gin.Context) {
	sheetId := c.Query("id")
	sheetName := c.Query("name")
	_, skillDatas, err := getSheetData(sheetId, sheetName)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, skillDatas)
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
