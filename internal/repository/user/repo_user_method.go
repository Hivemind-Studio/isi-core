package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant/loglevel"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"strings"
	"time"

	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) Create(ctx context.Context, tx *sqlx.Tx, name string, email string, password *string,
	roleId int64, phoneNumber *string, gender string, address string, status int, googleId *string,
) (id int64, err error) {
	if err := r.checkExistingData(ctx, tx, email, phoneNumber, nil); err != nil {
		return 0, err
	}

	var hashedPassword *string
	if password != nil {
		hashed, hashErr := hash.HashPassword(*password)
		if hashErr != nil {
			return 0, httperror.Wrap(fiber.StatusInternalServerError, hashErr, "failed to hash password")
		}
		hashedPassword = &hashed
	}

	phoneValue := interface{}(nil)
	if phoneNumber != nil {
		phoneValue = *phoneNumber
	}

	insertUserQuery := `INSERT INTO users (name, email, password, role_id, phone_number, status, gender,
                   address, verification, version, google_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, insertUserQuery, name, email, hashedPassword, roleId, phoneValue, status,
		gender, address, 0, 0, googleId)

	if err != nil {
		return 0, httperror.New(fiber.StatusConflict, "failed to insert user")
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return 0, httperror.New(fiber.StatusInternalServerError, "failed to retrieve last inserted user ID")
	}

	return userId, nil
}

func (r *Repository) CreateStaff(ctx context.Context, tx *sqlx.Tx, name string, email string,
	password string, address string, phoneNumber *string, status int, gender string, role string) (id int64, err error) {
	if err := r.checkExistingData(ctx, tx, email, phoneNumber, nil); err != nil {
		return 0, err
	}

	hashedPassword, hashErr := hash.HashPassword(password)
	if hashErr != nil {
		return 0, httperror.Wrap(fiber.StatusInternalServerError, hashErr, "failed to hash password")
	}

	insertUserQuery := `INSERT INTO users (name, email, password, role_id, phone_number, status, 
                  address, gender, verification, version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, insertUserQuery, name, email, hashedPassword, role, phoneNumber, status,
		address, gender, 0, 0)

	if err != nil {
		return 0, httperror.Wrap(fiber.StatusConflict, err, "failed to insert user")
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return 0, httperror.New(fiber.StatusInternalServerError, "failed to retrieve last inserted user ID")
	}

	return userId, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (User, error) {
	var result User

	query := `
		SELECT 
			*
		FROM 
			users 
		LEFT JOIN 
			roles 
		ON 
			users.role_id = roles.id 
		WHERE 
			users.email = ?
	`

	err := r.GetConnDb().QueryRowxContext(ctx, query, &email).StructScan(&result)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, httperror.New(fiber.StatusNotFound, "user not found")
		}
		return User{}, httperror.New(fiber.StatusInternalServerError, err.Error())
	}

	return result, nil
}

func (r *Repository) GetUsers(ctx context.Context, params dto.GetUsersDTO, page int64, perPage int64) ([]User, pagination.Pagination, error) {
	var users []User
	var totalRecords int64
	var args []interface{}

	baseQuery := `FROM users AS users LEFT JOIN roles ON users.role_id = roles.id`
	whereClause := " WHERE 1=1"

	if params.CampaignId != nil && *params.CampaignId != "" {
		baseQuery += " JOIN users_registration ur ON users.id = ur.user_id"
		whereClause += " AND ur.campaign_id = ?"
		args = append(args, *params.CampaignId)

		if params.StartDate != nil {
			whereClause += " AND ur.registration_date >= ?"
			args = append(args, *params.StartDate)
		}
		if params.EndDate != nil {
			whereClause += " AND ur.registration_date <= ?"
			args = append(args, *params.EndDate)
		}
	}

	if params.Role != nil {
		whereClause += " AND users.role_id = ?"
		args = append(args, *params.Role)
	} else {
		args = append(args, constant.RoleIDAdmin, constant.RoleIDStaff)
		whereClause += " AND users.role_id IN (?,?)"
	}

	if params.Name != "" {
		whereClause += " AND users.name LIKE ?"
		args = append(args, "%"+params.Name+"%")
	}
	if params.Email != "" {
		whereClause += " AND users.email LIKE ?"
		args = append(args, "%"+params.Email+"%")
	}
	if params.Status != "" {
		whereClause += " AND users.status = ?"
		args = append(args, params.Status)
	}
	if params.PhoneNumber != "" {
		whereClause += " AND users.phone_number LIKE ?"
		args = append(args, "%"+params.PhoneNumber+"%")
	}
	if params.StartDate != nil {
		whereClause += " AND users.created_at >= ?"
		args = append(args, *params.StartDate)
	}
	if params.EndDate != nil {
		whereClause += " AND users.created_at <= ?"
		args = append(args, *params.EndDate)
	}

	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	err := r.GetConnDb().GetContext(ctx, &totalRecords, countQuery, args...)
	if err != nil {
		return nil, pagination.Pagination{}, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to count users")
	}

	dataQuery := "SELECT users.*, roles.name AS role_name " + baseQuery + whereClause + " ORDER BY users.name ASC LIMIT ? OFFSET ?"
	queryArgs := append(args, perPage, (page-1)*perPage)

	err = r.GetConnDb().SelectContext(ctx, &users, dataQuery, queryArgs...)
	if err != nil {
		return nil, pagination.Pagination{}, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve users")
	}

	totalPages := (totalRecords + perPage - 1) / perPage

	paginate := pagination.Pagination{
		CurrentPage:  page,
		PerPage:      perPage,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
	}

	return users, paginate, nil
}

