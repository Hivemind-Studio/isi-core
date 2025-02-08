package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"time"
)

type CreateUserStaffUseCaseInterface interface {
	Execute(ctx context.Context, body auth.RegistrationStaffDTO) (err error)
}

type GetUsersUseCaseInterface interface {
	Execute(ctx context.Context, name string, email string, phoneNumber string, status string,
		startDate *time.Time, endDate *time.Time, page int64, perPage int64,
	) ([]user.UserDTO, pagination.Pagination, error)
}

type GetUserByIDUseCaseInterface interface {
	Execute(ctx context.Context, id int64) (result *user.UserDTO, err error)
}

type UpdateUserStatusUseCaseInterface interface {
	Execute(ctx context.Context, ids []int64, status string) error
}

type UpdateUserRoleUseCaseInterface interface {
	Execute(ctx context.Context, id int64, role int64) error
}

type UpdateUserEmailInterface interface {
	Execute(ctx context.Context, token string, newEmail string, oldEmail string) (err error)
}

type SendChangeEmailVerificationInterface interface {
	Execute(ctx context.Context, email string) error
}

type SendConfirmationChangeNewEmailInterface interface {
	Execute(ctx context.Context, token string, newEmail string, oldEmail string) (err error)
}

type ConfirmChangeNewEmail interface {
	Execute(ctx context.Context, token string, newEmail string, oldEmail string) (err error)
}
