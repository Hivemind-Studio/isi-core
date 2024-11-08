package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

func (r *Repository) Create(ctx context.Context, tx *sqlx.Tx, name string, email string,
	password string, roleId int64, phoneNumber string) (err error) {
	if err := r.checkExistingData(ctx, tx, email, phoneNumber); err != nil {
		return err
	}

	hashedPassword, hashErr := hash.HashPassword(password)
	if hashErr != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, hashErr, "failed to hash password")
	}

	insertQuery := `INSERT INTO users (name, email, password, role_id, phone_number, status, verification) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, insertQuery, name, email, hashedPassword, roleId, phoneNumber, 1, 0)

	if err != nil {
		return httperror.New(fiber.StatusConflict, "failed to insert")
	}

	return nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (User, error) {
	var result User

	query := `
		SELECT 
			users.id, 
			users.name, 
			users.email, 
			users.password, 
			users.role_id, 
			roles.name AS role_name 
		FROM 
			users 
		LEFT JOIN 
			roles 
		ON 
			users.role_id = roles.id 
		WHERE 
			users.email = ?
	`

	err := r.GetConnDb().QueryRowxContext(ctx, query, email).StructScan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, httperror.New(fiber.StatusNotFound, "user not found")
		}
		return User{}, httperror.New(fiber.StatusInternalServerError, err.Error())
	}

	return result, nil
}

func (r *Repository) GetUsers(ctx context.Context, name string, email string,
	startDate, endDate *time.Time, page int64, perPage int64,
) ([]User, error) {
	var users []User
	query := "SELECT * FROM users WHERE 1=1"
	var args []interface{}

	if name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}
	if email != "" {
		query += " AND email LIKE ?"
		args = append(args, "%"+email+"%")
	}
	if startDate != nil {
		query += " AND created_at >= ?"
		args = append(args, *startDate)
	}
	if endDate != nil {
		query += " AND created_at <= ?"
		args = append(args, *endDate)
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, perPage, (page-1)*perPage)

	err := r.GetConnDb().SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve users")
	}

	return users, nil
}

func (r *Repository) checkForDuplicate(ctx context.Context, tx *sqlx.Tx, column, value string) error {
	var exists string
	query := `SELECT 1 FROM users WHERE ` + column + ` = ?`
	err := tx.QueryRowContext(ctx, query, value).Scan(&exists)

	if err == nil {
		return httperror.New(fiber.StatusBadRequest, column+" already exists")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to check duplicate")
	}
	return nil
}

func (r *Repository) checkExistingData(ctx context.Context, tx *sqlx.Tx, email string, phoneNumber string) error {
	if err := r.checkForDuplicate(ctx, tx, "email", email); err != nil {
		return err
	}

	if err := r.checkForDuplicate(ctx, tx, "phone_number", phoneNumber); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int64) (User, error) {
	var result User

	query := "SELECT * FROM users WHERE id = ? LIMIT 1"

	err := r.GetConnDb().QueryRowxContext(ctx, query, id).StructScan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, httperror.New(fiber.StatusNotFound, "user not found")
		}
		return User{}, httperror.New(fiber.StatusInternalServerError, err.Error())
	}

	return result, nil
}

func (r *Repository) SuspendUsers(ctx context.Context, tx *sqlx.Tx, ids []int64) error {
	if len(ids) == 0 {
		return httperror.New(fiber.StatusBadRequest, "no user IDs provided")
	}

	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf("UPDATE users SET status = 1 WHERE id IN (%s)", strings.Join(placeholders, ","))

	_, err := tx.ExecContext(ctx, query, utils.ToInterfaceSlice(ids)...)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to suspend users")
	}

	return nil
}
