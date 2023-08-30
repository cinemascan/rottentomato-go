package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/cinemascan/rottentomato-go/internal/pkg/rotten_tomato"
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
	router.GET("/movie", movieInfoHandler())
	router.GET("/search/movie", movieInfoSearchHandler(proxyUrl))
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
`/movie` handler: returns movie data that matches provided params EXACTLY
*/
func movieInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		title := ctx.Query("title")
		yearParam := ctx.Query("year")
		rating := ctx.Query("rating")

		if title == "" {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "`imdb_id` or `title` required",
			})
			return
		}

		// handle title
		year := -1
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
		if year != -1 && rating != "" {
			// title, year, rating
			movieInfoDb, err = rtdb.GetMovieByTitleYearRating(ctx, db, title, year, rating)
		} else if year != -1 && rating == "" {
			// title, year
			movieInfoDb, err = rtdb.GetMovieByTitleYear(ctx, db, title, year)
		} else if year == -1 && rating != "" {
			// title, rating
			movieInfoDb, err = rtdb.GetMovieByTitleRating(ctx, db, title, rating)
		} else {
			// title
			movieInfoDb, err = rtdb.GetMovieByTitle(ctx, db, title)
		}

		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if movieInfoDb.Title == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "no valid results found",
			})
		} else {
			ctx.JSON(http.StatusOK, movieInfoDb)
		}
	}
}

/*
`/search/movie` handler: returns existing movie with EXACT matching params,
otherwise return first result from rotten tomato search (with matching year if year provided)
*/
func movieInfoSearchHandler(proxyUrl string) gin.HandlerFunc {
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

		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
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
