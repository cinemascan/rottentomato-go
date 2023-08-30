package rotten_tomato

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/cinemascan/rottentomato-go/internal/pkg/utils"
)

var REQUEST_HEADERS http.Header = http.Header{
	"User-Agent":      {"Mozilla/5.0 (X11; Linux x86_64; rv:12.0) Gecko/20100101 Firefox/12.0"},
	"Accept-Language": {"en-US"},
	"Accept":          {"text/html"},
	"Referer":         {"https://www.google.com"},
}

const RT_HOST string = "www.rottentomatoes.com"

func getMovieUrl(movieName string) string {
	movieUrl := url.URL{
		Scheme: "https",
		Host:   RT_HOST,
		Path:   "/m/" + strings.ReplaceAll(strings.ToLower(movieName), " ", "_"),
	}
	return movieUrl.String()
}

func getSearchUrl(movieName string) string {
	searchUrl := url.URL{
		Scheme: "https",
		Host:   RT_HOST,
		Path:   "/search",
		RawQuery: url.Values{
			"search": {movieName},
		}.Encode(),
	}
	return searchUrl.String()
}

func extract(content string, start_str string, end_str string) (string, error) {
	start_idx := strings.Index(content, start_str)
	if start_idx == -1 {
		return "", errors.New("Scraper.Extract: `start_str` not found in `content`")
	}
	end_idx := strings.Index(content[start_idx:], end_str) + start_idx
	extract_start_idx := start_idx + len(start_str)
	extracted := content[extract_start_idx:end_idx]
	return extracted, nil
}

// Retrieves the schema.org data model for a movie.
// This data typically contains Tomatometer score, genre etc.
func getSchemaJsonLD(content string) (*RTSchemaJson, error) {
	var json_data RTSchemaJson
	extracted, err := extract(content, `<script type="application/ld+json">`, `</script>`)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(extracted), &json_data)
	if err != nil {
		return nil, err
	}
	return &json_data, nil
}

func getMovieInfoFromContent(content string) (*RTMovieInfo, error) {
	var score_details RTScoreDetails
	extracted, err := extract(content, `<script id="scoreDetails" type="application/json">`, `</script>`)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(extracted), &score_details)
	if err != nil {
		return nil, err
	}
	split := strings.Split(score_details.Scoreboard.Info, ", ")
	if len(split) != 3 {
		return nil, errors.New("Scraper.GetScoreBard: score_details.Scoreboard.Info has missing data (either year, genres or runtime info missing)")
	}
	year, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, err
	}
	genres := strings.Split(split[1], "/")
	runtime := split[2]
	return &RTMovieInfo{
		AudienceScore:    score_details.Scoreboard.AudienceScore,
		Rating:           score_details.Scoreboard.Rating,
		TomatometerScore: score_details.Scoreboard.TomatometerScore,
		Title:            score_details.Scoreboard.Title,
		Year:             year,
		Runtime:          runtime,
		Genres:           genres,
	}, nil
}

func scrapeSearchData(movieName string, client *http.Client) (string, error) {
	req, _ := http.NewRequest("GET", getSearchUrl(movieName), nil)
	req.Header = REQUEST_HEADERS
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	strBody := string(body)
	strBody = strBody[2 : len(strBody)-1]
	return strBody, nil
}

func getSearchResults(movieName string, client *http.Client) ([]SearchListing, error) {
	content, err := scrapeSearchData(movieName, client)
	if err != nil {
		return []SearchListing{}, err
	}
	start_tag := `<search-page-media-row`
	end_tag := `</search-page-media-row>`
	chunks := utils.GetChunksFromString(start_tag, end_tag, content)
	results := []SearchListing{}
	for _, chunk := range chunks {
		result := FromHtml(chunk)
		results = append(results, result)
	}
	return results, nil
}

func filterSearchResults(results []SearchListing, year int) []SearchListing {
	filtered := []SearchListing{}
	for _, res := range results {
		if res.IsMovie {
			if res.Year == year || year == -1 {
				filtered = append(filtered, res)
			}
		}
	}
	return filtered
}

