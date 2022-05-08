package repo_customers

import (
	"context"
	"database/sql"
	"errors"
	"semesta-ban/pkg/crashy"

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

func (q *SqlRepository) Login(ctx context.Context, email, password string) (nama, errCode string, err error) {
	const query = `SELECT name, password FROM customers where email = ? AND deleted_at IS NULL`
	row := q.db.DB.QueryRowContext(ctx, query, email)

	var res Customers
	err = row.Scan(
		&res.Name,
		&res.Password,
	)
	if err != nil && err == sql.ErrNoRows {
		errCode = crashy.ErrInvalidUser
		return
	}
	if err != nil && err != sql.ErrNoRows {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	nama = res.Name

	// Comparing the password with the hash
	if err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password)); err != nil {
		errCode = crashy.ErrInvalidUser
		return
	}

	return
}

func (q *SqlRepository) CheckEmailExist(ctx context.Context, email string) (res bool, errCode string, err error) {
	const query = `select EXISTS(select name from customers where email = ? and deleted_at IS NULL)`
	row := q.db.DB.QueryRowContext(ctx, query, email)
	err = row.Scan(&res)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}

func (q *SqlRepository) Register(ctx context.Context, name, email, emailToken, password string) (errCode string, err error) {
	const query = `insert into customers (name, password, email, email_verified_token, is_active) 
	VALUES (?, ?, ?, ?, true) `

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	_, err = q.db.ExecContext(ctx, query,
		name,
		string(hashedPass), email, emailToken,
	)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}
	return
}

func (q *SqlRepository) VerifyEmail(ctx context.Context, emailToken string) (errCode string, err error) {
	const query = `update customers set email_verified_at = now() where email_verified_token = ?`
	res, err := q.db.ExecContext(ctx, query, emailToken)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		err = errors.New(crashy.ErrInvalidTokenEmail)
		errCode = crashy.ErrInvalidTokenEmail
		return
	}

	return
}
