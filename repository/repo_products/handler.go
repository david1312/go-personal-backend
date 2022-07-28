package repo_products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"semesta-ban/pkg/crashy"
	"strings"

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

func (q *SqlRepository) GetListProducts(ctx context.Context, fp ProductsParamsTemp, custId int) (res []Products, totalData int, errCode string, err error) {
	var (
		args        = make([]interface{}, 0)
		whereParams = ""
		offsetNum   = (fp.Page - 1) * fp.Limit
		orderBy     = "a.NamaBarang"
	)

	if len(fp.Name) > 0 {
		lowerName := strings.ToLower(fp.Name)
		whereParams += "and LOWER(a.NamaBarang) LIKE CONCAT('%', ?, '%') "
		args = append(args, lowerName)
	}

	if len(fp.MerkBan) > 0 {
		inTotal := ""
		for _, v := range fp.MerkBan {
			inTotal += "?,"
			args = append(args, v)
		}
		trimmed := inTotal[:len(inTotal)-1]
		whereParams += "and a.IDMerk in (" + trimmed + ") "
	}

	if fp.MinPrice > 0 {
		whereParams += "and a.HargaJualFinal >= ? "
		args = append(args, fp.MinPrice)
	}

	if fp.MaxPrice > 0 {
		whereParams += "and a.HargaJualFinal <= ? "
		args = append(args, fp.MaxPrice)
	}

	if len(fp.UkuranBan) > 0 && len(fp.ArrUkuran) == 0 {
		inTotal := ""
		for _, v := range fp.UkuranBan {
			inTotal += "?,"
			args = append(args, v)
		}
		trimmed := inTotal[:len(inTotal)-1]
		whereParams += "and a.IDUkuranRing in (" + trimmed + ") "
	}

	if len(fp.ArrUkuran) > 0 {
		inTotal := ""
		for _, v := range fp.ArrUkuran {
			inTotal += "?,"
			args = append(args, v)
		}
		trimmed := inTotal[:len(inTotal)-1]
		whereParams += "and a.IDUkuranRing in (" + trimmed + ") "

	}

	if len(fp.OrderBy) > 0 {
		switch fp.OrderBy {
		case "max_price":
			orderBy = "a.HargaJualFinal desc"
		case "min_price":
			orderBy = "a.HargaJualFinal asc"
		case "latest":
			orderBy = "a.CreateDate desc"
		default:
			orderBy = "a.NamaBarang asc"
		}
	}

	queryRecords := `
	SELECT count(a.kodeplu)
	from tblmasterplu a
	inner join tblmerkban b on a.IDMerk = b.IDMerk
	inner join tblposisiban d on a.IDPosisi = d.IDPosisi
	 where 1 = 1 ` + whereParams
	err = q.db.QueryRowContext(ctx, queryRecords, args...).Scan(&totalData)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	args = append(args, fp.Limit, offsetNum)

	query := `
	select a.KodePLU, a.NamaBarang, a.Disc, a.HargaJual, a.HargaJualFinal, a.IDUkuranRing, e.URL, a.JenisBan
	from tblmasterplu a
	inner join tblmerkban b on a.IDMerk = b.IDMerk
	inner join tblposisiban d on a.IDPosisi = d.IDPosisi
	left join tblurlgambar e on a.KodeBarang = e.KodeBarang and e.IsDisplay = true
	where 1 = 1 ` + whereParams + ` and e.DeletedAt IS NULL
	` + fmt.Sprintf("order by %v", orderBy) + `  limit ? offset ? `

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i Products

		if err = rows.Scan(
			&i.KodePLU,
			&i.NamaBarang,
			&i.Disc,
			&i.HargaJual,
			&i.HargaJualFinal,
			&i.NamaUkuran,
			&i.DisplayImage,
			&i.JenisBan,
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

func (q *SqlRepository) GetProductDetail(ctx context.Context, id, custId int) (res Products, errCode string, err error) {
	query := `
	select a.KodePLU, a.KodeBarang,  a.NamaBarang, a.Disc, a.HargaJual, a.HargaJualFinal, a.IDUkuranRing, a.JenisBan, d.Posisi, a.Deskripsi,
	(select exists(select x.product_id from wishlists x where x.customer_id = ? and x.product_id = a.KodePLU)) as isWishlist
	from tblmasterplu a
	inner join tblmerkban b on a.IDMerk = b.IDMerk
	inner join tblposisiban d on a.IDPosisi = d.IDPosisi
	where a.KodePLU = ? `

	row := q.db.DB.QueryRowContext(ctx, query, custId, id)

	err = row.Scan(
		&res.KodePLU,
		&res.KodeBarang,
		&res.NamaBarang,
		&res.Disc,
		&res.HargaJual,
		&res.HargaJualFinal,
		&res.NamaUkuran,
		&res.JenisBan,
		&res.NamaPosisi,
		&res.Deskripsi,
		&res.IsWishlist,
	)

	if err != nil {
		errCode = crashy.ErrCodeDataRead
		return
	}

	return
}

func (q *SqlRepository) GetProductImage(ctx context.Context, productCode string) (res []ProductImage, errCode string, err error) {
	query := `select URL, IsDisplay from tblurlgambar where KodeBarang = ? order by IsDisplay desc`

	rows, err := q.db.QueryContext(ctx, query, productCode)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i ProductImage

		if err = rows.Scan(
			&i.Url,
			&i.IsDisplay,
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

func (q *SqlRepository) GetCustomerId(ctx context.Context, uid string) (custId int, errCode string, err error) {
	const query = `SELECT id FROM customers where uid = ? AND deleted_at IS NULL`
	row := q.db.DB.QueryRowContext(ctx, query, uid)

	err = row.Scan(&custId)

	if err != nil && err != sql.ErrNoRows {
		errCode = crashy.ErrCodeUnexpected
		return
	} else if err != nil && err == sql.ErrNoRows {
		err = nil
		return
	}

	return
}

func (q *SqlRepository) WishlistAdd(ctx context.Context, custId, productId int) (errCode string, err error) {
	var checkExist int
	const query = `SELECT customer_id FROM wishlists where customer_id = ? AND product_id = ?`
	row := q.db.DB.QueryRowContext(ctx, query, custId, productId)

	err = row.Scan(&checkExist)
	if err != nil && err != sql.ErrNoRows {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	if checkExist > 0 {
		errCode = crashy.ErrExistsWishlist
		err = errors.New(crashy.ErrExistsWishlist)
		return
	}

	const queryInsert = `insert into wishlists (customer_id, product_id) VALUES (?, ?) `

	_, err = q.db.ExecContext(ctx, queryInsert, custId, productId)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) WishlistRemove(ctx context.Context, custId, productId int) (errCode string, err error) {
	const queryInsert = `delete from wishlists where customer_id = ? and product_id = ?`

	_, err = q.db.ExecContext(ctx, queryInsert, custId, productId)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) WishlistMe(ctx context.Context, custId int, fp ProductsParamsTemp) (res []Products, totalData int, errCode string, err error) {
	var (
		args      = make([]interface{}, 0)
		offsetNum = (fp.Page - 1) * fp.Limit
	)
	queryRecords := `
	select count(a.KodePLU)
	from tblmasterplu a
	inner join tblmerkban b on a.IDMerk = b.IDMerk
	inner join tblposisiban d on a.IDPosisi = d.IDPosisi
	left join tblurlgambar e on a.KodeBarang = e.KodeBarang and e.IsDisplay = true
        inner join wishlists f on f.product_id = a.KodePLU
		where f.customer_id = ?`
	args = append(args, custId)

	err = q.db.QueryRowContext(ctx, queryRecords, args...).Scan(&totalData)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	args = append(args, fp.Limit, offsetNum)

	query := `
	select a.KodePLU, a.NamaBarang, a.Disc, a.HargaJual, a.HargaJualFinal, a.IDUkuranRing, e.URL, a.JenisBan
	from tblmasterplu a
	inner join tblmerkban b on a.IDMerk = b.IDMerk
	inner join tblposisiban d on a.IDPosisi = d.IDPosisi
	left join tblurlgambar e on a.KodeBarang = e.KodeBarang and e.IsDisplay = true
        inner join wishlists f on f.product_id = a.KodePLU
		where f.customer_id = ?
		order by f.created_at desc limit ? offset ?`

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i Products

		if err = rows.Scan(
			&i.KodePLU,
			&i.NamaBarang,
			&i.Disc,
			&i.HargaJual,
			&i.HargaJualFinal,
			&i.NamaUkuran,
			&i.DisplayImage,
			&i.JenisBan,
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

func (q *SqlRepository) CartCheck(ctx context.Context, custUid string) (cartId int, errCode string, err error) {
	const query = `SELECT id FROM carts where customer_uid = ?`
	row := q.db.DB.QueryRowContext(ctx, query, custUid)
	err = row.Scan(&cartId)

	if err != nil && err != sql.ErrNoRows {
		errCode = crashy.ErrCodeUnexpected
		return
	} else if err != nil && err == sql.ErrNoRows {
		err = nil
		return
	}
	return
}

func (q *SqlRepository) CartAdd(ctx context.Context, custUid string) (cartId int, errCode string, err error) {
	const query = `INSERT INTO carts (customer_uid) VALUES (?)`
	res, err := q.db.ExecContext(ctx, query, custUid)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}

	cartId = int(lastId)

	return
}

func (q *SqlRepository) CartItemCheck(ctx context.Context, cartId, productId int) (cartItemId, qty int, errCode string, err error) {
	const query = `SELECT id, qty FROM carts_item where carts_id = ? and product_id = ?`
	row := q.db.DB.QueryRowContext(ctx, query, cartId, productId)
	err = row.Scan(&cartItemId, &qty)

	if err != nil && err != sql.ErrNoRows {
		errCode = crashy.ErrCodeUnexpected
		return
	} else if err != nil && err == sql.ErrNoRows {
		err = nil
		return
	}
	return
}

func (q *SqlRepository) CartItemAdd(ctx context.Context, cartId, productId int) (errCode string, err error) {
	const query = `INSERT INTO carts_item (carts_id, product_id, qty, is_selected) VALUES (?, ?, 1, true)`
	_, err = q.db.ExecContext(ctx, query, cartId, productId)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}

	return
}

func (q *SqlRepository) CartItemUpdate(ctx context.Context, cartItemId, qty int, isSelected bool) (errCode string, err error) {
	const queryUpdate = `update carts_item set qty =  ?, is_selected = ? where id = ?`
	_, err = q.db.ExecContext(ctx, queryUpdate, qty, isSelected, cartItemId)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) CartItemRemove(ctx context.Context, cartItemId int) (errCode string, err error) {
	const queryUpdate = `delete from carts_item where id = ?`
	_, err = q.db.ExecContext(ctx, queryUpdate, cartItemId)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) CartSelectDeselectAll(ctx context.Context, cartId int, isSelectAll bool) (errCode string, err error) {
	const queryUpdate = `update carts_item set is_selected = ? where carts_id = ?`
	_, err = q.db.ExecContext(ctx, queryUpdate, isSelectAll, cartId)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) CartMe(ctx context.Context, cartId int, fp ProductsParamsTemp) (res []Products, totalData int, errCode string, err error) {
	var (
		args      = make([]interface{}, 0)
		offsetNum = (fp.Page - 1) * fp.Limit
	)
	queryRecords := `
	select count(a.KodePLU)
	from tblmasterplu a
	inner join tblmerkban b on a.IDMerk = b.IDMerk
	inner join tblposisiban d on a.IDPosisi = d.IDPosisi
	left join tblurlgambar e on a.KodeBarang = e.KodeBarang and e.IsDisplay = true
        inner join carts_item f on f.product_id = a.KodePLU
		where f.carts_id = ?`
	args = append(args, cartId)

	err = q.db.QueryRowContext(ctx, queryRecords, args...).Scan(&totalData)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	args = append(args, fp.Limit, offsetNum)

	query := `
	select a.KodePLU, a.NamaBarang, a.Disc, a.HargaJual, a.HargaJualFinal, a.IDUkuranRing, e.URL, a.JenisBan, f.id, f.qty, f.is_selected
	from tblmasterplu a
	inner join tblmerkban b on a.IDMerk = b.IDMerk
	inner join tblposisiban d on a.IDPosisi = d.IDPosisi
	left join tblurlgambar e on a.KodeBarang = e.KodeBarang and e.IsDisplay = true
	inner join carts_item f on f.product_id = a.KodePLU
	where f.carts_id = ?
		order by f.created_at desc limit ? offset ?`

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i Products

		if err = rows.Scan(
			&i.KodePLU,
			&i.NamaBarang,
			&i.Disc,
			&i.HargaJual,
			&i.HargaJualFinal,
			&i.NamaUkuran,
			&i.DisplayImage,
			&i.JenisBan,
			&i.CartItemId,
			&i.CartItemQty,
			&i.CartItemIsSelected,
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
