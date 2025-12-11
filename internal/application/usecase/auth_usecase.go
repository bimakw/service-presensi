package usecase

import (
	"context"
	"errors"

	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/repository"
	"github.com/okinn/service-presensi/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrEmailAlreadyExists = errors.New("email sudah terdaftar")
	ErrUserNotActive      = errors.New("user tidak aktif")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
)

type RegisterInput struct {
	Email    string
	Password string
	Nama     string
	Role     string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	Token string      `json:"token"`
	User  *UserOutput `json:"user"`
}

type UserOutput struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Nama     string `json:"nama"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type AuthUseCase interface {
	Register(ctx context.Context, input RegisterInput) (*AuthOutput, error)
	Login(ctx context.Context, input LoginInput) (*AuthOutput, error)
	GetProfile(ctx context.Context, userID string) (*UserOutput, error)
}

type authUseCase struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.JWTManager
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtManager *jwt.JWTManager) AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (uc *authUseCase) Register(ctx context.Context, input RegisterInput) (*AuthOutput, error) {
	existingUser, _ := uc.userRepo.GetByEmail(ctx, input.Email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	role := entity.UserRole(input.Role)
	if !role.IsValid() {
		role = entity.RoleEmployee
	}

	user, err := entity.NewUser(input.Email, input.Password, input.Nama, role)
	if err != nil {
		return nil, err
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := uc.jwtManager.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		Token: token,
		User:  toUserOutput(user),
	}, nil
}

func (uc *authUseCase) Login(ctx context.Context, input LoginInput) (*AuthOutput, error) {
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.ComparePassword(input.Password) {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	token, err := uc.jwtManager.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		Token: token,
		User:  toUserOutput(user),
	}, nil
}

func (uc *authUseCase) GetProfile(ctx context.Context, userID string) (*UserOutput, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return toUserOutput(user), nil
}

func toUserOutput(u *entity.User) *UserOutput {
	return &UserOutput{
		ID:       u.ID,
		Email:    u.Email,
		Nama:     u.Nama,
		Role:     string(u.Role),
		IsActive: u.IsActive,
	}
}
