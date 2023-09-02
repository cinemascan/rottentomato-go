package rotten_tomato

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type SearchListing struct {
	Title          string
	HasTomatometer bool
	IsMovie        bool
	Year           int
	Url            string
}

func FromHtml(htmlSnippet string) (*SearchListing, error) {
	tomato_query := "tomatometerscore"
	tomato_loc := strings.Index(htmlSnippet, tomato_query) + len(tomato_query)
	tomato_snippet := htmlSnippet[tomato_loc : tomato_loc+5]
	meter := strings.Split(tomato_snippet, `"`)[1]
	hasTomatometer := false
	if len(meter) > 0 {
		hasTomatometer = true
	}

	titleReg := regexp.MustCompile(`alt="(.*?)"`)
	rawTitles := titleReg.FindAllString(htmlSnippet, -1)
	titles := []string{}
	for _, rawTitle := range rawTitles {
		parsed := strings.ReplaceAll(strings.Split(rawTitle, "alt=")[1], `"`, "")
		titles = append(titles, parsed)
	}
	if len(titles) == 0 {
		return nil, fmt.Errorf("no valid RT titles found")
	}
	title := titles[0]

	// get url
	hrefReg := regexp.MustCompile(`a href="(.*?)"`)
	hrefs := hrefReg.FindAllString(htmlSnippet, -1)
	urls := []string{}
	for _, href := range hrefs {
		url := strings.ReplaceAll(strings.Split(href, "href=")[1], `"`, "")
		urls = append(urls, url)
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("no valid RT urls found")
	}
	url := urls[0]

	// get year
	yearReg := regexp.MustCompile(`releaseyear="(.*?)"`)
	scrapedYears := yearReg.FindAllString(htmlSnippet, -1)
	var years []int
	for _, scrapedYear := range scrapedYears {
		cleaned := strings.ReplaceAll(strings.Split(scrapedYear, "releaseyear=")[1], `"`, "")
		if cleaned == "" {
			continue
		}
		year, err := strconv.Atoi(cleaned)
		if err != nil {
			log.Printf(`error occurred while converting scraped html release year to int %v`, err)
			continue
		}
		years = append(years, year)
	}
	if len(years) == 0 {
		return nil, fmt.Errorf("no valid RT years found")
	}
	year := years[0]

	// get movie type
	isMovie := false
	movieIdx := strings.Index(url, `/m/`)
	if movieIdx != -1 {
		isMovie = true
	}

	return &SearchListing{
		Title:          title,
		HasTomatometer: hasTomatometer,
		IsMovie:        isMovie,
		Url:            url,
		Year:           year,
	}, nil
}