func getTopResult(movieName string, year int, client *http.Client) (*SearchListing, error) {
	results, err := getSearchResults(movieName, client)
	if err != nil {
		return nil, err
	}
	filtered := filterSearchResults(results, year)
	if len(filtered) > 0 {
		return &filtered[0], nil
	}
	return nil, errors.New("no valid results found")
}

func scrapeViaNameYear(movieName string, year int, proxyUrl string) (string, error) {
	var client http.Client
	if proxyUrl != "" {
		proxy, err := url.Parse(os.Getenv("SCRAPE_PROXY_URL"))
		if err != nil {
			return "", err
		}
		client = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
	} else {
		client = http.Client{}
	}
	search_res, err := getTopResult(movieName, year, &client)
	if err != nil {
		return "", err
	}
	search_url := search_res.Url
	req, _ := http.NewRequest("GET", search_url, nil)
	req.Header = REQUEST_HEADERS
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return "", fmt.Errorf("Scraper.ScrapeMovieData: unable to find movie on rotten tomato, try this link to source movie manually %s", search_url)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	strBody := string(body)
	return strBody, nil
}

func scrapeViaUrl(movieName string) (string, error) {
	rt_url := getMovieUrl(movieName)
	client := http.Client{}
	req, _ := http.NewRequest("GET", rt_url, nil)
	req.Header = REQUEST_HEADERS
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return "", fmt.Errorf("Scraper.ScrapeMovieData: unable to find movie on rotten tomato, try this link to source movie manually %s", rt_url)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	strBody := string(body)
	return strBody, nil
}

func GetMovieInfo(movieName string, year int, proxyUrl string) (*RTMovieInfo, error) {
	content, err := scrapeViaNameYear(movieName, year, proxyUrl)
	if err != nil {
		return nil, err
	}
	details, err := getMovieInfoFromContent(content)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func GetMovieTitle(movieName string, year int, proxyUrl string) (string, error) {
	content, err := scrapeViaNameYear(movieName, year, proxyUrl)
	if err != nil {
		return "", err
	}
	find_str := `<meta property="og:title" content=`
	loc := strings.Index(content, find_str) + len(find_str)
	substring := content[loc : loc+100]
	subs := strings.Split(substring, ">")
	title := strings.ReplaceAll(subs[0], `"`, "")
	return title, nil
}

func GetGenres(movieName string, year int, proxyUrl string) ([]string, error) {
	content, err := scrapeViaNameYear(movieName, year, proxyUrl)
	if err != nil {
		return []string{}, err
	}
	schemaJson, err := getSchemaJsonLD(content)
	if err != nil {
		return []string{}, err
	}
	return schemaJson.Genre, nil
}

func GetActors(movieName string, year int, maxActors int, proxyUrl string) ([]RTSchemaPerson, error) {
	content, err := scrapeViaNameYear(movieName, year, proxyUrl)
	if err != nil {
		return []RTSchemaPerson{}, err
	}
	schemaJson, err := getSchemaJsonLD(content)
	if err != nil {
		return []RTSchemaPerson{}, err
	}
	var filtered []RTSchemaPerson
	if maxActors <= 0 {
		filtered = schemaJson.Actors
	} else {
		filtered = schemaJson.Actors[:maxActors]
	}
	return filtered, nil
}

func GetDirectors(movieName string, year int, maxDirectors int, proxyUrl string) ([]RTSchemaPerson, error) {
	content, err := scrapeViaNameYear(movieName, year, proxyUrl)
	if err != nil {
		return []RTSchemaPerson{}, err
	}
	schemaJson, err := getSchemaJsonLD(content)
	if err != nil {
		return []RTSchemaPerson{}, err
	}
	var filtered []RTSchemaPerson
	if maxDirectors <= 0 {
		filtered = schemaJson.Director
	} else {
		filtered = schemaJson.Director[:maxDirectors]
	}
	return filtered, nil
}
