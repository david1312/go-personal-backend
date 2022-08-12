package repo_ratings

import "context"

type RatingsRepository interface {
	SubmitRatingProduct(ctx context.Context, custId int, productId, comment, rate string, photoList []string) (errCode string, err error)
	SubmitRatingOutlet(ctx context.Context, custId int, outletId, comment, rate string, photoList []string) (errCode string, err error)

	GetListRatingOutlet(ctx context.Context, fp GetListRatingOutletRequestParam, outletId int) (res DataInfoRating, totalData int, errCode string, err error)
}
