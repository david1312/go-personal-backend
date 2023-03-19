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

func (q *SqlRepository) SubmitRatingOutlet(ctx context.Context, custId int, outletId, comment, rate string, photoList []string, invoiceID string) (errCode string, err error) {
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

	const queryFinishTransaction = `update tbltransaksihead set StatusPembayaran = 'LUNAS', StatusTransaksi = 'Berhasil' where NoFaktur = ?`
	_, err = tx.ExecContext(ctx, queryFinishTransaction, invoiceID)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	if err = tx.Commit(); err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) GetRatingSummary(ctx context.Context, outletId int) (res DataInfoRating, errCode string, err error) {
	var (
		args = make([]interface{}, 0)
	)
	for i := 1; i <= 9; i++ {
		args = append(args, outletId)
	}

	const query = `select (select count(id) from outlet_ratings where outlet_id = ?) as rate_all, 
	(select count(id) from outlet_ratings where outlet_id = ? and rating = 5) as rate_five,
	(select count(id) from outlet_ratings where outlet_id = ? and rating = 4) as rate_four,
	(select count(id) from outlet_ratings where outlet_id = ? and rating = 3) as rate_three,
	(select count(id) from outlet_ratings where outlet_id = ? and rating = 2) as rate_two,
	(select count(id) from outlet_ratings where outlet_id = ? and rating = 1) as rate_one,
	(select count(id) from outlet_ratings where outlet_id = ? and comment != '') as with_comment,
	(select count(distinct b.id_ratings) 
        from outlet_ratings a join outlet_ratings_img b on a.id = b.id_ratings 
        where  a.outlet_id = ?) as with_media,
	(Select sum(rating) from outlet_ratings where outlet_id = ?) as total_rating;`

	row := q.db.DB.QueryRowContext(ctx, query, args...)

	err = row.Scan(
		&res.All,
		&res.RateFive,
		&res.RateFour,
		&res.RateThree,
		&res.RateTwo,
		&res.RateOne,
		&res.WithComment,
		&res.WithMedia,
		&res.SumRating,
	)
	if err != nil {
		errCode = crashy.ErrCodeDataRead
		return
	}
	return
}

func (q *SqlRepository) GetListRatingOutlet(ctx context.Context, fp GetListRatingOutletRequestParam, outletId int) (res []GetListRatingResponse, totalData int, listRatingId []int, errCode string, err error) {
	var (
		args        = make([]interface{}, 0)
		whereParams = ""
		offsetNum   = (fp.Page - 1) * fp.Limit
	)

	if len(fp.Ratings) > 0 {
		inTotal := ""
		for _, v := range fp.Ratings {
			inTotal += "?,"
			args = append(args, v)
		}
		trimmed := inTotal[:len(inTotal)-1]
		whereParams += " and a.rating in(" + trimmed + ") "
	}

	if fp.WithComment {
		whereParams += " and a.comment != '' "
	}

	queryRecords := `
	select count(a.id) as totalData
	from outlet_ratings a
	join customers b on a.customer_id = b.id
	join outlets c on a.outlet_id = c.id
	where 1=1 ` + whereParams
	err = q.db.QueryRowContext(ctx, queryRecords, args...).Scan(&totalData)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	args = append(args, fp.Limit, offsetNum)
	query := `
	select a.id, b.name as customer_name, b.avatar, a.rating, c.name as outlet_name, a.comment, a.created_at
	from outlet_ratings a
	join customers b on a.customer_id = b.id
	join outlets c on a.outlet_id = c.id
	where 1=1 ` + whereParams + `order by a.created_at desc limit ? offset ? `

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i GetListRatingResponse

		if err = rows.Scan(
			&i.IdRating,
			&i.CustomerName,
			&i.CustomerAvatar,
			&i.Rating,
			&i.OutletName,
			&i.Comment,
			&i.CreatedAt,
		); err != nil {
			errCode = crashy.ErrCodeUnexpected
			return
		}
		listRatingId = append(listRatingId, i.IdRating)
		res = append(res, i)
	}
	if err = rows.Close(); err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	if err = rows.Err(); err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	return
}

func (q *SqlRepository) GetListRatingImage(ctx context.Context, listOutletId []int) (res []GetListImageResponse, errCode string, err error) {
	var (
		args        = make([]interface{}, 0)
		whereParams = ""
	)

	inTotal := ""
	if len(listOutletId) > 0 {
		for _, v := range listOutletId {
			inTotal += "?,"
			args = append(args, v)
		}
		trimmed := inTotal[:len(inTotal)-1]
		whereParams += " and id_ratings in(" + trimmed + ") "
	}

	query := `select id_ratings, image from outlet_ratings_img where 1=1 ` + whereParams
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i GetListImageResponse

		if err = rows.Scan(
			&i.IdRating,
			&i.ImageName,
		); err != nil {
			errCode = crashy.ErrCodeUnexpected
			return
		}
		res = append(res, i)
	}
	if err = rows.Close(); err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	if err = rows.Err(); err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}
