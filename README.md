![tests](https://github.com/cinemascan/rottentomato-go/actions/workflows/ci.yml/badge.svg)


# :movie_camera: Rotten Tomatoes in Golang (and API) :clapper:

> **Note**
> This is a golang rewrite of the [rottentomatoes-python](https://github.com/preritdas//blob/master/README.md) project, made to be used for [cinemascan.org](https://cinemascan.org)

> **Disclaimer**
> If at any point in your project this library stops working, 99% of the time it's due to Rotten Tomatoes IP-blocking the server (every request scrapes Rotten Tomatoes /search endpoint) OR because the Rotten Tomatoes site schema has changed, meaning some changes to web scraping and extraction under the hood will be necessary to make everything work again.

This package allows you to easily fetch Rotten Tomatoes scores and other movie data such as genres, without the use of the official Rotten Tomatoes API. The package scrapes their website for the data. This package is a golang rewrite of [rottentomatoes-python](https://github.com/preritdas//blob/master/README.md) for higher performance and to be used for storing movie ratings info for [cinemascan.org](cinemascan.org)

The package now, by default, scrapes the Rotten Tomatoes search page to find the true url of the first valid movie response (is a movie and has a tomatometer). This means queries that previously didn't work because their urls had a unique identifier or a year-released prefix, now work. The limitation of this new mechanism is that you only get the top response, and when searching for specific movies (sequels, by year, etc.) Rotten Tomatoes seems to return the same results as the original query. So, it's difficult to use specific queries to try and get the desired result movie as the top response. See #4 for more info on this.

There is now an API deployed to query movies and getting responses easier. The endpoint is https://rottentomato.cinemascan.org and it's open and free to use. Visit the [swagger docs](https://rottentomato.cinemascan.org/swagger/index.html) in the browser to view the endpoints. Both endpoints live right now are browser accessible meaning you don't need an HTTP client to use the API.

- https://rottentomato.cinemascan.org/movies?title=the+matrix&year=1999 for JSON response of the top result


## Usage

Basic usage example:

```go
import (¯
    rotten_tomato "github.com/cinemascan/rottentomato-go"
)

movieName := "The Matrix"
currentYear := 1999
proxyUrl := os.Getenv("PROXY_URL")
scrapedRtInfo, err := rotten_tomato.GetMovieInfo(title, year, proxyUrl)

fmt.Printf("%v", scrapedRtInfo)
//// OUTPUT
// {
//     "audienceScore": {
//         "averageRating": "4.5",
//         "bandedRatingCount": "10,000+",
//         "likedCount": 12460,
//         "notLikedCount": 1248,
//         "ratingCount": 13708,
//         "reviewCount": 5583,
//         "state": "upright",
//         "value": 91
//     },
//     "rating": "R",
//     "tomatometerScore": {
//         "averageRating": "8.60",
//         "bandedRatingCount": "",
//         "likedCount": 429,
//         "notLikedCount": 31,
//         "ratingCount": 460,
//         "reviewCount": 460,
//         "state": "certified-fresh",
//         "value": 93
//     },
//     "title": "Oppenheimer",
//     "year": 2023,
//     "runtime": "3h 0m",
//     "genres": [
//         "History",
//         "Drama"
//     ]
// }
```

## Performance

Since every request queries the Rotten Tomatoes search endpoint, response times can range from 2-3s up to 10s in rarer cases.

If performance is important, you may use cinemascan's private API [https://api.cinemascan.org/search/movies](https://api.cinemascan.org/search/movies?title=the+matrix&year=2003)

We store the ratings for the top movies from 1999 - 2023 in our DB, hence response times range from 50-100+ms depending on location (view response times [here](https://cinemascan-api.planetfall.io))

## API

Try out via swagger: [https://rottentomato.cinemascan.org/swagger/index.html#/](https://rottentomato.cinemascan.org/swagger/index.html#/)