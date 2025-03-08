package getprofileuser

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (s *UseCase) Execute(ctx context.Context, id int64, role string) (*dto.UserDTO, error) {
	module := "get_profile_user_by_id"
	functionName := "Execute"
	requestId := ctx.Value("request_id").(string) // Assuming request_id is stored in the context

	logger.Print("info", requestId, module, functionName,
		fmt.Sprintf("Executing with ID: %d, Role: %s", id, role), nil)

	r := constant.GetRoleID(role)
	logger.Print("info", requestId, module, functionName,
		fmt.Sprintf("Resolved role ID: %d for role: %s", r, role), nil)

	if r == constant.RoleIDCoach {
		logger.Print("info", requestId, module, functionName,
			fmt.Sprintf("Fetching coach details for ID: %d", id), nil)

		result, err := s.repoCoach.GetCoachById(ctx, id)
		if err != nil {
			logger.Print("error", requestId, module, functionName,
				fmt.Sprintf("Failed to get coach by ID: %d, error: %v", id, err), nil)
			return nil, err
		}

		coach := user.ConvertUserToDTO(result)
		logger.Print("info", requestId, module, functionName,
			fmt.Sprintf("Successfully fetched coach details for ID: %d", id), nil)

		return &coach, nil
	}

	logger.Print("info", requestId, module, functionName,
		fmt.Sprintf("Fetching user details for ID: %d", id), nil)

	res, err := s.repoUser.GetUserByID(ctx, id)
	if err != nil {
		logger.Print("error", requestId, module, functionName,
			fmt.Sprintf("User not found, ID: %d, error: %v", id, err), nil)
		return nil, httperror.New(fiber.StatusNotFound, "user not found")
	}

	userDto := user.ConvertUserToDTO(res)

	logger.Print("info", requestId, module, functionName,
		fmt.Sprintf("Successfully fetched user details for ID: %d", id), nil)

	return &userDto, nil
}
