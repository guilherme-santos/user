package mysql

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/guilherme-santos/user"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s UserStorage) Create(ctx context.Context, u *user.User) error {
	id := xid.New().String()
	passwd, err := hashPassword(u.Password)
	if err != nil {
		return &user.FieldError{
			Err: user.Error{
				Type:    user.InvalidArgument,
				Code:    "invalid_password",
				Message: err.Error(),
			},
			Field: "password",
		}
	}

	query := `
		INSERT INTO user
			(id, first_name, last_name, nickname, password, email, country)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = s.db.ExecContext(ctx, query,
		id,
		u.FirstName,
		u.LastName,
		u.Nickname,
		passwd,
		u.Email,
		u.Country,
	)
	if err != nil {
		if IsDuplicateError(err, "user.email") {
			return &user.FieldError{
				Err: user.Error{
					Type:    user.InvalidArgument,
					Code:    "email_already_exists",
					Message: "Provided e-mail already exists",
				},
				Field: "email",
			}
		}
		return err
	}
	u.ID = id
	return nil
}

func (s UserStorage) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE user
		SET
			first_name = ?,
			last_name = ?,
			nickname = ?,
			email = ?,
			country = ?
	`
	args := []interface{}{
		u.FirstName,
		u.LastName,
		u.Nickname,
		u.Email,
		u.Country,
	}
	if u.Password != "" {
		query += ", password = ?"
		passwd, err := hashPassword(u.Password)
		if err != nil {
			return &user.FieldError{
				Err: user.Error{
					Type:    user.InvalidArgument,
					Code:    "invalid_password",
					Message: err.Error(),
				},
				Field: "password",
			}
		}
		args = append(args, passwd)
	}
	query += " WHERE id = ?"
	args = append(args, u.ID)

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		if IsDuplicateError(err, "user.email") {
			return &user.FieldError{
				Err: user.Error{
					Type:    user.InvalidArgument,
					Code:    "email_already_exists",
					Message: "Provided e-mail already exists",
				},
				Field: "email",
			}
		}
	}
	return err
}

func (s UserStorage) Delete(ctx context.Context, id string) error {
	query := `UPDATE user SET removed_at = NOW() WHERE id = ? AND removed_at IS NULL`
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return user.ErrNotFound
	}
	return nil
}

func (s UserStorage) Get(ctx context.Context, id string) (*user.User, error) {
	query := `
		SELECT id, first_name, last_name, nickname, email, country, created_at, updated_at, removed_at
		FROM user
		WHERE id = ?
	`

	row := s.db.QueryRowContext(ctx, query, id)
	u, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s UserStorage) List(ctx context.Context, opts *user.ListOptions) (*user.ListResponse, error) {
	query := `
		SELECT id, first_name, last_name, nickname, email, country, created_at, updated_at, removed_at
		FROM user
	`
	where := []string{"removed_at IS NULL"}
	args := []interface{}{}

	if opts.Country != "" {
		where = append(where, "country = ?")
		args = append(args, opts.Country)
	}
	if opts.Cursor != "" {
		// TODO: implement cursor based
	}
	query += " WHERE " + strings.Join(where, " AND ")
	if opts.Sort != "" {
		var mode, field string
		if strings.HasPrefix(opts.Sort, "-") {
			mode = "DESC"
			field = strings.TrimPrefix(opts.Sort, "-")
		} else {
			mode = "ASC"
			field = opts.Sort
		}
		query += " ORDER BY " + field + " " + mode
	}
	query += " LIMIT " + strconv.FormatInt(opts.PerPage, 10)
	if opts.Cursor == "" {
		query += " OFFSET " + strconv.FormatInt(int64(opts.Page)*opts.PerPage, 10)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*user.User, 0)

	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	// TODO: get total of users using COUNT(*) from the SELECT

	return &user.ListResponse{
		Total:   int64(len(users)),
		PerPage: opts.PerPage,
		Users:   users,
	}, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanUser(row rowScanner) (*user.User, error) {
	u := new(user.User)
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Nickname,
		&u.Email,
		&u.Country,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.RemovedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// hashPassword returns a hashed version from the password
func hashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(bytes), err
}
