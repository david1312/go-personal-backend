package repo_customers

import (
	"context"
	"database/sql"
	"errors"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/helper"

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

func (q *SqlRepository) Login(ctx context.Context, email, password string) (res Customers, errCode string, err error) {
	const query = `SELECT name, password, uid FROM customers where email = ? AND deleted_at IS NULL`
	row := q.db.DB.QueryRowContext(ctx, query, email)

	err = row.Scan(
		&res.Name,
		&res.Password,
		&res.Uid,
	)
	if err != nil && err == sql.ErrNoRows {
		errCode = crashy.ErrInvalidUser
		return
	}
	if err != nil && err != sql.ErrNoRows {
		errCode = crashy.ErrCodeUnexpected
		return
	}

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

func (q *SqlRepository) Register(ctx context.Context, name, email, emailToken, password, uid string) (errCode string, err error) {
	const query = `insert into customers (uid, name, password, email, email_verified_token, is_active, email_verified_sent) 
	VALUES (?, ?, ?, ?, ?, true, 1) `

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	_, err = q.db.ExecContext(ctx, query, uid,
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
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		err = errors.New(crashy.ErrInvalidTokenEmail)
		errCode = crashy.ErrInvalidTokenEmail
		return
	}

	return
}

func (q *SqlRepository) GetCustomer(ctx context.Context, uid string) (res Customers, errCode string, err error) {
	const query = `SELECT name, email, email_verified_at, gender, phone, phone_verified_at, avatar, birthdate FROM customers where uid = ? AND deleted_at IS NULL`
	row := q.db.DB.QueryRowContext(ctx, query, uid)

	err = row.Scan(
		&res.Name,
		&res.Email,
		&res.EmailVerifiedAt,
		&res.Gender,
		&res.Phone,
		&res.PhoneVerifiedAt,
		&res.Avatar,
		&res.Birthdate,
	)

	if err != nil {
		errCode = crashy.ErrCodeDataRead
		return
	}
	return
}

func (q *SqlRepository) ChangePassword(ctx context.Context, uid, oldPass, newPass string) (errCode string, err error) {
	const query = `SELECT password FROM customers where uid = ? AND deleted_at IS NULL`
	row := q.db.DB.QueryRowContext(ctx, query, uid)

	var res Customers

	err = row.Scan(
		&res.Password,
	)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	// Comparing the password with the hash
	if err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(oldPass)); err != nil {
		errCode = crashy.ErrInvalidOldPassword
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	//update the password
	const queryUpdate = `update customers set password =  ? where uid = ?`
	_, err = q.db.ExecContext(ctx, queryUpdate, hashedPass, uid)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	return
}

func (q *SqlRepository) ResendEmail(ctx context.Context, uid, email string) (emailToken, errCode string, err error) {
	const query = `SELECT email_verified_token FROM customers where email = ? and uid = ? AND deleted_at IS NULL`
	row := q.db.DB.QueryRowContext(ctx, query, email, uid)

	err = row.Scan(&emailToken)

	if err != nil && err == sql.ErrNoRows {
		errCode = crashy.ErrInvalidEmail
		return
	}
	if err != nil && err != sql.ErrNoRows {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	return
}

func (q *SqlRepository) RequestPinEmail(ctx context.Context, uid, email string) (pin, errCode string, err error) {
	var exist bool
	const query = `SELECT EXISTS(SELECT email FROM customers where email = ? and uid = ? AND deleted_at IS NULL)`
	row := q.db.DB.QueryRowContext(ctx, query, email, uid)
	err = row.Scan(&exist)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	if !exist {
		err = errors.New(crashy.ErrInvalidEmail)
		errCode = crashy.ErrInvalidEmail
		return
	}

	pin = helper.RandomNumber(6)

	//update the password
	const queryUpdate = `update customers set email_change_code = ?, email_change_eligible = true where uid = ?`
	_, err = q.db.ExecContext(ctx, queryUpdate, pin, uid)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	return
}

func (q *SqlRepository) ChangeEmail(ctx context.Context, uid, oldEmail, newEmail, hashedTokenEmail, code string) (errCode string, err error) { //simplify
	var codeDB string
	const query = `SELECT email_change_code FROM customers where email = ? and uid = ? AND email_change_eligible = true AND deleted_at IS NULL`
	row := q.db.DB.QueryRowContext(ctx, query, oldEmail, uid)
	err = row.Scan(&codeDB)

	if err != nil && err == sql.ErrNoRows {
		errCode = crashy.ErrInvalidEmail
		return
	}
	if err != nil && err != sql.ErrNoRows {
		errCode = crashy.ErrCodeUnexpected
		return
	}

	if code != codeDB {
		err = errors.New(crashy.ErrInvalidCode)
		errCode = crashy.ErrInvalidCode
		return
	}

	//update the password
	const queryUpdate = `update customers set email =  ?, email_verified_at = NULL, email_verified_token = ?, email_change_code = NULL, email_change_eligible = false where uid = ?`
	_, err = q.db.ExecContext(ctx, queryUpdate, newEmail, hashedTokenEmail, uid)

	if err != nil {
		errCode = crashy.ErrCodeUnexpected
		return
	}
	return
}
