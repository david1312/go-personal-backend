package repo_master_data

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

func (q *SqlRepository) GetListMerkBan(ctx context.Context) (res []MerkBan, errCode string, err error) {
	const query = `SELECT IdMerk, Merk, Icon from tblmerkban order by Ranking asc`
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i MerkBan

		if err = rows.Scan(
			&i.IdMerk,
			&i.Merk,
			&i.Icon,
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

func (q *SqlRepository) GetListUkuranBan(ctx context.Context) (res []UkuranRingBan, errCode string, err error) {
	const query = `select a.id, b.UkuranRing, a.id_ring_ban
					from tblbanukuranring a
					join tblmasterringban b on a.id_ring_ban = b.IDRing
					order by a.id_ring_ban asc`
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i UkuranRingBan

		if err = rows.Scan(
			&i.Id,
			&i.UkuranRing,
			&i.IdRingBan,
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

func (q *SqlRepository) GetListUkuranBanByBrandMotor(ctx context.Context, idBrandMotor []int) (res []UkuranRingBan, errCode string, err error) {
	var (
		args        = make([]interface{}, 0)
		whereParams = ""
		inTotal     = ""
	)

	for _, v := range idBrandMotor {
		inTotal += "?,"
		args = append(args, v)
	}
	trimmed := inTotal[:len(inTotal)-1]
	whereParams += " id_brand_motor in (" + trimmed + ") "

	query := `select distinct(id_ukuran_ring) as id from motor_x_size_ban where ` + whereParams + ` order by id_ukuran_ring asc;`
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i UkuranRingBan

		if err = rows.Scan(
			&i.Id,
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

func (q *SqlRepository) GetListUkuranBanByMotor(ctx context.Context, idMotor int) (res []UkuranRingBan, errCode string, err error) {
	const query = `select distinct(id_ukuran_ring) as id from motor_x_size_ban where id_motor = ? order by id_ukuran_ring asc;`
	rows, err := q.db.QueryContext(ctx, query, idMotor)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i UkuranRingBan

		if err = rows.Scan(
			&i.Id,
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

func (q *SqlRepository) GetListMerkMotor(ctx context.Context) (res []MerkMotor, errCode string, err error) {
	const query = `SELECT id, nama, icon from tblmerkmotor order by id asc`
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i MerkMotor

		if err = rows.Scan(
			&i.Id,
			&i.Nama,
			&i.Icon,
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

func (q *SqlRepository) GetListMotorByBrand(ctx context.Context, idBrandMotor int) (res []Motor, errCode string, err error) {

	const query = `select a.id,a.nama,a.icon, b.nama as category_name
					from tblmotor a join tblkategorimotor b on a.id_kategori_motor = b.id
					where a.id_merk_motor = ?`
	rows, err := q.db.QueryContext(ctx, query, idBrandMotor)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i Motor

		if err = rows.Scan(
			&i.Id,
			&i.Name,
			&i.Icon,
			&i.CategoryName,
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

func (q *SqlRepository) GetListPaymentMethod(ctx context.Context) (res []PaymentMethod, errCode string, err error) {
	const query = `select a.id, a.description, a.is_default, a.icon,b.name as category
	from payment_method a
	join payment_category b on a.id_payment_category = b.id
	order by b.id asc, a.is_default desc`
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i PaymentMethod

		if err = rows.Scan(
			&i.Id,
			&i.Description,
			&i.IsDefault,
			&i.Icon,
			&i.CategoryName,
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
