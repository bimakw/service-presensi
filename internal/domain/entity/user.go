package entity

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("email tidak valid")
	ErrInvalidPassword = errors.New("password minimal 6 karakter")
	ErrInvalidName     = errors.New("nama tidak boleh kosong")
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleEmployee UserRole = "employee"
)

func (r UserRole) IsValid() bool {
	return r == RoleAdmin || r == RoleEmployee
}

type User struct {
	ID        string
	Email     string
	Password  string
	Nama      string
	Role      UserRole
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, password, nama string, role UserRole) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if len(password) < 6 {
		return nil, ErrInvalidPassword
	}
	if nama == "" {
		return nil, ErrInvalidName
	}
	if !role.IsValid() {
		role = RoleEmployee
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		Email:     email,
		Password:  string(hashedPassword),
		Nama:      nama,
		Role:      role,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) UpdatePassword(newPassword string) error {
	if len(newPassword) < 6 {
		return ErrInvalidPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}
