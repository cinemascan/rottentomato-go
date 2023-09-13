package rotten_tomato

import (
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func testScrapeActors(t *testing.T, proxyUrl string) {
	expected := []RTSchemaPerson{
		{
			Name:  "Keanu Reeves",
			Url:   "https://www.rottentomatoes.com/celebrity/keanu_reeves",
			Image: "https://resizing.flixster.com/YARxkSH8c59kDC2pA87rGSQ8uX0=/100x120/v2/https://flxt.tmsimg.com/assets/1443_v9_bc.jpg",
		},
		{
			Name:  "Laurence Fishburne",
			Url:   "https://www.rottentomatoes.com/celebrity/larry_fishburne",
			Image: "https://resizing.flixster.com/K9EdHynqolK6JE8ZB_9-fNb0KhA=/100x120/v2/https://flxt.tmsimg.com/assets/71229_v9_bb.jpg",
		},
		{
			Name:  "Carrie-Anne Moss",
			Url:   "https://www.rottentomatoes.com/celebrity/carrie_anne_moss",
			Image: "https://resizing.flixster.com/o1J5kMouS0pTlmH7Zp4NEpBuJD0=/100x120/v2/https://flxt.tmsimg.com/assets/78172_v9_bb.jpg",
		},
		{
			Name:  "Hugo Weaving",
			Url:   "https://www.rottentomatoes.com/celebrity/hugo_weaving",
			Image: "https://resizing.flixster.com/l4YJKIcNoGzuyVJXq_egJIc25Lw=/100x120/v2/https://flxt.tmsimg.com/assets/27163_v9_bb.jpg",
		},
		{
			Name:  "Joe Pantoliano",
			Url:   "https://www.rottentomatoes.com/celebrity/joe_pantoliano",
			Image: "https://resizing.flixster.com/G0MMX7KZ0DZG0rxtg6UwxdK-sMM=/100x120/v2/https://flxt.tmsimg.com/assets/32287_v9_bb.jpg",
		},
	}
	actors, err := GetActors("the matrix", 1999, 5, proxyUrl)
	if err != nil {
		t.Fatalf("error occurred while scraping `the matrix` actors: %+v", err)
	}
	assert.Equal(t, actors, expected)
}

func testScrapeDirectors(t *testing.T, proxyUrl string) {
	expected := []RTSchemaPerson{
		{
			Name:  "Lilly Wachowski",
			Url:   "https://www.rottentomatoes.com/celebrity/lilly_wachowski",
			Image: "https://resizing.flixster.com/tj_RivcrCUlgfuV8xJb6koyPRYo=/100x120/v2/https://flxt.tmsimg.com/assets/150670_v9_ba.jpg",
		},
		{
			Name:  "Lana Wachowski",
			Url:   "https://www.rottentomatoes.com/celebrity/lana_wachowski",
			Image: "https://resizing.flixster.com/nIZRovwZGWwrmpbqjmsUrFa9HgI=/100x120/v2/https://flxt.tmsimg.com/assets/150673_v9_ba.jpg",
		},
	}
	directors, err := GetDirectors("the matrix", 1999, 0, proxyUrl)
	if err != nil {
		t.Fatalf("error occurred while scraping `the matrix` directors: %+v", err)
	}
	assert.Equal(t, directors, expected)
}

func TestScrapeMovieInfo(t *testing.T) {
	expectedMoveInfos := []RTMovieInfo{
		{
			Title:   "The Matrix",
			Year:    1999,
			Rating:  "R",
			Runtime: "2h 16m",
			Genres:  []string{"Sci-fi", "Action"},
			AudienceScore: RTScore{
				AverageRating:     "3.6",
				BandedRatingCount: "250,000+",
				LikedCount:        142778,
				NotLikedCount:     24632,
				RatingCount:       33324202,
				ReviewCount:       1307885,
				State:             "upright",
				Value:             85,
			},
			TomatometerScore: RTScore{
				AverageRating:     "7.70",
				BandedRatingCount: "",
				LikedCount:        171,
				NotLikedCount:     36,
				RatingCount:       207,
				ReviewCount:       207,
				State:             "certified-fresh",
				Value:             83,
			},
		},
	}

	for _, expectedMoveInfo := range expectedMoveInfos {
		movieInfo, err := GetMovieInfo(expectedMoveInfo.Title, expectedMoveInfo.Year, "")
		if err != nil {
			log.Fatalf("Error: %v", err.Error())
		}
		assert.Equal(t, cmp.Equal(*movieInfo, expectedMoveInfo), true)
	}
}
