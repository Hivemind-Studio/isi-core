package user

import (
	"database/sql"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
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
		return result, fmt.Errorf("email already exists")
	} else if err != sql.ErrNoRows {
		return result, fmt.Errorf("failed to check for duplicate email: %w", err)
	}

	hashedPassword, _ := utils.HashPassword(body.Password)

	insertQuery := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
	_, err = tx.Exec(insertQuery, body.Name, body.Email, hashedPassword)
	if err != nil {
		return result, fmt.Errorf("failed to insert user: %w", err)
	}

	result = &user.RegisterResponse{
		Name:  body.Name,
		Email: body.Email,
	}

	return result, nil
}

func (r *Repository) FindByEmail(ctx *fiber.Ctx, tx *sqlx.Tx, body *user.LoginDTO) (User, error) {
	var result User

	// Query only the fields that are needed
	query := `SELECT users.* FROM users WHERE email = ?`
	err := tx.QueryRowx(query, body.Email).StructScan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return result, nil
}