func (r *Repository) checkForDuplicate(ctx context.Context, tx *sqlx.Tx, column, value string, id *int64) error {
	var exists string
	query := `SELECT 1 FROM users WHERE ` + column + ` = ?`
	args := []interface{}{value}

	if id != nil {
		query += ` AND id != ?`
		args = append(args, *id)
	}

	err := tx.QueryRowContext(ctx, query, args...).Scan(&exists)

	if err == nil {
		return httperror.New(fiber.StatusBadRequest, column+" already exists")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to check duplicate")
	}
	return nil
}

func (r *Repository) checkExistingData(ctx context.Context, tx *sqlx.Tx, email string, phoneNumber *string, id *int64) error {
	if err := r.checkForDuplicate(ctx, tx, "email", email, id); err != nil {
		return err
	}
	if phoneNumber != nil {
		if err := r.checkForDuplicate(ctx, tx, "phone_number", *phoneNumber, id); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int64) (User, error) {
	var result User

	var query string
	var args []interface{}

	query = `SELECT * FROM users WHERE id = ? LIMIT 1;`
	args = append(args, id)

	err := r.GetConnDb().QueryRowxContext(ctx, query, args...).StructScan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, httperror.New(fiber.StatusNotFound, "user not found")
		}
		return User{}, httperror.New(fiber.StatusInternalServerError, err.Error())
	}

	return result, nil
}

func (r *Repository) UpdateUserStatus(ctx context.Context, tx *sqlx.Tx, ids []int64, updatedStatus string, versions []int64) error {
	if len(ids) == 0 {
		return httperror.New(fiber.StatusBadRequest, "no user IDs provided")
	}

	versionConditions := make([]string, len(ids))
	args := make([]interface{}, 0)

	for i, id := range ids {
		versionConditions[i] = "(id = ? AND version = ?)"
		args = append(args, id, versions[i])
	}

	query := fmt.Sprintf(`
		UPDATE users 
		SET status = ?, version = version + 1, updated_at = CURRENT_TIMESTAMP
		WHERE %s`,
		strings.Join(versionConditions, " OR "),
	)

	args = append([]interface{}{constant.GetStatusFromString(updatedStatus)}, args...)

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update users")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return httperror.New(fiber.StatusConflict, "update conflict: users may have been updated by another transaction")
	}

	return nil
}

func (r *Repository) GetEmailVerificationTrialRequestByDate(ctx context.Context, email string, queryDate time.Time, tokenType string) (*int8, error) {
	filterDate := queryDate.Format("2006-01-02")

	query := `SELECT trial FROM email_verifications 
             WHERE email = ? 
             AND DATE(created_at) = ?
             AND type = ?`
	var trial int8 = 0
	err := r.GetConnDb().QueryRowContext(ctx, query, email, filterDate, tokenType).Scan(&trial)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to fetch verification record")
		}
	}

	return &trial, nil
}

func (r *Repository) InsertEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string, token string, expiredAt time.Time, tokenType string) error {
	insertQuery := `
			INSERT INTO email_verifications (email, verification_token, expired_at, type, trial, version)
			VALUES (?, ?, ?, ?, 1, 0)
		`
	_, err := tx.ExecContext(ctx, insertQuery, email, token, expiredAt, tokenType)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err,
			"failed to insert verification record")
	}

	return nil
}

