package repo_ratings

import "context"

type RatingsRepository interface {
	SubmitRatingProduct(ctx context.Context, custId int, productId, comment, rate string, photoList []string) (errCode string, err error)
}
