package repo_ratings

import (
	"context"
	"semesta-ban/pkg/crashy"

	"github.com/jmoiron/sqlx"
)

type SqlRepository struct {
	db *sqlx.DB
}

func NewSqlRepository(db *sqlx.DB) *SqlRepository {
	return &SqlRepository{
		db: db,
	}
}

func (q *SqlRepository) SubmitRatingProduct(ctx context.Context, custId int, productId, comment, rate string, photoList []string) (errCode string, err error) {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `insert into product_ratings (customer_id, product_id, comment, rating)
	values (?,?,?,?)`, custId, productId, comment, rate)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	//insert list image if included
	if len(photoList) > 0 {
		for _, v := range photoList {
			_, err = tx.ExecContext(ctx, "insert into product_ratings_img (id_ratings, image) values (?,?)",
				lastInsertId, v)
			if err != nil {
				errCode = crashy.ErrCodeUnexpected
				return
			}
		}
	}

	if err = tx.Commit(); err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}
