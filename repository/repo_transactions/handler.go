package repo_transactions

import (
	"context"
	"database/sql"
	"errors"
	"semesta-ban/pkg/crashy"
	"semesta-ban/repository/repo_products"

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
func (q *SqlRepository) SubmitTransaction(ctx context.Context, fp SubmitTransactionsParam) (errCode string, err error) {
	var (
		argsCheckQty    = make([]interface{}, 0)
		mapStockProduct = make(map[int]int)
		whereQty        = ""
		tempProduct     = []repo_products.Products{}
		totalBayar      = 0
	)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer tx.Rollback()

	//check stock every products
	inTotal := ""
	for _, v := range fp.ListProduct {
		inTotal += "?,"
		argsCheckQty = append(argsCheckQty, v.ProductId)
		//mapping qty per product
		mapStockProduct[v.ProductId] = v.Qty
		totalBayar += int(v.Total)
	}
	trimmed := inTotal[:len(inTotal)-1]
	whereQty += " KodePLU in (" + trimmed + ") "
	queryCheckStock := `select KodePLU, StokAll from tblmasterplu where ` + whereQty
	rows, err := tx.QueryContext(ctx, queryCheckStock, argsCheckQty...)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {
		var i repo_products.Products
		if err = rows.Scan(
			&i.KodePLU,
			&i.StockAll,
		); err != nil {
			errCode = crashy.ErrCodeUnexpected
			return
		}
		tempProduct = append(tempProduct, i)
	}
	if err = rows.Close(); err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	if err = rows.Err(); err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	for _, v := range tempProduct {
		if mapStockProduct[int(v.KodePLU)] > v.StockAll {
			errCode = crashy.ErrInsufficientStock
			err = errors.New(crashy.ErrInsufficientStock)
			return
		}
		//update the stock
		_, err = tx.ExecContext(ctx, "update tblmasterplu set StokAll = StokAll - ? where KodePLU = ?",
			mapStockProduct[int(v.KodePLU)], v.KodePLU)
		if err != nil {
			errCode = crashy.ErrCodeUnexpected
			return
		}
	}

	//insert into tbl transcation head
	_, err = tx.ExecContext(ctx, `insert into tbltransaksihead (NoFaktur, Tagihan, TglTrans, IdOutlet, TipeTransaksi, MetodePembayaran, JadwalPemasangan, CustomerId, Catatan, Source ,CreateBy)
	values (?,?,?,?,?,?,?,?,?,?, 'Customer')`, fp.NoFaktur, totalBayar, fp.ScheduleDate, fp.IdOutlet, fp.TranType, fp.PaymentMethod, fp.ScheduleTime, fp.CustomerId, fp.Notes, fp.Source)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	//insert into tbl transaction detail
	for _, v := range fp.ListProduct {
		_, err = tx.ExecContext(ctx, "insert into tbltransaksidetail values (?,?,?,?,?)",
			fp.NoFaktur, v.ProductId, v.Qty, v.Price, v.Total)
		if err != nil {
			errCode = crashy.ErrCodeUnexpected
			return
		}
	}

	if err = tx.Commit(); err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	return
}

func (q *SqlRepository) GetLastTransactionId(ctx context.Context) (res, errCode string, err error) {
	const query = `select NoFaktur from tbltransaksihead order by CreateDate desc limit 1`
	row := q.db.DB.QueryRowContext(ctx, query)
	err = row.Scan(&res)
	if err != nil && err != sql.ErrNoRows {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	if err != nil && err == sql.ErrNoRows {
		err = nil
		return
	}
	return
}

func (q *SqlRepository) InquirySchedule(ctx context.Context, startDate, endDate string) (res []ScheduleCount, errCode string, err error) {
	const query = `select  TglTrans, JadwalPemasangan, count(NoFaktur) as totalOrder
	from tbltransaksihead
	where TglTrans BETWEEN ? AND ?
	group by JadwalPemasangan,TglTrans
	order by TglTrans asc, JadwalPemasangan asc`

	rows, err := q.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {
		var i ScheduleCount
		if err = rows.Scan(
			&i.ScheduleDate,
			&i.ScheduleTime,
			&i.OrderCount,
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
