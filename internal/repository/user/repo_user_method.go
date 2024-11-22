package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/mail"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"os"
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

	insertUserQuery := `INSERT INTO users (name, email, password, role_id, phone_number, status, verification) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := tx.ExecContext(ctx, insertUserQuery, name, email, hashedPassword, roleId, phoneNumber, 1, 0)
	if err != nil {
		return httperror.New(fiber.StatusConflict, "failed to insert user")
	}

	userId, err := res.LastInsertId()
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to get last inserted ID")
	}

	token := utils.GenerateVerificationToken()
	expiresAt := time.Now().Add(1 * time.Hour)

	insertVerificationQuery := `INSERT INTO user_verified_account (user_id, verification_token, expires_at) VALUES (?, ?, ?)`
	_, err = tx.ExecContext(ctx, insertVerificationQuery, userId, token, expiresAt)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to insert verification record")
	}

	err, done := r.EmailVerification(name, token, err, email)
	if done {
		return err
	}

	return nil
}

func (r *Repository) EmailVerification(name string, token string, err error, email string) (error, bool) {
	emailClient := mail.NewEmailClient(&mail.EmailConfig{
		Host:        os.Getenv("MAIL_SMTP_HOST"),
		Port:        os.Getenv("MAIL_SMTP_PORT"),
		Username:    os.Getenv("MAIL_SMTP_USERNAME"),
		Password:    os.Getenv("MAIL_SMTP_PASSWORD"),
		SenderEmail: os.Getenv("MAIL_AUTH_EMAIL"),
	})

	emailData := struct {
		Name            string
		VerificationURL string
	}{
		Name:            name,
		VerificationURL: fmt.Sprintf("%stoken=%s", os.Getenv("CALLBACK_VERIFICATION_URL"), token),
	}

	err = emailClient.SendMail(
		[]string{email},
		"Verify Your Email",
		"template/verification_email.html",
		emailData,
	)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send verification email"), true
	}
	return err, false
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

func (r *Repository) GetUsers(ctx context.Context, params dto.GetUsersDTO, page int64, perPage int64,
) ([]User, error) {
	var users []User
	query := "SELECT * FROM users WHERE 1=1"
	var args []interface{}

	if params.Name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+params.Name+"%")
	}
	if params.Email != "" {
		query += " AND email LIKE ?"
		args = append(args, "%"+params.Email+"%")
	}
	if params.Role != nil {
		query += " AND role_id = ?"
		args = append(args, params.Role)
	}
	if params.StartDate != nil {
		query += " AND created_at >= ?"
		args = append(args, *params.StartDate)
	}
	if params.EndDate != nil {
		query += " AND created_at <= ?"
		args = append(args, *params.EndDate)
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
