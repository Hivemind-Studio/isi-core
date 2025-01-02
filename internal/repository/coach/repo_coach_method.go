package coach

import (
	"context"
	"database/sql"
	"errors"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"time"
)

func (r *Repository) GetUsers(ctx context.Context, params dto.GetUsersDTO, page int64, perPage int64,
) ([]user.User, error) {
	var users []user.User
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
	if params.PhoneNumber != "" {
		query += " AND phone_number LIKE ?"
		args = append(args, "%"+params.PhoneNumber+"%")
	}
	if params.Level != "" {
		query += " AND level LIKE ?"
		args = append(args, "%"+params.Level+"%")
	}
	if params.Status != "" {
		query += " AND status LIKE ?"
		args = append(args, "%"+params.Status+"%")
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

func (r *Repository) CreateCoach(ctx context.Context, tx *sqlx.Tx, id int64) (err error) {
	insertCoachQuery := `INSERT INTO coaches (user_id) VALUES (?)`
	_, err = tx.ExecContext(ctx, insertCoachQuery, id)
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "failed to insert coach")
	}

	return nil
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

func (r *Repository) InsertEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string,
	token string, expiredAt time.Time,
) error {
	insertQuery := `
			INSERT INTO email_verifications (email, verification_token, expired_at, trial)
			VALUES (?, ?, ?, 1)
		`
	_, err := tx.ExecContext(ctx, insertQuery, email, token, expiredAt)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err,
			"failed to insert verification record")
	}

	return nil
}

func (r *Repository) UpdateEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string,
	targetDate string, token string, expiredAt time.Time,
) error {
	updateQuery := `
			UPDATE email_verifications
			SET verification_token = ?, expired_at = ?, trial = trial + 1, updated_at = NOW()
			WHERE email = ? AND DATE(created_at) = ?
		`
	_, err := tx.ExecContext(ctx, updateQuery, token, expiredAt, email, targetDate)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err,
			"failed to update verification record")
	}

	return nil
}

func (r *Repository) GetEmailVerificationTrialRequestByDate(ctx context.Context, email string, queryDate time.Time,
) (*int8, error) {
	filterDate := queryDate.Format("2006-01-02")

	query := `SELECT trial FROM email_verifications WHERE email = ? AND DATE(created_at) = ?`
	var trial int8 = 0
	err := r.GetConnDb().QueryRowContext(ctx, query, email, filterDate).Scan(&trial)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to fetch verification record")
		}
	}

	return &trial, nil
}

func (r *Repository) GetTokenEmailVerification(token string) (string, error) {
	query := `
		SELECT email, expired_at 
		FROM email_verifications 
		WHERE verification_token = ?
	`
	var email string
	var expiredAt time.Time

	err := r.GetConnDb().QueryRowxContext(context.Background(), query, token).Scan(&email, &expiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", httperror.New(fiber.StatusNotFound, "verification token not found")
		}
		return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to fetch verification record")
	}

	if time.Now().After(expiredAt) {
		return "", httperror.Wrap(fiber.StatusBadRequest, nil, "verification token is expired")
	}

	return email, nil
}

func (r *Repository) UpdateCoachPassword(ctx context.Context, tx *sqlx.Tx, password string, email string) error {
	hashedPassword, hashErr := hash.HashPassword(password)
	if hashErr != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, hashErr, "failed to hash password")
	}

	query := `UPDATE users SET password = ? WHERE email = ?`
	_, err := tx.ExecContext(ctx, query, hashedPassword, email)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update user password")
	}

	return nil
}
