package repo_ratings

import "context"

type RatingsRepository interface {
	SubmitRatingProduct(ctx context.Context, custId int, productId, comment, rate string, photoList []string) (errCode string, err error)
	SubmitRatingOutlet(ctx context.Context, custId int, outletId, comment, rate string, photoList []string, invoiceID string) (errCode string, err error)

	GetRatingSummary(ctx context.Context, outletId int) (res DataInfoRating, errCode string, err error)
	GetListRatingOutlet(ctx context.Context, fp GetListRatingOutletRequestParam, outletId int) (res []GetListRatingResponse, totalData int, listRatingId []int, errCode string, err error)
	GetListRatingImage(ctx context.Context, listOutletId []int) (res []GetListImageResponse, errCode string, err error)
}
