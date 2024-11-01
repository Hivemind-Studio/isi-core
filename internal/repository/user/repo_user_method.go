package user

import (
	"database/sql"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"strconv"
)

func (r *Repository) GetTest(ctx *fiber.Ctx, id int) (result string, err error) {
	return "Test Get Repo" + strconv.FormatInt(int64(id), 10), nil
}

func (r *Repository) Create(ctx *fiber.Ctx, tx *sqlx.Tx, body *user.RegisterDTO,
) (result *user.RegisterResponse, err error) {
	var existingID int
	checkEmailQuery := `SELECT id FROM users WHERE email = ?`
	err = tx.QueryRow(checkEmailQuery, body.Email).Scan(&existingID)
	if err == nil {
		return result, httperror.Wrap(fiber.StatusBadRequest, err, "email already exists")
	} else if err != sql.ErrNoRows {
		return result, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to check duplicate")
	}

	hashedPassword, _ := utils.HashPassword(body.Password)

	insertQuery := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
	_, err = tx.Exec(insertQuery, body.Name, body.Email, hashedPassword)
	if err != nil {
		return result, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to insert")
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
		if err == sql.ErrNoRows {
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
