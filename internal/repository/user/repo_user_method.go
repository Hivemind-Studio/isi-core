package user

import (
	"database/sql"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (r *UserRepository) GetTest(ctx *fiber.Ctx, id int) (result string, err error) {
	return "Test Get Repo" + strconv.FormatInt(int64(id), 10), nil
}

func (r *UserRepository) Create(ctx *fiber.Ctx, body *user.RegisterDTO) (result *user.RegisterResponse, err error) {
	// Start the transaction
	err = r.StartTx()
	if err != nil {
		return result, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Get the transaction context
	tx, err := r.GetTx()
	if err != nil {
		_ = r.RollbackTx()
		return result, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Check if the email already exists in the database
	var existingID int
	checkEmailQuery := `SELECT id FROM users WHERE email = ?`
	err = tx.QueryRow(checkEmailQuery, body.Email).Scan(&existingID)
	if err == nil {
		_ = r.RollbackTx()
		return result, fmt.Errorf("email already exists")
	} else if err != sql.ErrNoRows {
		_ = r.RollbackTx()
		return result, fmt.Errorf("failed to check for duplicate email: %w", err)
	}

	// Hash the password with the salt
	hashedPassword, _ := utils.HashPassword(body.Password)

	// Insert the user into the database with the hashed password and salt
	insertQuery := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
	_, err = tx.Exec(insertQuery, body.Name, body.Email, hashedPassword)
	if err != nil {
		_ = r.RollbackTx()
		return result, fmt.Errorf("failed to insert user: %w", err)
	}

	// Commit the transaction
	err = r.CommitTx()
	if err != nil {
		_ = r.RollbackTx()
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Don't return the password in the response
	result = &user.RegisterResponse{
		Name:  body.Name,
		Email: body.Email,
	}

	return result, nil
}

func (r *UserRepository) Login(ctx *fiber.Ctx, body *user.LoginDTO) (result string, err error) {
	// Start a new transaction
	err = r.StartTx()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %w", err)
	}

	// Retrieve the transaction context
	tx, err := r.GetTx()
	if err != nil {
		_ = r.RollbackTx() // Rollback in case of error
		return "", fmt.Errorf("failed to get transaction: %w", err)
	}

	// Retrieve the user record by email, including the salt and hashed password
	var storedUser user.LoginDTO

	query := `SELECT email, password FROM users WHERE email = ?`
	err = tx.QueryRowx(query, body.Email).StructScan(&storedUser)
	if err != nil {
		if err == sql.ErrNoRows {
			_ = r.RollbackTx() // Rollback in case of error
			return "", fmt.Errorf("user not found")
		}
		_ = r.RollbackTx() // Rollback in case of other errors
		return "", fmt.Errorf("failed to retrieve user: %w", err)
	}

	// Compare the hashed password using the stored salt
	isValidPassword, _ := utils.ComparePassword(storedUser.Password, body.Password)
	if !isValidPassword {
		_ = r.RollbackTx() // Rollback if password doesn't match
		return "", fmt.Errorf("invalid password")
	}

	// Commit the transaction if everything is successful
	err = r.CommitTx()
	if err != nil {
		_ = r.RollbackTx() // Rollback if commit fails
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return success message
	return "Success", nil
}
