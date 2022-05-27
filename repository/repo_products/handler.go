package repo_products

import (
	"context"
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

func (q *SqlRepository) GetListProducts(ctx context.Context, fp ProductsParamsTemp) (res []Products, totalData int, errCode string, err error) {
	var (
		args        = make([]interface{}, 0)
		whereParams = ""
		offsetNum   = (fp.Page - 1) * fp.Limit
		orderBy     = "a.NamaBarang"
		orderType   = "asc"
	)

	if len(fp.Name) > 0 {
		lowerName := strings.ToLower(fp.Name)
		whereParams += "and LOWER(a.NamaBarang) LIKE CONCAT('%', ?, '%') "
		args = append(args, lowerName)
	}

	if len(fp.Posisi) > 0 {
		whereParams += "and a.IDPosisi = ? "
		args = append(args, fp.Posisi)
	}

	if len(fp.MerkBan) > 0 {
		lowerName := strings.ToLower(fp.MerkBan)
		whereParams += "and LOWER(a.IDMerk) = ? "
		args = append(args, lowerName)
	}

	if fp.MinPrice > 0 {
		whereParams += "and a.HargaJualFinal >= ? "
		args = append(args, fp.MinPrice)
	}

	if fp.MaxPrice > 0 {
		whereParams += "and a.HargaJualFinal <= ? "
		args = append(args, fp.MaxPrice)
	}

	if len(fp.OrderBy) > 0 {
		switch fp.OrderBy {
		case "price":
			orderBy = "a.HargaJualFinal"
		case "time":
			orderBy = "a.CreateDate"
		default:
			orderBy = "a.NamaBarang"
		}
	}

	if len(fp.OrderType) > 0 && strings.ToLower(fp.OrderType) == "desc" {
		orderType = "desc"
	}

	queryRecords := `
	SELECT count(a.KodePLU)
	from tblmasterplu a
	inner join tblMerkBan b on a.IDMerk = b.IDMerk
	inner join tblmasterukuranban c on a.IDUkuran = c.IDUkuranBan
	inner join tblPosisiBan d on a.IDPosisi = d.IDPosisi where 1 = 1 ` + whereParams
	err = q.db.QueryRowContext(ctx, queryRecords, args...).Scan(&totalData)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	args = append(args, fp.Limit, offsetNum)

	query := `
	select a.KodePLU, a.NamaBarang, a.Disc, a.HargaJual, a.HargaJualFinal, c.Ukuran
	from tblmasterplu a
	inner join tblMerkBan b on a.IDMerk = b.IDMerk
	inner join tblmasterukuranban c on a.IDUkuran = c.IDUkuranBan
	inner join tblPosisiBan d on a.IDPosisi = d.IDPosisi
	where 1 = 1 ` + whereParams + ` 
	` + fmt.Sprintf("order by %v %v", orderBy, orderType) + `  limit ? offset ? `
	fmt.Println(query)

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		fmt.Println(err)
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
