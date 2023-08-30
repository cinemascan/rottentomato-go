package rotten_tomato

type RTSchemaPerson struct {
	Name  string `json:"name"`
	Url   string `json:"sameAs"`
	Image string `json:"image"`
}

type RTSchemaAggregateRating struct {
	Type        string `json:"@type"`
	BestRating  string `json:"bestRating"`
	Description string `json:"description"`
	Name        string `json:"name"`
	RatingCount int    `json:"ratingCount"`
	RatingValue string `json:"ratingValue"`
	ReviewCount int    `json:"reviewCount"`
	WorstRating string `json:"worstRating"`
}

type RTSchemaCompany struct {
	Type string `json:"@type"`
	Name string `json:"name"`
}

type RTSchemaJson struct {
	Context           string                  `json:"@context"`
	Type              string                  `json:"@type"`
	Actors            []RTSchemaPerson        `json:"actors"`
	AggregateRating   RTSchemaAggregateRating `json:"aggregateRating"`
	Author            []RTSchemaPerson        `json:"author"`
	Character         []string                `json:"character"`
	ContentRating     string                  `json:"contentRating"`
	DateCreated       string                  `json:"dateCreated"`
	DateModified      string                  `json:"dateModified"`
	Director          []RTSchemaPerson        `json:"director"`
	Genre             []string                `json:"genre"`
	Image             string                  `json:"image"`
	Name              string                  `json:"name"`
	ProductionCompany RTSchemaCompany         `json:"productionCompany"`
	Url               string                  `json:"url"`
}

type RTScore struct {
	AverageRating     string `json:"averageRating"`
	BandedRatingCount string `json:"bandedRatingCount"`
	LikedCount        int    `json:"likedCount"`
	NotLikedCount     int    `json:"notLikedCount"`
	RatingCount       int    `json:"ratingCount"`
	ReviewCount       int    `json:"reviewCount"`
	State             string `json:"state"`
	Value             int    `json:"value"`
}

type RTScoreboard struct {
	AudienceCountHref string  `json:"audienceCountHref"`
	AudienceScore     RTScore `json:"audienceScore"`
	Rating            string  `json:"rating"`
	TomatometerScore  RTScore `json:"tomatometerScore"`
	Title             string  `json:"title"`
	Info              string  `json:"info"`
}

type RTMovieInfo struct {
	AudienceScore    RTScore  `json:"audienceScore"`
	Rating           string   `json:"rating"`
	TomatometerScore RTScore  `json:"tomatometerScore"`
	Title            string   `json:"title"`
	Year             int      `json:"year"`
	Runtime          string   `json:"runtime"`
	Genres           []string `json:"genres"`
}

type RTScoreDetails struct {
	MediaType       string       `json:"mediaType"`
	PrimaryImageUrl string       `json:"primaryImageUrl"`
	Scoreboard      RTScoreboard `json:"scoreboard"`
}
