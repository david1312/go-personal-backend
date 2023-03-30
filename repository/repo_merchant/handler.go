package repo_merchant

import (
	"context"
	"database/sql"
	"libra-internal/pkg/crashy"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type SqlRepository struct {
	db *sqlx.DB
}

func NewSqlRepository(db *sqlx.DB) *SqlRepository {
	return &SqlRepository{
		db: db,
	}
}

func (q *SqlRepository) Login(ctx context.Context, username, password string) (res MerchantData, errCode string, err error) {
	const query = `select username, outlet_id, password from merchants where username = ? and deleted_at IS NULL`
	row := q.db.DB.QueryRowContext(ctx, query, username)

	err = row.Scan(
		&res.Username,
		&res.OutletId,
		&res.Password,
	)
	if err != nil && err == sql.ErrNoRows {
		errCode = crashy.ErrInvalidUserMerchant
		return
	}
	if err != nil && err != sql.ErrNoRows {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	// Comparing the password with the hash
	if err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password)); err != nil {
		errCode = crashy.ErrInvalidUserMerchant
		return
	}
	return
}

func (q *SqlRepository) GetMerchantProfile(ctx context.Context, username string) (res MerchantData, errCode string, err error) {
	const query = `select a.username, a.email, a.phone, a.avatar, a.outlet_id, b.name, b.address, b.city, b.gmap_url
	from merchants a
	join outlets b on a.outlet_id = b.id
	where a.username = ?`
	row := q.db.DB.QueryRowContext(ctx, query, username)

	err = row.Scan(
		&res.Username,
		&res.OutletEmail,
		&res.CsNumber,
		&res.OutletAvatar,
		&res.OutletId,
		&res.OutletName,
		&res.OutletAddress,
		&res.OutletCity,
		&res.OutletGmapUrl,
	)
	if err != nil && err == sql.ErrNoRows {
		errCode = crashy.ErrInvalidUserMerchant
		return
	}
	if err != nil && err != sql.ErrNoRows {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}
