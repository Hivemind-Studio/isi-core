package forgotpassword

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	FindByEmail(ctx context.Context, email string) (user.User, error)
}

type userEmailService interface {
	ValidateEmail(ctx context.Context, email string) bool
	HandleTokenGeneration(ctx context.Context, email string, trial int8) (string, error)
	ValidateTrialByDate(ctx context.Context, email string) (*int8, error)
	SendEmail(recipients []string, subject, templatePath string, emailData interface{}) error
}
