package campaign

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"time"
)

func (r *Repository) Create(ctx context.Context, tx *sqlx.Tx, name, channel, link string, startDate, endDate time.Time, status bool) (err error) {
	insertUserQuery := `INSERT INTO campaign (name, channel, start_date, end_date, link, status) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, insertUserQuery, name, channel, startDate, endDate, link, status)

	if err != nil {
		return httperror.New(fiber.StatusConflict, "failed to insert campaign")
	}

	return nil
}
