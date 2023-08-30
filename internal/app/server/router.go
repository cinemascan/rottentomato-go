package server

import (
	"net/http"
	"os"
	"strconv"

	_ "github.com/cinemascan/rottentomato-go/internal/pkg/docs"
	"github.com/cinemascan/rottentomato-go/rotten_tomato"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath /

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(mwLogger())

	proxyUrl := os.Getenv("SCRAPE_PROXY_URL")

	// TODO: add routes
	router.NoRoute(invalidRouteHandler())
	router.GET("/ping", pingHandler())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL("http://localhost:8081/swagger/doc.json")))
	router.GET("/movie", movieInfoHandler(proxyUrl))
	return router
}

func invalidRouteHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "Not Found")
	}
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} pong
// @Router /ping [get]
func pingHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	}
}

// MovieSearch godoc
// @Summary top search result scraped from rotten tomato
// @Schemes
// @Description scrapes https://rottentomatoes.com/search url and returns the top result that matches title, year params provided
// @Tags search
// @Accept json
// @Produce json
// @Param title query string false "movie title" minlength(1)
// @Param year query int false "movie release year" minimum(1972) maximum(2023)
// @Success 200 {object} rotten_tomato.RTMovieInfo "title year"
// @Router /movie [get]
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
