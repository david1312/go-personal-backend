package repo_master_data

import (
	"context"
	"semesta-ban/pkg/constants"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/helper"
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
	const query = `select a.id, b.UkuranRing, b.ranking
					from tblbanukuranring a
					join tblmasterringban b on a.id_ring_ban = b.IDRing
					order by b.ranking asc`
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
			&i.Ranking,
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

func (q *SqlRepository) GetListTopRankpMotor(ctx context.Context) (res []Motor, errCode string, err error) {
	const query = `select a.id,a.nama,a.icon from tblmotor a order by ranking asc limit 8`
	rows, err := q.db.QueryContext(ctx, query)
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

func (q *SqlRepository) AddBrandMotor(ctx context.Context, name, icon string) (errCode string, err error) {
	const queryInsert = `insert into tblmerkmotor (nama, icon) VALUES (?, ?) `

	_, err = q.db.ExecContext(ctx, queryInsert, name, icon)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) CheckBrandMotorUsed(ctx context.Context, idMotor int) (exists bool, errCode string, err error) {
	const query = `SELECT EXISTS(SELECT * FROM tblmotor WHERE id_merk_motor = ?)`
	row := q.db.DB.QueryRowContext(ctx, query, idMotor)
	err = row.Scan(&exists)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) RemoveBrandMotor(ctx context.Context, idMotor int, uploadPath, dirFile string) (errCode string, err error) {
	err = q.removeBrandMotor(ctx, idMotor, uploadPath, dirFile)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	const query = `delete from tblmerkmotor where id = ? `

	_, err = q.db.ExecContext(ctx, query, idMotor)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) UpdateBrandMotor(ctx context.Context, idMotor int, name string) (errCode string, err error) {
	const query = `update tblmerkmotor set nama = ? where id = ? `
	_, err = q.db.ExecContext(ctx, query, name, idMotor)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) UpdateBrandMotorImage(ctx context.Context, idMotor int, fileName, uploadPath, dirFile string) (errCode string, err error) {
	_ = q.removeBrandMotor(ctx, idMotor, uploadPath, dirFile)

	const query = `update tblmerkmotor set icon = ? where id = ? `
	_, err = q.db.ExecContext(ctx, query, fileName, idMotor)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) removeBrandMotor(ctx context.Context, id int, uploadPath, dirFile string) error {
	var temp = Motor{}
	const querySelect = `SELECT icon from tblmerkmotor where id = ?`
	row := q.db.DB.QueryRowContext(ctx, querySelect, id)

	err := row.Scan(
		&temp.Icon,
	)
	if err != nil {
		return err
	}
	if temp.Icon == constants.DefaultImgPng {
		return nil
	}
	helper.RemoveFile(temp.Icon, uploadPath, dirFile)
	return nil
}

func (q *SqlRepository) CheckBrandMotorExist(ctx context.Context, idMotor int) (exists bool, errCode string, err error) {
	const query = `SELECT EXISTS(SELECT id FROM tblmerkmotor WHERE id = ?)`
	row := q.db.DB.QueryRowContext(ctx, query, idMotor)
	err = row.Scan(&exists)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) AddTireBrand(ctx context.Context, id, name, icon, ranking string) (errCode string, err error) {
	const queryInsert = `insert into tblmerkban (IDMerk, Merk, Icon, Ranking) VALUES (?, ?, ?, ?) `

	_, err = q.db.ExecContext(ctx, queryInsert, id, name, icon, ranking)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) CheckTireBrandUsed(ctx context.Context, idMerkBan string) (exists bool, errCode string, err error) {
	const query = `select exists(select KodePLU from tblmasterplu where IDMerk = ? limit 1)`
	row := q.db.DB.QueryRowContext(ctx, query, idMerkBan)
	err = row.Scan(&exists)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) RemoveTireBrand(ctx context.Context, idMerkBan, uploadPath, dirFile string) (errCode string, err error) {
	_ = q.removeTireBrand(ctx, idMerkBan, uploadPath, dirFile)

	const query = `delete from tblmerkban where IDMerk = ? `
	_, err = q.db.ExecContext(ctx, query, idMerkBan)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) removeTireBrand(ctx context.Context, id, uploadPath, dirFile string) error {
	var temp = Motor{}
	const querySelect = `SELECT Icon from tblmerkban where IDMerk = ?`
	row := q.db.DB.QueryRowContext(ctx, querySelect, id)

	err := row.Scan(
		&temp.Icon,
	)
	if err != nil {
		return err
	}
	if temp.Icon == constants.DefaultImgPng {
		return nil
	}
	helper.RemoveFile(temp.Icon, uploadPath, dirFile)
	return nil
}

func (q *SqlRepository) UpdateTireBrand(ctx context.Context, idMerkBan, name string, ranking int) (errCode string, err error) {
	const query = `update tblmerkban set Merk = ?, Ranking = ? where IDMerk = ? `
	_, err = q.db.ExecContext(ctx, query, name, ranking, idMerkBan)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) CheckTireBrandExist(ctx context.Context, idMerkBan string) (exists bool, errCode string, err error) {
	const query = `select exists(select IDMerk from tblmerkban where IDMerk = ?)`
	row := q.db.DB.QueryRowContext(ctx, query, idMerkBan)
	err = row.Scan(&exists)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) UpdateTireBrandImage(ctx context.Context, idMerkBan, fileName, uploadPath, dirFile string) (errCode string, err error) {
	_ = q.removeTireBrand(ctx, idMerkBan, uploadPath, dirFile)

	const query = `update tblmerkban set Icon = ? where IDMerk = ? `
	_, err = q.db.ExecContext(ctx, query, fileName, idMerkBan)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) GetListMotor(ctx context.Context, fp ListMotorRequestRepo) (res []MotorMD, totalData int, errCode string, err error) {
	var (
		args        = make([]interface{}, 0)
		whereParams = ""
		offsetNum   = (fp.Page - 1) * fp.Limit
	)

	if len(fp.Name) > 0 {
		lowerName := strings.ToLower(fp.Name)
		whereParams += "and LOWER(a.nama) LIKE CONCAT('%', ?, '%') "
		args = append(args, lowerName)
	}

	if fp.IdBrandMotor > 0 {
		whereParams += "and a.id_merk_motor = ? "
		args = append(args, fp.IdBrandMotor)
	}

	if fp.IdCategoryMotor > 0 {
		whereParams += "and a.id_kategori_motor = ? "
		args = append(args, fp.IdCategoryMotor)
	}

	queryRecords := `select count(a.id)
	from tblmotor a 
	inner join tblmerkmotor b on a.id_merk_motor = b.id
	inner join tblkategorimotor c on a.id_kategori_motor = c.id
	where 1=1 ` + whereParams
	err = q.db.QueryRowContext(ctx, queryRecords, args...).Scan(&totalData)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	query := `select a.id, a.nama, b.nama as brand, c.nama as kategori, a.icon, a.id_merk_motor, a.id_kategori_motor
	from tblmotor a 
	inner join tblmerkmotor b on a.id_merk_motor = b.id
	inner join tblkategorimotor c on a.id_kategori_motor = c.id
	where 1=1 ` + whereParams + ` 
	order by a.nama asc
	limit ? offset ?`
	args = append(args, fp.Limit, offsetNum)

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i MotorMD

		if err = rows.Scan(
			&i.Id,
			&i.Name,
			&i.BrandMotor,
			&i.CategoryMotor,
			&i.Icon,
			&i.IdBrandMotor,
			&i.IdCategoryMotor,
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

func (q *SqlRepository) GetListCategoryMotor(ctx context.Context) (res []CategoryMotor, errCode string, err error) {
	const query = `select id, nama, icon from tblkategorimotor`

	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i CategoryMotor

		if err = rows.Scan(
			&i.Id,
			&i.Name,
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

func (q *SqlRepository) MotorAdd(ctx context.Context, name, idBrandMotor, idCategoryMotor, icon string) (errCode string, err error) {
	const queryInsert = `insert into tblmotor (nama, id_merk_motor, id_kategori_motor, icon) VALUES (?, ?, ?, ?) `

	_, err = q.db.ExecContext(ctx, queryInsert, name, idBrandMotor, idCategoryMotor, icon)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}

	return
}

func (q *SqlRepository) MotorUpdate(ctx context.Context, idMotor int, name string, idBrandMotor, idCategoryMotor int) (errCode string, err error) {
	const query = `update tblmotor set nama = ?, id_merk_motor = ?, id_kategori_motor = ? where id = ? `
	_, err = q.db.ExecContext(ctx, query, name, idBrandMotor, idCategoryMotor, idMotor)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) MotorCheckExists(ctx context.Context, idMotor string) (exists bool, errCode string, err error) {
	const query = `SELECT EXISTS(SELECT * FROM tblmotor WHERE id = ?)`
	row := q.db.DB.QueryRowContext(ctx, query, idMotor)
	err = row.Scan(&exists)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) MotorUpdateImage(ctx context.Context, idMotor, fileName, uploadPath, dirFile string) (errCode string, err error) {
	_ = q.removeMotorImage(ctx, idMotor, uploadPath, dirFile)

	const query = `update tblmotor set icon = ? where id = ? `
	_, err = q.db.ExecContext(ctx, query, fileName, idMotor)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) removeMotorImage(ctx context.Context, id string, uploadPath, dirFile string) error {
	var temp = Motor{}
	const querySelect = `SELECT icon from tblmotor where id = ?`
	row := q.db.DB.QueryRowContext(ctx, querySelect, id)

	err := row.Scan(
		&temp.Icon,
	)
	if err != nil {
		return err
	}
	if temp.Icon == constants.DefaultImgPng {
		return nil
	}
	helper.RemoveFile(temp.Icon, uploadPath, dirFile)
	return nil
}

func (q *SqlRepository) MotorCheckUsed(ctx context.Context, idMotor string) (exists bool, errCode string, err error) {
	const query = `select exists(select id from motor_x_size_ban where id_motor = ? limit 1)`
	row := q.db.DB.QueryRowContext(ctx, query, idMotor)
	err = row.Scan(&exists)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) MotorRemove(ctx context.Context, idMotor, uploadPath, dirFile string) (errCode string, err error) {
	err = q.removeMotorImage(ctx, idMotor, uploadPath, dirFile)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	const query = `delete from tblmotor where id = ? `

	_, err = q.db.ExecContext(ctx, query, idMotor)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) GetListUkuranBanRaw(ctx context.Context) (res []UkuranRingBan, errCode string, err error) {
	const query = `select IDUkuranBan from tblmasterukuranban order by IDUkuranBan asc`
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	defer rows.Close()

	for rows.Next() {

		var i UkuranRingBan

		if err = rows.Scan(
			&i.IdUkuranBan,
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

func (q *SqlRepository) TireSizeExist(ctx context.Context, id string) (exists bool, errCode string, err error) {
	const query = `SELECT EXISTS(SELECT * FROM tblbanukuranring WHERE id = ?)`
	row := q.db.DB.QueryRowContext(ctx, query, id)
	err = row.Scan(&exists)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}
func (q *SqlRepository) TireSizeAdd(ctx context.Context, id, idRing, idSize string) (errCode string, err error) {
	const queryInsert = `insert into tblbanukuranring (id, id_ring_ban, id_ukuran_ban) VALUES (?, ?, ?) `

	_, err = q.db.ExecContext(ctx, queryInsert, id, idRing, idSize)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) TireSizeUsed(ctx context.Context, id string) (exists bool, errCode string, err error) {
	var (
		matriksExists = false
		productExists = false
	)
	exists = false
	const query = `select exists(select id from motor_x_size_ban where id_ukuran_ring = ? limit 1)`
	row := q.db.DB.QueryRowContext(ctx, query, id)
	err = row.Scan(&matriksExists)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	const querySec = `select exists(select KodePlu from tblmasterplu where IDUkuranRing = ? limit 1)`
	rowSec := q.db.DB.QueryRowContext(ctx, querySec, id)
	err = rowSec.Scan(&productExists)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	if productExists || matriksExists {
		exists = true
	}
	return
}

func (q *SqlRepository) TireSizeDelete(ctx context.Context, id string) (errCode string, err error) {
	const query = `delete from tblbanukuranring where id = ? `

	_, err = q.db.ExecContext(ctx, query, id)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) TireRingExist(ctx context.Context, id int) (exists bool, errCode string, err error) {
	const query = `SELECT EXISTS(SELECT * FROM tblmasterringban WHERE IDRing = ?)`
	row := q.db.DB.QueryRowContext(ctx, query, id)
	err = row.Scan(&exists)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) TireRingAdd(ctx context.Context, id int, nameRing string) (errCode string, err error) {
	const queryInsert = `insert into tblmasterringban (IDRing, UkuranRing, ranking) VALUES (?, ?, ?) `

	_, err = q.db.ExecContext(ctx, queryInsert, id, nameRing, 99)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}
