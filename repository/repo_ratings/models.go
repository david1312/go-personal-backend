package repo_ratings

import (
	"database/sql"
	"time"
)

type GetListRatingOutletRequestParam struct {
	Limit       int   `json:"limit"`
	Page        int   `json:"page"`
	WithMedia   bool  `json:"with_media"`
	WithComment bool  `json:"trans_status"`
	Ratings     []int `json:"ratings"`
}

type DataInfoRating struct {
	All         int `json:"all"`
	WithMedia   int `json:"with_media"`
	WithComment int `json:"with_comment"`
	RateFive    int `json:"rate_five"`
	RateFour    int `json:"rate_four"`
	RateThree   int `json:"rate_three"`
	RateTwo     int `json:"rate_two"`
	RateOne     int `json:"rate_one"`
	SumRating   int
}

type GetListRatingResponse struct {
	IdRating       int
	CustomerName   string
	CustomerAvatar sql.NullString
	Rating         int
	OutletName     string
	Comment        string
	CreatedAt      time.Time
}

type GetListImageResponse struct {
	IdRating  int
	ImageName string
}
