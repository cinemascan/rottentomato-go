package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/cinemascan/rottentomato-go/rotten_tomato"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(mwLogger())

	proxyUrl := os.Getenv("SCRAPE_PROXY_URL")

	// TODO: add routes
	router.NoRoute(invalidRouteHandler())
	router.GET("/ping", pingHandler())
	router.GET("/movie", movieInfoHandler(proxyUrl))
	return router
}

func invalidRouteHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "Not Found")
	}
}

func pingHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	}
}

/*
`/search/movie` handler: returns existing movie with EXACT matching params,
otherwise return first result from rotten tomato search (with matching year if year provided)
*/
func movieInfoHandler(proxyUrl string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		title := ctx.Query("title")
		yearParam := ctx.Query("year")

		if title == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "`title` required",
			})
			return
		}

		year := -1
		var err error
		if yearParam != "" {
			year, err = strconv.Atoi(yearParam)
			if err != nil {
				ctx.Error(err)
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid year",
				})
				return
			}
		}

		// attempt to scrape
		scrapedRtInfo, err := rotten_tomato.GetMovieInfo(title, year, proxyUrl)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if scrapedRtInfo == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "no valid results found",
			})
			return
		}
		// reject if incorrect year
		if year != -1 && scrapedRtInfo.Year != year {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "no valid results found",
			})
			return
		}
		ctx.JSON(http.StatusOK, *scrapedRtInfo)
	}
}
