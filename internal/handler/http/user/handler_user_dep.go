package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"time"
)

type serviceUserInterface interface {
	Create(ctx context.Context, body *auth.RegistrationDTO) (result *auth.RegisterResponse, err error)
	GetUsers(ctx context.Context, name string, email string, startDate, endDate *time.Time, page int64, perPage int64,
	) ([]user.UserDTO, error)
	GetUserByID(ctx context.Context, id int64) (result *user.UserDTO, err error)
	UpdateUserStatus(ctx context.Context, ids []int64, status string) (err error)
}
