package repo_ratings

import (
	"context"
	"fmt"
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

func (q *SqlRepository) SubmitRatingOutlet(ctx context.Context, custId int, outletId, comment, rate string, photoList []string) (errCode string, err error) {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `insert into outlet_ratings (customer_id, outlet_id, comment, rating)
	values (?,?,?,?)`, custId, outletId, comment, rate)
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
			_, err = tx.ExecContext(ctx, "insert into outlet_ratings_img (id_ratings, image) values (?,?)",
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

func (q *SqlRepository) GetListRatingOutlet(ctx context.Context, fp GetListRatingOutletRequestParam, outletId int) (res DataInfoRating, totalData int, errCode string, err error){

	const query = `select (select count(id) from outlet_ratings) as rate_all, 
	(select count(id) from outlet_ratings where rating = 5) as rate_five,
	(select count(id) from outlet_ratings where rating = 4) as rate_four,
	(select count(id) from outlet_ratings where rating = 3) as rate_three,
	(select count(id) from outlet_ratings where rating = 2) as rate_two,
	(select count(id) from outlet_ratings where rating = 1) as rate_one,
	(select count(id) from outlet_ratings where comment != '') as with_comment,
	(select count(id_ratings) from outlet_ratings_img group by id_ratings) as with_media,
	(Select sum(rating) from outlet_ratings) as total_rating`

	fmt.Println(query)
	return
}
