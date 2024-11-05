package user

import (
	"database/sql"
	"errors"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) Create(ctx *fiber.Ctx, tx *sqlx.Tx, body *user.RegistrationDTO, roleId int) (result *user.RegisterResponse, err error) {
	var existingID int
	checkEmailQuery := `SELECT id FROM users WHERE email = ?`
	err = tx.QueryRow(checkEmailQuery, body.Email).Scan(&existingID)

	if err == nil {
		return nil, httperror.New(fiber.StatusBadRequest, "email already exists")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to check duplicate")
	}

	hashedPassword, hashErr := utils.HashPassword(body.Password)
	if hashErr != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, hashErr, "failed to hash password")
	}

	insertQuery := `INSERT INTO users (name, email, password, role_id) VALUES (?, ?, ?, ?)`
	_, err = tx.Exec(insertQuery, body.Name, body.Email, hashedPassword, roleId)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to insert")
	}

	result = &user.RegisterResponse{
		Name:  body.Name,
		Email: body.Email,
	}

	return result, nil
}

func (r *Repository) FindByEmail(ctx *fiber.Ctx, tx *sqlx.Tx, email string,
) (User, error) {
	var result User

	query := `SELECT users.* FROM users WHERE email = ?`
	err := tx.QueryRowx(query, email).StructScan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, httperror.New(fiber.StatusNotFound, "user not found")
		}
		return User{}, httperror.New(fiber.StatusInternalServerError, err.Error())
	}

	return result, nil
}

func (r *Repository) GetUsers(ctx *fiber.Ctx, tx *sqlx.Tx, name string, email string,
) ([]User, error) {
	var users []User
	var query string
	var args []interface{}

	query = "SELECT * FROM users WHERE 1=1"

	if name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if email != "" {
		query += " AND email LIKE ?"
		args = append(args, "%"+email+"%")
	}

	err := tx.Select(&users, query, args...)
	if err != nil {
		return nil, httperror.
			Wrap(fiber.StatusInternalServerError, err, "failed to retrieve users")
	}

	return users, nil
}
