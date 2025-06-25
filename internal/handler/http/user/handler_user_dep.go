package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"time"
)

type CreateUserStaffUseCaseInterface interface {
	CreateUserStaff(ctx context.Context, body auth.RegistrationStaffDTO) (err error)
}

type GetUsersUseCaseInterface interface {
	GetCoachees(ctx context.Context, name string, email string, phoneNumber string, status string,
		startDate *time.Time, endDate *time.Time, campaignId string, page int64, perPage int64,
	) ([]user.UserDTO, pagination.Pagination, error)
}

type GetUserByIDUseCaseInterface interface {
	GetCoacheeByID(ctx context.Context, id int64) (result *user.UserDTO, err error)
}

type UpdateUserStatusUseCaseInterface interface {
	UpdateUserStatus(ctx context.Context, ids []int64, status int64) error
}

type UpdateUserRoleUseCaseInterface interface {
	UpdateUserRole(ctx context.Context, id int64, role int64) error
}

type UpdateUserEmailInterface interface {
	VerifyRegistrationToken(ctx context.Context, token string, oldEmail string) (err error)
}

type SendChangeEmailVerificationInterface interface {
	SendVerificationUseCase(ctx context.Context, email string) error
}

type SendConfirmationChangeNewEmailInterface interface {
	SendConfirmationChangeNewEmail(ctx context.Context, token string, newEmail string, oldEmail string) (err error)
}

type ConfirmChangeNewEmail interface {
	VerifyRegistrationToken(ctx context.Context, token string, oldEmail string) (err error)
}
