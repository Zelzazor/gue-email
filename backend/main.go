package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
)

type Query struct {
	Term string `json:"term"`
}
type Search struct {
	SearchType  string    `json:"search_type"`
	Query       Query     `json:"query"`
	From        int       `json:"from"`
	Max_results int       `json:"max_results"`
	Highlight   Highlight `json:"highlight"`
}

type Highlight struct {
	Fields Fields `json:"fields"`
}

type Fields struct {
	Content Content `json:"Content"`
}

type Content struct {
	PreTags  []string `json:"pre_tags"`
	PostTags []string `json:"post_tags"`
}

func search(c *gin.Context) {

	from := 0
	if c.Query("offset") != "" {
		fromInt, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			from = 0
		} else {
			from = fromInt
		}

	}
	limit := 20
	if c.Query("limit") != "" {
		limitInt, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 20
		} else {
			limit = limitInt
		}

	}
	searchType := "matchphrase"
	if c.Query("q") == "" {
		searchType = "matchall"
	}

	body := Search{
		SearchType: searchType,
		Query: Query{
			Term: c.Query("q"),
		},
		From:        from,
		Max_results: limit,
		Highlight: Highlight{
			Fields: Fields{
				Content: Content{
					PreTags:  []string{"<b>"},
					PostTags: []string{"</b>"},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:4080/api/email/_search", strings.NewReader(string(jsonBody)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, gin.MIMEJSON, responseBody)

}

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		println("Error loading .env file")
		return
	}
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/search", search)

	router.Run("localhost:8000")
}
