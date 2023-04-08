package repo_reports

import (
	"context"
	"fmt"
	"libra-internal/internal/models"
	"libra-internal/pkg/constants"
	"libra-internal/pkg/crashy"
	"libra-internal/pkg/helper"
	"libra-internal/pkg/log"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slices"
)

type SqlRepository struct {
	db *sqlx.DB
}

func NewSqlRepository(db *sqlx.DB) *SqlRepository {
	return &SqlRepository{
		db: db,
	}
}

func (q *SqlRepository) SyncUpSales(ctx context.Context, fileName, dir string) (err error) {
	//prepare transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()

	//open excel file
	f, err := excelize.OpenFile(fmt.Sprintf("%v%v", dir, fileName))
	if err != nil {
		return
	}
	rows := f.GetRows("Data")
	style, _ := f.NewStyle(`{"number_format":22}`)
	f.SetCellStyle("Data", "A2", fmt.Sprintf("A%v", len(rows)), style)

	// Get value from cell by given worksheet name and cell reference.
	// fmt.Println(len(rows))
	counter := 0
	rowsData := f.GetRows("Data")
	for index, row := range rowsData {
		var args = make([]interface{}, 0)
		counter++
		if counter == 1 {
			continue
		}

		var dateColumn = 0
		for i, colCell := range row {
			// a := slices.Contains()
			isNumberColumn := slices.Contains(constants.TypeDataNumberReportsSales, i)
			if i == dateColumn {
				args = append(args, helper.ConvertDateTimeReportExcel(colCell))
			} else if isNumberColumn {
				args = append(args, helper.GetDefaultNumberDBVal(colCell))
			} else {
				args = append(args, colCell)
			}
		}
		for _, valArg := range args {
			args = append(args, valArg)
		}

		//insert data to db
		const query = `insert into sales 
		(tanggal, tipe_transaksi, ref, no_pesanan, status, channel, nama_toko, pelanggan, 
		sub_total, diskon, diskon_lainnya, potongan_biaya, biaya_lain, termasuk_pajak,
		pajak, ongkir, asuransi, nett_sales, hpp, gross_profit, is_calculated_profit)
		values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,false) ON DUPLICATE KEY UPDATE 
		tanggal = ?, tipe_transaksi = ?, ref = ?, no_pesanan=?, status=?, channel=?, nama_toko=?, pelanggan=?, 
		sub_total=?, diskon=?, diskon_lainnya=?, potongan_biaya=?, biaya_lain=?, termasuk_pajak=?,
		pajak=?, ongkir=?, asuransi=?, nett_sales=?, hpp=?, gross_profit=?, is_calculated_profit = false`
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Errorf("error on row %v with invoice number:%v\r\n", (index + 1), args[3])
			return
		}
	}

	if err = tx.Commit(); err != nil {
		return
	}

	err = q.UpdateNetProfit(ctx, false)
	return
}

func (q *SqlRepository) UpdateNetProfit(ctx context.Context, limit bool) (err error) {
	var (
		tempSales []SalesModel
	)

	query := `select no_pesanan,gross_profit, channel
	from sales
	where is_calculated_profit = false
	limit 10000`

	if !limit {
		query = `select no_pesanan,gross_profit, channel
		from sales
		where is_calculated_profit = false`
	}

	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i SalesModel

		if err = rows.Scan(
			&i.NoPesanan,
			&i.GrossProfit,
			&i.Channel,
		); err != nil {
			return
		}
		tempSales = append(tempSales, i)
	}
	if err = rows.Close(); err != nil {
		return
	}
	if err = rows.Err(); err != nil {
		return
	}

	if len(tempSales) == 0 {
		log.Infof("all sales profit already calculated good job!\r\n")
		return
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()

	for _, val := range tempSales {
		feePrice := helper.CalculateFeeMarketPlace(val.GrossProfit, val.Channel)
		fixedFee := feePrice
		if feePrice < 0 {
			fixedFee = feePrice * -1
		}
		netProfit := val.GrossProfit - fixedFee
		const query = `update sales set potongan_marketplace = ?, net_profit = ?, is_calculated_profit = true where no_pesanan = ?`
		_, err = tx.ExecContext(ctx, query, fixedFee, netProfit, val.NoPesanan)
		if err != nil {
			log.Errorf("error calculate profit with invoice number:%v, due to :%v\r\n", val.NoPesanan, err.Error())
			return
		}
	}

	log.Infof("success calculate profit total data updated :%v\r\n", len(tempSales))
	if err = tx.Commit(); err != nil {
		return
	}
	return
}

func (q *SqlRepository) GetDetailInvoice(ctx context.Context, noPesanan string) (err error) {
	return
}

func (q *SqlRepository) GetAllSalesReport(ctx context.Context, params models.GetAllSalesRequest) (res []SalesModel, pageData models.Pagination, summary models.SummarySales, errCode string, err error) {
	var (
		args        = make([]interface{}, 0)
		whereParams = ""
		offsetNum   = (params.Page - 1) * params.Limit
		orderBy     = " order by a.tanggal asc "
		totalData = 0
	)

	whereParams += "and a.tanggal between ? "
	args = append(args, fmt.Sprintf("%v 00:00:00", params.StartDate))

	whereParams += "and ? "
	args = append(args, fmt.Sprintf("%v 23:59:59", params.EndDate))

	if len(params.NoPesanan) > 0 {
		args = nil
		whereParams = "and a.no_pesanan = ? "
		args = append(args, params.NoPesanan)
	}

	querySummary := `select  COUNT(a.id),
	COALESCE(SUM(nett_sales),0) AS total_sales,
	COALESCE(Sum( Case 
				When status != 'FAILED' AND status != 'RETURNED' then gross_profit
				Else 0 End ),0) as total_gross,
	COALESCE(Sum(Case 
				When status != 'FAILED' AND status != 'RETURNED' then potongan_marketplace
				Else 0 End ),0) as total_fee,
	COALESCE(Sum(Case 
				When status != 'FAILED' AND status != 'RETURNED' then net_profit
				Else 0 End ),0) as total_net_profit
	
	from sales a
	where 1 = 1 ` + whereParams

	err = q.db.QueryRowContext(ctx, querySummary, args...).Scan(&totalData, &summary.TotalNettSales, &summary.TotalGross, &summary.TotalPotonganMarketplace, &summary.TotalNetProfit)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	args = append(args, params.Limit, offsetNum)

	queryData := `select a.id, a.no_pesanan, a.tanggal, a.status, a.channel, a.nett_sales, a.gross_profit, a.potongan_marketplace, a.net_profit
	from sales a
	where 1=1 ` + whereParams + orderBy + ` limit ? offset ? `
	rows, err := q.db.QueryContext(ctx, queryData, args...)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i SalesModel

		if err = rows.Scan(
			&i.ID,
			&i.NoPesanan,
			&i.Tanggal,
			&i.Status,
			&i.Channel,
			&i.NettSales,
			&i.GrossProfit,
			&i.PotonganMarketPlace,
			&i.NetProfit,
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


	pageData = helper.CalculatePaginationData(params.Page, params.Limit, totalData)

	return
}
