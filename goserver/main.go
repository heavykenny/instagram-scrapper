package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"math"
	"net/http"
	"strings"
)

type InstagramData struct {
	CountryCode string `json:"country_code"`
	EntryData   struct {
		ProfilePage []struct {
			Graphql struct {
				User struct {
					EdgeFollowedBy struct {
						Count int `json:"count"`
					} `json:"edge_followed_by"`
					EdgeFollow struct {
						Count int `json:"count"`
					} `json:"edge_follow"`
					FullName                 string `json:"full_name"`
					IsPrivate                bool   `json:"is_private"`
					IsVerified               bool   `json:"is_verified"`
					ProfilePicURL            string `json:"profile_pic_url"`
					Username                 string `json:"username"`
					EdgeOwnerToTimelineMedia struct {
						Count int `json:"count"`
						Edges []struct {
							Node `json:"node"`
						} `json:"edges"`
					} `json:"edge_owner_to_timeline_media"`
				} `json:"user"`
			} `json:"graphql"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
}

type Node struct {
	IsVideo            bool `json:"is_video"`
	VideoViewCount     int  `json:"video_view_count"`
	EdgeMediaToComment struct {
		Count int `json:"count"`
	} `json:"edge_media_to_comment"`
	EdgeLikedBy struct {
		Count int `json:"count"`
	} `json:"edge_liked_by"`
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.Use(CORSMiddleware())

	r.GET("/api/username/:u", GetInstagramDetails)
	err := r.Run()
	if err != nil {
		return
	}
}

func CalculateEngagement(receiver InstagramData) float64 {
	var average float64
	var totalFollowers int

	if receiver.EntryData.ProfilePage != nil {
		totalFollowers = receiver.EntryData.ProfilePage[0].Graphql.User.EdgeFollowedBy.Count

		likes := receiver.EntryData.ProfilePage[0].Graphql.User.EdgeOwnerToTimelineMedia.Edges
		count := 1.0
		for i := range likes {
			count++
			average += (float64(likes[i].EdgeLikedBy.Count) + float64(likes[i].EdgeMediaToComment.Count)) / float64(totalFollowers) * 100
		}

		return math.Round((average/count)*100) / 100
	}
	return 0
}

func GetInstagramDetails(context *gin.Context) {
	username := context.Param("u")
	var jsonMap InstagramData

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("Referrer", "https://www.instagram.com/"+username)
	})

	c.OnHTML("script[type=\"text/javascript\"]", func(e *colly.HTMLElement) {

		val := e.Text

		if strings.HasPrefix(val, "window._sharedData = {") {
			err := json.Unmarshal([]byte(val[21:len(val)-1]), &jsonMap)
			if err != nil {
				return
			}
		}
	})

	err := c.Visit("https://www.instagram.com/" + username)
	if err != nil {
		return
	}

	//sampleResponse, _ := os.Open("sample.json")
	//byteValue, _ := ioutil.ReadAll(sampleResponse)
	//var result InstagramData
	//json.Unmarshal(byteValue, &result)

	context.JSON(http.StatusOK, CalculateEngagement(jsonMap))
}
