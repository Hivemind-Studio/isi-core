package user

import (
	"database/sql"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (r *Repository) GetTest(ctx *fiber.Ctx, id int) (result string, err error) {
	return "Test Get Repo" + strconv.FormatInt(int64(id), 10), nil
}

//func (r *Repository) Create(ctx *fiber.Ctx, body *user.RegisterDTO) (result *user.RegisterDTO, err error) {
//	// Start the transaction
//	err = r.StartTx()
//	if err != nil {
//		return result, fmt.Errorf("failed to start transaction: %w", err)
//	}
//
//	// Get the transaction context
//	tx, err := r.GetTx()
//	if err != nil {
//		// Rollback if unable to get transaction
//		r.RollbackTx()
//		return result, fmt.Errorf("failed to get transaction: %w", err)
//	}
//
//	// Check if the email already exists in the database
//	var existingID int
//	checkEmailQuery := `SELECT id FROM users WHERE email = ?`
//	err = tx.QueryRow(checkEmailQuery, body.Email).Scan(&existingID)
//	if err == nil {
//		// Email already exists, rollback and return error
//		r.RollbackTx()
//		return result, fmt.Errorf("email already exists")
//	} else if err != sql.ErrNoRows {
//		// Some other error occurred, rollback and return error
//		r.RollbackTx()
//		return result, fmt.Errorf("failed to check for duplicate email: %w", err)
//	}
//
//	// Insert the user
//	insertQuery := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
//	_, err = tx.Exec(insertQuery, body.Name, body.Email, body.Password)
//	if err != nil {
//		// Rollback in case of an error
//		r.RollbackTx()
//		return result, fmt.Errorf("failed to insert user: %w", err)
//	}
//
//	// Commit the transaction if everything is successful
//	err = r.CommitTx()
//	if err != nil {
//		// Rollback if commit fails
//		r.RollbackTx()
//		return result, fmt.Errorf("failed to commit transaction: %w", err)
//	}
//
//	result = body
//
//	return result, nil
//}

func (r *Repository) Create(ctx *fiber.Ctx, body *user.RegisterDTO) (result *user.RegisterDTO, err error) {
	// Start the transaction
	err = r.StartTx()
	if err != nil {
		return result, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Get the transaction context
	tx, err := r.GetTx()
	if err != nil {
		// Rollback if unable to get transaction
		_ = r.RollbackTx()
		return result, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Check if the email already exists in the database
	var existingID int
	checkEmailQuery := `SELECT id FROM users WHERE email = ?`
	err = tx.QueryRow(checkEmailQuery, body.Email).Scan(&existingID)
	if err == nil {
		// Email already exists, rollback and return error
		_ = r.RollbackTx()
		return result, fmt.Errorf("email already exists")
	} else if err != sql.ErrNoRows {
		// Some other error occurred, rollback and return error
		_ = r.RollbackTx()
		return result, fmt.Errorf("failed to check for duplicate email: %w", err)
	}

	// Insert the user
	insertQuery := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
	_, err = tx.Exec(insertQuery, body.Name, body.Email, body.Password)
	if err != nil {
		// Rollback in case of an error
		_ = r.RollbackTx()
		return result, fmt.Errorf("failed to insert user: %w", err)
	}

	// Commit the transaction if everything is successful
	err = r.CommitTx()
	if err != nil {
		// Rollback if commit fails
		_ = r.RollbackTx()
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return the result
	result = body
	return result, nil
}
