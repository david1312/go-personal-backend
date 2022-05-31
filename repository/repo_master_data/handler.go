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
	const query = `SELECT IdMerk, Merk, Icon from tblmerkban`
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
