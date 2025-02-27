package campaign

import (
	"context"
	"database/sql"
	"github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

func (r *Repository) Create(ctx context.Context, tx *sqlx.Tx, name, channel, link string, status int8,
	startDate, endDate *time.Time) (err error) {
	insertUserQuery := `INSERT INTO campaign (name, channel, start_date, end_date, link, status) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, insertUserQuery, name, channel, &startDate, &endDate, link, status)

	if err != nil {
		return httperror.New(fiber.StatusConflict, "failed to insert campaign")
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, params campaign.Params, page int64, perPage int64,
) ([]Campaign, pagination.Pagination, error) {
	var campaigns []Campaign
	var totalRecords int64
	var args []interface{}

	baseQuery := "SELECT * FROM campaign"

	if params.Name != "" || params.Status != "" || params.StartDate != nil || params.EndDate != nil {
		baseQuery += " WHERE"
	}

	if params.Name != "" {
		baseQuery += " name LIKE ?"
		args = append(args, "%"+params.Name+"%")
		if params.Status != "" || params.StartDate != nil || params.EndDate != nil {
			baseQuery += " AND"
		}
	}
	if params.Status != "" {
		baseQuery += " status = ?"
		args = append(args, params.Status)
		if params.StartDate != nil || params.EndDate != nil {
			baseQuery += " AND"
		}
	}
	if params.StartDate != nil {
		baseQuery += " created_at >= ?"
		args = append(args, *params.StartDate)
		if params.EndDate != nil {
			baseQuery += " AND"
		}
	}
	if params.EndDate != nil {
		baseQuery += " created_at <= ?"
		args = append(args, *params.EndDate)
	}

	countQuery := "SELECT COUNT(*) " + strings.Replace(baseQuery, "SELECT * FROM campaign", "FROM campaign", 1)
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
	query := "SELECT * FROM campaign WHERE id = ?"
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

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := "SELECT * FROM campaign WHERE id = ?"
	var c Campaign

	err := r.GetConnDb().GetContext(ctx, &c, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(fiber.StatusNotFound, "campaign not found")
		}
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve campaign")
	}

	deleteQuery := "DELETE FROM campaign WHERE id = ?"
	_, err = r.GetConnDb().ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to delete campaign")
	}

	return nil
}
