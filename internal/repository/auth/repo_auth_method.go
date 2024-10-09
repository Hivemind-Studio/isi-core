package user

import (
	"database/sql"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
)

func (r *Repository) Login(ctx *fiber.Ctx, body *user.LoginDTO) (result string, err error) {
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
