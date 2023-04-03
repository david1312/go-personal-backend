package repo_reports

import (
	"context"

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

func (q *SqlRepository) SyncUpSales(ctx context.Context, fileName string) (err error) {
	return nil
}
