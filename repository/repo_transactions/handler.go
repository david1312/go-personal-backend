package repo_transactions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"semesta-ban/pkg/crashy"
	"semesta-ban/repository/repo_products"
	"time"

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
	paymentDue := time.Now().Add(24 * time.Hour) // updated later +24 from response midtrans api
	_, err = tx.ExecContext(ctx, `insert into tbltransaksihead (NoFaktur, Tagihan, TglTrans, IdOutlet, TipeTransaksi, MetodePembayaran, JadwalPemasangan, CustomerId, Catatan, Source ,CreateBy, PaymentDue)
	values (?,?,?,?,?,?,?,?,?,?, 'Customer', ?)`, fp.NoFaktur, totalBayar, fp.ScheduleDate, fp.IdOutlet, fp.TranType, fp.PaymentMethod, fp.ScheduleTime, fp.CustomerId, fp.Notes, fp.Source, paymentDue)
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

func (q *SqlRepository) GetHistoryTransaction(ctx context.Context, fp GetListTransactionsParam) (res []Transactions, totalData int, listInvoice []string, errCode string, err error) {
	var (
		args        = make([]interface{}, 0)
		whereParams = ""
		offsetNum   = (fp.Page - 1) * fp.Limit
		orderBy     = "CreateDate desc "
	)

	whereParams += " and a.CustomerId = ? "
	args = append(args, fp.CustomerId)

	if len(fp.StatusTransactions) > 0 {
		inTotal := ""
		for _, v := range fp.StatusTransactions {
			inTotal += "?,"
			args = append(args, v)
		}
		trimmed := inTotal[:len(inTotal)-1]
		whereParams += " and a.StatusTransaksi in(" + trimmed + ") "
	}

	queryRecords := `
	select count(a.NoFaktur)
	from tbltransaksihead a
	join payment_method b on a.MetodePembayaran = b.id
	where 1=1 ` + whereParams
	err = q.db.QueryRowContext(ctx, queryRecords, args...).Scan(&totalData)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	args = append(args, fp.Limit, offsetNum)

	query := `
	select a.NoFaktur, a.StatusTransaksi, a.Tagihan,a.CreateDate, b.description as payment_desc, b.icon, a.PaymentDue
	from tbltransaksihead a
	join payment_method b on a.MetodePembayaran = b.id
	where 1=1` + whereParams + `
	` + fmt.Sprintf("order by %v", orderBy) + `  limit ? offset ? `

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i Transactions

		if err = rows.Scan(
			&i.InvoiceId,
			&i.Status,
			&i.TotalAmount,
			&i.CreatedAt,
			&i.PaymentMethodDesc,
			&i.PaymentMethodIcon,
			&i.PaymentDue,
		); err != nil {
			errCode = crashy.ErrCodeUnexpected
			return
		}
		listInvoice = append(listInvoice, i.InvoiceId)
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

func (q *SqlRepository) GetProductByInvoices(ctx context.Context, listInvoiceId []string) (res []ProductsData, errCode string, err error) {
	var (
		args        = make([]interface{}, 0)
		whereParams = ""
	)

	inTotal := ""
	for _, v := range listInvoiceId {
		inTotal += "?,"
		args = append(args, v)
	}
	trimmed := inTotal[:len(inTotal)-1]
	whereParams += " a.NoFaktur in(" + trimmed + ") "

	query := `select a.NoFaktur, b.NamaBarang, b.IdUkuranRing, a.HargaSatuan, b.Deskripsi, c.Url, a.QtyItem, a.Total, a.IdBarang
	from tbltransaksidetail a
	join tblmasterplu b on a.IdBarang = b.KodePlu
	left join tblurlgambar c on b.KodeBarang = c.KodeBarang and c.IsDisplay = true
	where ` + whereParams

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i ProductsData

		if err = rows.Scan(
			&i.InvoiceId,
			&i.NamaBarang,
			&i.NamaUkuran,
			&i.Harga,
			&i.Deskripsi,
			&i.DisplayImage,
			&i.Qty,
			&i.HargaTotal,
			&i.KodePLU,
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
