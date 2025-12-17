package usecase

import (
	"context"

	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/repository"
)

type UserListOutput struct {
	Users []*UserOutput `json:"users"`
	Total int64         `json:"total"`
}

type UserUseCase interface {
	GetAllUsers(ctx context.Context, page, limit int) (*UserListOutput, error)
	GetUserByID(ctx context.Context, id string) (*UserOutput, error)
	UpdateUserStatus(ctx context.Context, id string, isActive bool) error
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (uc *userUseCase) GetAllUsers(ctx context.Context, page, limit int) (*UserListOutput, error) {
	users, total, err := uc.userRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	output := make([]*UserOutput, len(users))
	for i, user := range users {
		output[i] = toUserOutput(user)
	}

	return &UserListOutput{
		Users: output,
		Total: total,
	}, nil
}

func (uc *userUseCase) GetUserByID(ctx context.Context, id string) (*UserOutput, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return toUserOutput(user), nil
}

func (uc *userUseCase) UpdateUserStatus(ctx context.Context, id string, isActive bool) error {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if isActive {
		user.Activate()
	} else {
		user.Deactivate()
	}

	return uc.userRepo.Update(ctx, user)
}

// Helper function moved from auth_usecase.go to be reusable
func toUserOutputFromEntity(u *entity.User) *UserOutput {
	return &UserOutput{
		ID:       u.ID,
		Email:    u.Email,
		Nama:     u.Nama,
		Role:     string(u.Role),
		IsActive: u.IsActive,
	}
}