func (r *Repository) UpdateEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string,
	targetDate string, token string, expiredAt time.Time, version int64, tokenType string,
) error {
	updateQuery := `
			UPDATE email_verifications
			SET verification_token = ?, 
			    expired_at = ?, 
			    trial = trial + 1, 
			    updated_at = CURRENT_TIMESTAMP,
			    version = ?
			WHERE email = ? AND DATE(created_at) = ?
			AND version = ?
			AND type = ?
		`
	_, err := tx.ExecContext(ctx, updateQuery, token, expiredAt, version+1, email, targetDate, version, tokenType)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err,
			"failed to update verification record")
	}

	return nil
}

func (r *Repository) GetByVerificationTokenAndEmail(ctx context.Context,
	verificationToken, email string,
) (*EmailVerification, error) {
	query := `SELECT id, email, verification_token, trial, expired_at, created_at, updated_at, version
			  FROM email_verifications 
			  WHERE verification_token = ? AND email = ?`

	var emailVerification EmailVerification
	err := r.GetConnDb().QueryRowxContext(ctx, query, verificationToken, email).StructScan(&emailVerification)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err,
			"failed to fetch verification record")
	}

	return &emailVerification, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string, tokenType string,
) (*EmailVerification, error) {
	query := `SELECT id, email, verification_token, trial, expired_at, created_at, updated_at, version
			  FROM email_verifications 
			  WHERE email = ?
			  AND type = ?
			  ORDER BY created_at DESC LIMIT 1`

	var emailVerification EmailVerification
	err := r.GetConnDb().QueryRowxContext(ctx, query, email, tokenType).StructScan(&emailVerification)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err,
			"failed to fetch verification record")
	}

	return &emailVerification, nil
}

func (r *Repository) DeleteEmailTokenVerification(ctx context.Context, tx *sqlx.Tx, email string) error {
	query := `DELETE FROM email_verifications WHERE email = ?`

	_, err := tx.ExecContext(ctx, query, email)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to delete verification record")
	}

	return nil
}

func (r *Repository) UpdateUserRole(ctx context.Context, tx *sqlx.Tx, id int64, role int64, version int64) error {
	query := `UPDATE users SET role_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	result, err := tx.ExecContext(ctx, query, role, id)
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "failed to update user role")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "failed to update user role")
	}

	if rowsAffected == 0 {
		return httperror.New(fiber.StatusNotFound, "user not found")
	}

	return nil
}

func (r *Repository) UpdatePassword(ctx context.Context, tx *sqlx.Tx, password string, email string, version int64) error {
	hashedPassword, hashErr := hash.HashPassword(password)
	if hashErr != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, hashErr, "failed to hash password")
	}

	query := `UPDATE users SET password = ?, 
                 status = 1,
                 version = ?
                 WHERE email = ?
                 AND version = ?`
	_, err := tx.ExecContext(ctx, query, hashedPassword, version+1, email, version)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update user password")
	}

	return nil
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

func (r *Repository) UpdateUser(ctx context.Context, tx *sqlx.Tx, id int64, name, address, gender, phoneNumber string,
	occupation string, dateOfBirth string, version int64) (*User, error) {
	if err := r.checkForDuplicate(ctx, tx, "phone_number", phoneNumber, &id); err != nil {
		return nil, err
	}

	var dob interface{}

	if dateOfBirth != "" && dateOfBirth != "null" {
		parsedDOB, err := time.Parse("2006-01-02", dateOfBirth)
		if err != nil {
			return nil, httperror.Wrap(fiber.StatusBadRequest, err, "invalid date_of_birth format")
		}
		dob = &parsedDOB
	} else {
		dob = nil
	}

	query := `UPDATE users SET name = ?, phone_number = ?, address = ?, gender = ?,occupation = ?, date_of_birth = ?, 
                 version = ?, updated_at = CURRENT_TIMESTAMP
              WHERE id = ? AND version = ?`
	result, err := tx.ExecContext(ctx, query, name, phoneNumber, address, gender, occupation, dob, version+1, id, version)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve affected rows")
	}
	if rowsAffected == 0 {
		return nil, httperror.New(fiber.StatusConflict, "user was modified by another transaction")
	}

	var user User
	query = `SELECT * FROM users WHERE id = ?`
	err = tx.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve updated user")
	}

	return &user, nil
}

func (r *Repository) GetUserVersions(ctx context.Context, ids []int64) ([]int64, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))

	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("SELECT version FROM users WHERE id IN (%s)",
		strings.Join(placeholders, ","))

	rows, err := r.GetConnDb().QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to get user versions")
	}
	defer rows.Close()

	var versions []int64
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	if len(versions) != len(ids) {
		return nil, httperror.New(fiber.StatusNotFound, "some users not found")
	}

	return versions, nil
}

func (r *Repository) GetTokenEmailVerificationWithType(ctx context.Context, token string,
	tokenType string, emailParam string,
) (string, error) {
	query := `
		SELECT email, expired_at 
		FROM email_verifications 
		WHERE verification_token = ?
		AND type = ?
		AND email = ?
	`
	var email string
	var expiredAt time.Time

	err := r.GetConnDb().QueryRowxContext(ctx, query, token, tokenType, emailParam).Scan(&email, &expiredAt)
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

func (r *Repository) DeleteEmailTokenVerificationByTokenAndType(ctx context.Context, tx *sqlx.Tx, token string, tokenType string) error {
	query := `DELETE FROM email_verifications WHERE verification_token = ? AND type = ?`

	_, err := tx.ExecContext(ctx, query, token, tokenType)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update user password")
	}

	return nil
}

func (r *Repository) DeleteEmailTokenVerificationByToken(ctx context.Context, tx *sqlx.Tx, token string) error {
	query := `DELETE FROM email_verifications WHERE verification_token = ? `

	_, err := tx.ExecContext(ctx, query, token)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update user password")
	}

	return nil
}

func (r *Repository) UpdateUserEmail(ctx context.Context, tx *sqlx.Tx, newEmail string, oldEmail string) error {
	query := `UPDATE users SET email = ?, updated_at = CURRENT_TIMESTAMP WHERE email = ?`

	result, err := tx.ExecContext(ctx, query, newEmail, oldEmail)
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "failed to update user email")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "failed to update user email")
	}

	if rowsAffected == 0 {
		return httperror.New(fiber.StatusNotFound, "user not found")
	}

	return nil
}

func (r *Repository) UpdatePhoto(ctx context.Context, tx *sqlx.Tx, id int64, photo string) error {
	query := `UPDATE users SET photo = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := tx.ExecContext(ctx, query, photo, id)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update user photo")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update user photo")
	}

	if rowsAffected == 0 {
		return httperror.New(fiber.StatusNotFound, "user photo not found")
	}

	return nil
}

