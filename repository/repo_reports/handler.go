package repo_reports

import (
	"context"
	"fmt"
	"libra-internal/pkg/constants"
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
		pajak, ongkir, asuransi, nett_sales, hpp, gross_profit)
		values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE 
		tanggal = ?, tipe_transaksi = ?, ref = ?, no_pesanan=?, status=?, channel=?, nama_toko=?, pelanggan=?, 
		sub_total=?, diskon=?, diskon_lainnya=?, potongan_biaya=?, biaya_lain=?, termasuk_pajak=?,
		pajak=?, ongkir=?, asuransi=?, nett_sales=?, hpp=?, gross_profit=?`
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Errorf("error on row %v with invoice number:%v\r\n", (index + 1), args[3])
			return
		}
	}

	if err = tx.Commit(); err != nil {
		return
	}
	return
}

func (q *SqlRepository) UpdateNetProfit(ctx context.Context) (err error) {
	var (
		tempSales   []Sales
	)

	query := `select no_pesanan,gross_profit, channel
	from sales
	where potongan_marketplace is null
	limit 10000`

	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i Sales

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
		netProfit := val.GrossProfit - feePrice
		const query = `update sales set potongan_marketplace = ?, net_profit = ? where no_pesanan = ?`
		_, err = tx.ExecContext(ctx, query, feePrice, netProfit, val.NoPesanan)
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
