package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Hivemind-Studio/isi-core/internal/constant"
	pagination "github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) Create(ctx context.Context, tx *sqlx.Tx, name string, email string,
	password string, roleId int64, phoneNumber string, status int) (id int64, err error) {
	if err := r.checkExistingData(ctx, tx, email, phoneNumber); err != nil {
		return 0, err
	}

	hashedPassword, hashErr := hash.HashPassword(password)
	if hashErr != nil {
		return 0, httperror.Wrap(fiber.StatusInternalServerError, hashErr, "failed to hash password")
	}

	insertUserQuery := `INSERT INTO users (name, email, password, role_id, phone_number, status, verification) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, insertUserQuery, name, email, hashedPassword, roleId, phoneNumber, status, 0)
	if err != nil {
		return 0, httperror.New(fiber.StatusConflict, "failed to insert user")
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return 0, httperror.New(fiber.StatusInternalServerError, "failed to retrieve last inserted user ID")
	}

	return userId, nil
}

func (r *Repository) CreateStaff(ctx context.Context, tx *sqlx.Tx, user User) (id int64, err error) {
	if err := r.checkExistingData(ctx, tx, user.Email, *user.PhoneNumber); err != nil {
		return 0, err
	}

	hashedPassword, hashErr := hash.HashPassword(user.Password)
	if hashErr != nil {
		return 0, httperror.Wrap(fiber.StatusInternalServerError, hashErr, "failed to hash password")
	}

	insertUserQuery := `INSERT INTO users (name, email, password, role_id, phone_number, status, 
                   address, gender, verification) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, insertUserQuery, user.Name, user.Email, hashedPassword,
		constant.RoleIDStaff, user.PhoneNumber, user.Status, user.Gender, user.Address, 0)
	if err != nil {
		return 0, httperror.New(fiber.StatusConflict, "failed to insert user")
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
) ([]User, pagination.Pagination, error) {
	var users []User
	var totalRecords int64
	var args []interface{}

	baseQuery := `FROM users as users LEFT JOIN roles ON users.role_id = roles.id WHERE`
	if params.Role != nil {
		baseQuery += " u.role_id = ?"
		args = append(args, *params.Role)
	} else {
		args = append(args, constant.RoleIDAdmin, constant.RoleIDStaff)
		baseQuery += " role_id IN (?,?)"
	}

	if params.Name != "" {
		baseQuery += " AND users.name LIKE ?"
		args = append(args, "%"+params.Name+"%")
	}
	if params.Email != "" {
		baseQuery += " AND users.email LIKE ?"
		args = append(args, "%"+params.Email+"%")
	}
	if params.Status != "" {
		baseQuery += " AND users.status = ?"
		args = append(args, params.Status)
	}
	if params.PhoneNumber != "" {
		baseQuery += " AND users.phone_number = ?"
		args = append(args, params.PhoneNumber)
	}
	if params.StartDate != nil {
		baseQuery += " AND users.created_at >= ?"
		args = append(args, *params.StartDate)
	}
	if params.EndDate != nil {
		baseQuery += " AND users.created_at <= ?"
		args = append(args, *params.EndDate)
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.GetConnDb().GetContext(ctx, &totalRecords, countQuery, args...)
	if err != nil {
		return nil, pagination.Pagination{}, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to count users")
	}

	dataQuery := "SELECT users.*, roles.name AS role_name " + baseQuery + " LIMIT ? OFFSET ?"
	queryArgs := append(args, perPage, (page-1)*perPage)

	err = r.GetConnDb().SelectContext(ctx, &users, dataQuery, queryArgs...)
	if err != nil {
		return nil, pagination.Pagination{}, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve users")
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

func (r *Repository) GetUserByID(ctx context.Context, id int64, role *int64) (User, error) {
	var result User

	var query string
	var args []interface{}

	if role != nil {
		if *role == constant.RoleIDCoachee {
			query = `SELECT * 
						FROM users 
						WHERE id = ? AND role_id = ? 
						LIMIT 1;`
			args = append(args, id, *role)
		}
	} else {
		query = `SELECT * 
						FROM users 
						WHERE id = ? AND (role_id = ? OR role_id = ?) 
						LIMIT 1;`
		args = append(args, id, constant.RoleIDAdmin, constant.RoleIDStaff)
	}

	err := r.GetConnDb().QueryRowxContext(ctx, query, args...).StructScan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, httperror.New(fiber.StatusNotFound, "user not found")
		}
		return User{}, httperror.New(fiber.StatusInternalServerError, err.Error())
	}

	return result, nil
}

func (r *Repository) UpdateUserStatus(ctx context.Context, tx *sqlx.Tx, ids []int64, updatedStatus string) error {
	if len(ids) == 0 {
		return httperror.New(fiber.StatusBadRequest, "no user IDs provided")
	}

	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf("UPDATE users SET status = %d WHERE id IN (%s)",
		constant.GetStatusFromString(updatedStatus), strings.Join(placeholders, ","))

	_, err := tx.ExecContext(ctx, query, utils.ToInterfaceSlice(ids)...)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to suspend users")
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

func (r *Repository) GetByVerificationTokenAndEmail(ctx context.Context,
	verificationToken, email string,
) (*EmailVerification, error) {
	query := `SELECT id, email, verification_token, trial, expired_at, created_at, updated_at
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

func (r *Repository) DeleteEmailTokenVerification(ctx context.Context, tx *sqlx.Tx, email string) error {
	query := `DELETE FROM email_verifications WHERE email = ?`

	_, err := tx.ExecContext(ctx, query, email)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update user password")
	}

	return nil
}
