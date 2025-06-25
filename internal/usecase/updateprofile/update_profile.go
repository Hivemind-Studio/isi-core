package updateprofile

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/constant/loglevel"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) UpdateProfileUser(ctx context.Context, id int64, role string, payload dto.UpdateUserDTO) (*dto.UserDTO, error) {
	tx, err := uc.repoUser.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print(loglevel.INFO, requestId, "Profile service", "UpdateProfile", "function start", payload.Title)

	r := constant.GetRoleID(role)

	if err != nil {
		return nil, httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	findUser, err := uc.repoUser.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if r == constant.RoleIDCoach {
		res, e := uc.repoCoach.UpdateCoach(ctx, tx, findUser.ID, payload.Name, payload.Address, payload.Gender,
			payload.PhoneNumber, payload.DateOfBirth,
			payload.Title, payload.Bio, payload.Expertise, findUser.Version)

		if e != nil {
			return nil, e
		}

		coach := user.ConvertUserToDTO(*res)

		err = tx.Commit()
		if err != nil {
			return nil, httperror.New(fiber.StatusInternalServerError, "Failed to update password")
		}

		return &coach, nil
	}

	res, err := uc.repoUser.UpdateUser(ctx, tx, id, payload.Name, payload.Address, payload.Gender, payload.PhoneNumber,
		payload.Occupation, payload.DateOfBirth, findUser.Version)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, httperror.New(fiber.StatusInternalServerError, "Failed to update password")
	}

	result := user.ConvertUserToDTO(*res)

	return &result, nil
}
