package repo_ratings

type GetListRatingOutletRequestParam struct {
	Limit       int   `json:"limit"`
	Page        int   `json:"page"`
	WithMedia   bool  `json:"with_media"`
	WithComment bool  `json:"trans_status"`
	Ratings     []int `json:"ratings"`
}

type DataInfoRating struct {
	SummaryRating struct {
		All         int `json:"all"`
		WithMedia   int `json:"with_media"`
		WithComment int `json:"with_comment"`
		RateFive    int `json:"rate_five"`
		RateFour    int `json:"rate_four"`
		RateThree   int `json:"rate_three"`
		RateTwo     int `json:"rate_two"`
		RateOne     int `json:"rate_one"`
	} `json:"summary_rating"`
	CurrentPage int `json:"cur_page"`
	MaxPage     int `json:"max_page"`
	Limit       int `json:"limit"`
	TotalRecord int `json:"total_record"`
}
