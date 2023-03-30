package ratings

import (
	"libra-internal/internal/api/products"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RatingProductRequest struct {
}

type GetListRatingOutletRequest struct {
	Limit       int   `json:"limit"`
	Page        int   `json:"page"`
	OutletId    int   `json:"outlet_id"`
	WithMedia   bool  `json:"with_media"`
	WithComment bool  `json:"with_comment"`
	Ratings     []int `json:"ratings"`
}

func (m *GetListRatingOutletRequest) Bind(r *http.Request) error {
	return m.ValidateGetListRatingOutletRequest()
}

func (m *GetListRatingOutletRequest) ValidateGetListRatingOutletRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.OutletId, validation.Required))
}

type GetListRatingOutletResponse struct {
	DataInfo          products.DataInfo `json:"info"`
	SummaryRatingData SummaryRating     `json:"summary_rating"`
	ListRatingData    []ListRating      `json:"list"`
}

type SummaryRating struct {
	AvgRating    string `json:"avg_rating"`
	SatisfyLevel string `json:"satisfy_level"`
	All          int    `json:"all"`
	WithMedia    int    `json:"with_media"`
	WithComment  int    `json:"with_comment"`
	RateFive     int    `json:"rate_five"`
	RateFour     int    `json:"rate_four"`
	RateThree    int    `json:"rate_three"`
	RateTwo      int    `json:"rate_two"`
	RateOne      int    `json:"rate_one"`
}

type ListRating struct {
	IdRating       int               `json:"id"`
	CustomerName   string            `json:"customer_name"`
	CustomerAvatar string            `json:"customer_avatar"`
	Rating         int               `json:"rating"`
	OutletName     string            `json:"outlet_name"`
	Comment        string            `json:"comment"`
	CreatedAt      time.Time         `json:"created_at"`
	ListImages     []ListRatingImage `json:"list_images"`
}

type ListRatingImage struct {
	ImageUrl string `json:"image_url"`
}

type GetListImageResponse struct {
	IdRating  int
	ImageName string
}
