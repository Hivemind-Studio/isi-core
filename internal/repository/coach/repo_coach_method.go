package coach

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"time"
)

func (r *Repository) GetCoaches(ctx context.Context, params coach.QueryCoachDTO, page int64, perPage int64) ([]Coach, pagination.Pagination, error) {
	var users []Coach
	var totalRecords int64
	var args []interface{}

	baseQuery := `
        FROM
            users u
            LEFT JOIN coaches c ON u.id = c.user_id
            LEFT JOIN roles r ON u.role_id = r.id
        WHERE
            u.role_id = 3
    `

	if params.Name != "" {
		baseQuery += " AND u.name LIKE ?"
		args = append(args, "%"+params.Name+"%")
	}
	if params.Email != "" {
		baseQuery += " AND u.email LIKE ?"
		args = append(args, "%"+params.Email+"%")
	}
	if params.PhoneNumber != "" {
		baseQuery += " AND u.phone_number LIKE ?"
		args = append(args, "%"+params.PhoneNumber+"%")
	}
	if params.Level != "" {
		baseQuery += " AND c.level LIKE ?"
		args = append(args, "%"+params.Level+"%")
	}
	if params.Status != "" {
		baseQuery += " AND u.status = ?"
		args = append(args, params.Status == "true")
	}
	if params.StartDate != nil {
		baseQuery += " AND u.created_at >= ?"
		args = append(args, *params.StartDate)
	}
	if params.EndDate != nil {
		baseQuery += " AND u.created_at <= ?"
		args = append(args, *params.EndDate)
	}

	countQuery := "SELECT COUNT(DISTINCT u.id) " + baseQuery
	err := r.GetConnDb().GetContext(ctx, &totalRecords, countQuery, args...)
	if err != nil {
		return nil, pagination.Pagination{}, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to count coaches")
	}

	selectQuery := `
        SELECT
            u.id AS id,
            u.name AS name,
            u.email AS email,
            u.address AS address,
            u.phone_number AS phone_number,
            u.date_of_birth AS date_of_birth,
            u.gender AS gender,
            u.verification AS verification,
            u.occupation AS occupation,
            u.photo AS photo,
            u.status AS status,
            u.created_at AS created_at,
            u.updated_at AS updated_at,
            c.certifications AS certifications,
            c.experiences AS experiences,
            c.education AS educations,
            c.level AS level,
            r.id AS role_id,
            r.name AS role_name
        ` + baseQuery + ` LIMIT ? OFFSET ?`

	queryArgs := append(args, perPage, (page-1)*perPage)
	err = r.GetConnDb().SelectContext(ctx, &users, selectQuery, queryArgs...)
	if err != nil {
		return nil, pagination.Pagination{}, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve coaches")
	}

	// Calculate total pages
	totalPages := (totalRecords + perPage - 1) / perPage

	paginate := pagination.Pagination{
		CurrentPage:  page,
		PerPage:      perPage,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
	}

	return users, paginate, nil
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

func (r *Repository) GetCoachById(ctx context.Context, id int64) (Coach, error) {
	var result Coach

	query := `
		SELECT
			u.id AS id,
			u.name AS name,
			u.email AS email,
			u.address AS address,
			u.phone_number AS phone_number,
			u.date_of_birth AS date_of_birth,
			u.gender AS gender,
			u.verification AS verification,
			u.occupation AS occupation,
			u.photo AS photo,
			u.status AS status,
			u.created_at AS created_at,
			u.updated_at AS updated_at,
			c.certifications AS certifications,
			c.experiences AS experiences,
			c.education AS educations,
			c.level AS level,
			r.id AS role_id,
			r.name AS role_name
		FROM
			users u
				JOIN
			coaches c ON u.id = c.user_id
				JOIN
			roles r ON u.role_id = r.id
		WHERE
			u.role_id = 3 AND u.id = ?
		`

	err := r.GetConnDb().QueryRowxContext(ctx, query, id).StructScan(&result)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Coach{}, httperror.New(fiber.StatusNotFound, "user not found")
		}
		return Coach{}, httperror.New(fiber.StatusInternalServerError, err.Error())
	}

	return result, nil
}