func (r *Repository) DeletePhoto(ctx context.Context, tx *sqlx.Tx, id int64) error {
	query := `UPDATE users SET photo = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := tx.ExecContext(ctx, query, nil, id)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update user photo")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update user photo")
	}

	if rowsAffected == 0 {
		return httperror.New(fiber.StatusNotFound, "user photo not found")
	}

	return nil
}

func (r *Repository) GetByVerificationToken(ctx context.Context,
	verificationToken string,
) (*EmailVerification, error) {
	query := `SELECT id, email, verification_token, trial, expired_at, created_at, updated_at, version
			  FROM email_verifications 
			  WHERE verification_token = ?`

	var emailVerification EmailVerification
	err := r.GetConnDb().QueryRowxContext(ctx, query, verificationToken).StructScan(&emailVerification)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err,
			"failed to fetch verification record")
	}

	return &emailVerification, nil
}

func (r *Repository) GetByVerificationTokenAndTokenType(ctx context.Context, verificationToken string, tokenType string,
) (*EmailVerification, error) {
	query := `SELECT id, email, verification_token, trial, expired_at, created_at, updated_at, version
			  FROM email_verifications 
			  WHERE verification_token = ?
			  AND type = ?`

	var emailVerification EmailVerification
	err := r.GetConnDb().QueryRowxContext(ctx, query, verificationToken, tokenType).StructScan(&emailVerification)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err,
			"failed to fetch verification record")
	}

	return &emailVerification, nil
}

func (r *Repository) UpdateUserGoogleId(ctx context.Context, tx *sqlx.Tx, email string, googleId string,
) error {
	updateQuery := `
			UPDATE users SET google_id = ?, updated_at = CURRENT_TIMESTAMP
			WHERE email = ?`

	_, err := tx.ExecContext(ctx, updateQuery, googleId, email)
	if err != nil {
		requestId := ctx.Value("request_id").(string)
		logger.Print(loglevel.ERROR, requestId, "repo_user_method", "UpdateUserGoogleId",
			"failed to update user google id with exception:"+
				err.Error(), nil)
		return httperror.Wrap(fiber.StatusInternalServerError, err,
			"failed to update google id record")
	}

	return nil
}
