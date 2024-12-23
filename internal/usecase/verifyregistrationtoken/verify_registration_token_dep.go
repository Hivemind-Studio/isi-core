package verifyregistrationtoken

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	GetByVerificationTokenAndEmail(ctx context.Context, verificationToken, email string) (*user.EmailVerification, error)
}
