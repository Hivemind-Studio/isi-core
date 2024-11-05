package user

import (
	"database/sql"
	"errors"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"time"
)

func (r *Repository) Create(ctx *fiber.Ctx, tx *sqlx.Tx, name string, email string, password string, roleId int64,
) (err error) {
	var existingID int
	checkEmailQuery := `SELECT id FROM users WHERE email = ?`
	err = tx.QueryRow(checkEmailQuery, email).Scan(&existingID)

	if err == nil {
		// Email already exists
		return httperror.New(fiber.StatusBadRequest, "email already exists")
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Unexpected error during the duplicate check
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to check duplicate")
	}

	hashedPassword, hashErr := hash.HashPassword(password)
	if hashErr != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, hashErr, "failed to hash password")
	}

	insertQuery := `INSERT INTO users (name, email, password, role_id, status, verification) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(insertQuery, name, email, hashedPassword, roleId, 1, 0)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to insert")
	}

	return nil
}

func (r *Repository) FindByEmail(ctx *fiber.Ctx, email string) (Cookie, error) {
	var result Cookie

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

	err := r.GetConnDb().QueryRowx(query, email).StructScan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Cookie{}, httperror.New(fiber.StatusNotFound, "user not found")
		}
		return Cookie{}, httperror.New(fiber.StatusInternalServerError, err.Error())
	}

	return result, nil
}

func (r *Repository) GetUsers(ctx *fiber.Ctx, name string, email string,
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

	err := r.GetConnDb().Select(&users, query, args...)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve users")
	}

	return users, nil
}
