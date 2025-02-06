package verifyregistrationtoken

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	GetByVerificationToken(ctx context.Context, verificationToken string) (*user.EmailVerification, error)
}
