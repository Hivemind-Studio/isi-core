package campaign

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

func (r *Repository) Create(ctx context.Context, tx *sqlx.Tx, name, channel, link, campaignId string, status int8,
	startDate, endDate *time.Time) (c Campaign, err error) {
	insertUserQuery := `INSERT INTO campaigns (name, channel, start_date, end_date, link, campaign_id, status) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, insertUserQuery, name, channel, &startDate, &endDate, link, campaignId, status)

	if err != nil {
		return Campaign{}, httperror.New(fiber.StatusConflict, "failed to insert campaign")
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return Campaign{}, httperror.New(fiber.StatusInternalServerError, fmt.Sprintf("failed to get last insert ID: %v", err))
	}

	c = Campaign{
		ID:         lastInsertID,
		Name:       name,
		Channel:    channel,
		Link:       link,
		CampaignID: campaignId,
		Status:     status,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	return c, nil
}

func (r *Repository) Get(ctx context.Context, params campaign.Params, page int64, perPage int64,
) ([]Campaign, pagination.Pagination, error) {
	var campaigns []Campaign
	var totalRecords int64
	var args []interface{}

	baseQuery := "SELECT * FROM campaigns"
	whereConditions := []string{}

	if params.Name != "" {
		whereConditions = append(whereConditions, "name LIKE ?")
		args = append(args, "%"+params.Name+"%")
	}

	if params.Status != "" {
		whereConditions = append(whereConditions, "status = ?")
		args = append(args, params.Status)
	}

	if params.Channel != "" {
		whereConditions = append(whereConditions, "channel = ?")
		args = append(args, params.Channel)
	}

	if params.StartDate != nil {
		whereConditions = append(whereConditions, "created_at >= ?")
		args = append(args, *params.StartDate)
	}

	if params.EndDate != nil {
		whereConditions = append(whereConditions, "created_at <= ?")
		args = append(args, *params.EndDate)
	}

	if len(whereConditions) > 0 {
		baseQuery += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	countQuery := strings.Replace(baseQuery, "SELECT * FROM campaigns", "SELECT COUNT(*) FROM campaigns", 1)
	err := r.GetConnDb().GetContext(ctx, &totalRecords, countQuery, args...)
	if err != nil {
		return nil, pagination.Pagination{}, httperror.Wrap(fiber.StatusInternalServerError,
			err, "failed to count campaigns")
	}

	dataQuery := baseQuery + " ORDER BY name ASC LIMIT ? OFFSET ?"
	queryArgs := append(args, perPage, (page-1)*perPage)
	err = r.GetConnDb().SelectContext(ctx, &campaigns, dataQuery, queryArgs...)
	if err != nil {
		return nil, pagination.Pagination{}, httperror.Wrap(fiber.StatusInternalServerError,
			err, "failed to retrieve campaigns")
	}

	// Calculate pagination info
	totalPages := (totalRecords + perPage - 1) / perPage
	paginate := pagination.Pagination{
		CurrentPage:  page,
		PerPage:      perPage,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
	}

	return campaigns, paginate, nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (Campaign, error) {
	query := "SELECT * FROM campaigns WHERE id = ?"
	var c Campaign

	err := r.GetConnDb().GetContext(ctx, &c, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return Campaign{}, httperror.New(fiber.StatusNotFound, "campaign not found")
		}
		return Campaign{}, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve campaign")
	}

	return c, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, ids []int64, status int8) error {
	if len(ids) == 0 {
		return httperror.New(fiber.StatusBadRequest, "no campaign IDs provided")
	}

	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	query := fmt.Sprintf("SELECT id FROM campaigns WHERE id IN (%s)", idsStr)
	var campaigns []Campaign

	err := r.GetConnDb().SelectContext(ctx, &campaigns, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(fiber.StatusNotFound, "no campaigns found with the provided IDs")
		}
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve campaigns")
	}

	updateQuery := fmt.Sprintf("UPDATE campaigns SET status = ? WHERE id IN (%s)", idsStr)

	_, err = r.GetConnDb().ExecContext(ctx, updateQuery, status)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update campaign status")
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, tx *sqlx.Tx, id int64, name, channel, link *string, status *int8,
	startDate, endDate *time.Time, version int64) (Campaign, error) {
	setClauses := []string{}
	args := []interface{}{}

	if name != nil {
		setClauses = append(setClauses, "name = ?")
		args = append(args, *name)
	}
	if channel != nil {
		setClauses = append(setClauses, "channel = ?")
		args = append(args, *channel)
	}
	if link != nil {
		setClauses = append(setClauses, "link = ?")
		args = append(args, *link)
	}
	if status != nil {
		setClauses = append(setClauses, "status = ?")
		args = append(args, *status)
	}
	if startDate != nil {
		setClauses = append(setClauses, "start_date = ?")
		args = append(args, *startDate)
	}
	if endDate != nil {
		setClauses = append(setClauses, "end_date = ?")
		args = append(args, *endDate)
	}

	if len(setClauses) == 0 {
		return Campaign{}, httperror.New(fiber.StatusBadRequest, "no campaigns found with the provided id")
	}

	query := fmt.Sprintf(
		`UPDATE campaigns SET %s, version = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		strings.Join(setClauses, ", "),
	)
	args = append(args, version, id)

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return Campaign{}, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update campaign")
	}

	var updatedCampaign Campaign
	err = tx.GetContext(ctx, &updatedCampaign, `SELECT * FROM campaigns WHERE id = ?`, id)
	if err != nil {
		return Campaign{}, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to fetch updated campaign")
	}

	return updatedCampaign, nil
}

func (r *Repository) CreateUserCampaign(ctx context.Context, tx *sqlx.Tx, userId int64, campaignId, ipAddress, userAgent string,
	registrationDate string) error {
	insertUserQuery := `INSERT INTO campaign_registrations (user_id, campaign_id, registration_date, ip_address, 
                                user_agent) VALUES (?, ?, ?, ?, ?)`

	_, err := tx.ExecContext(ctx, insertUserQuery, userId, campaignId, registrationDate, ipAddress, userAgent)

	if err != nil {
		return httperror.New(fiber.StatusConflict, err.Error())
	}

	return nil
}

func (r *Repository) GetTotalRegistrants(ctx context.Context, id int64) (int64, error) {
	query := `
        SELECT COUNT(*) 
        FROM campaign_registrations 
        WHERE id = ?
    `
	var count int64
	err := r.GetConnDb().GetContext(ctx, &count, query, id)
	if err != nil {
		return 0, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to count registrations")
	}
	return count, nil
}
