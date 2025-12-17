package entity

import (
	"testing"
)

func TestUserRole_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{"admin role is valid", RoleAdmin, true},
		{"employee role is valid", RoleEmployee, true},
		{"invalid role", UserRole("invalid"), false},
		{"empty role", UserRole(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsValid(); got != tt.expected {
				t.Errorf("UserRole.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		nama        string
		role        UserRole
		wantErr     bool
		expectedErr error
	}{
		{
			name:     "valid user creation",
			email:    "test@example.com",
			password: "password123",
			nama:     "Test User",
			role:     RoleEmployee,
			wantErr:  false,
		},
		{
			name:     "valid admin creation",
			email:    "admin@example.com",
			password: "adminpass",
			nama:     "Admin User",
			role:     RoleAdmin,
			wantErr:  false,
		},
		{
			name:        "empty email",
			email:       "",
			password:    "password123",
			nama:        "Test User",
			role:        RoleEmployee,
			wantErr:     true,
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "short password",
			email:       "test@example.com",
			password:    "12345",
			nama:        "Test User",
			role:        RoleEmployee,
			wantErr:     true,
			expectedErr: ErrInvalidPassword,
		},
		{
			name:        "empty name",
			email:       "test@example.com",
			password:    "password123",
			nama:        "",
			role:        RoleEmployee,
			wantErr:     true,
			expectedErr: ErrInvalidName,
		},
		{
			name:     "invalid role defaults to employee",
			email:    "test@example.com",
			password: "password123",
			nama:     "Test User",
			role:     UserRole("invalid"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.password, tt.nama, tt.role)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error, got nil")
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("NewUser() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error = %v", err)
				return
			}

			if user.Email != tt.email {
				t.Errorf("NewUser() email = %v, want %v", user.Email, tt.email)
			}
			if user.Nama != tt.nama {
				t.Errorf("NewUser() nama = %v, want %v", user.Nama, tt.nama)
			}
			if !user.IsActive {
				t.Error("NewUser() user should be active by default")
			}
			if user.CreatedAt.IsZero() {
				t.Error("NewUser() CreatedAt should not be zero")
			}
			if user.UpdatedAt.IsZero() {
				t.Error("NewUser() UpdatedAt should not be zero")
			}
		})
	}
}

func TestUser_ComparePassword(t *testing.T) {
	user, err := NewUser("test@example.com", "password123", "Test User", RoleEmployee)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"correct password", "password123", true},
		{"incorrect password", "wrongpassword", false},
		{"empty password", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := user.ComparePassword(tt.password); got != tt.expected {
				t.Errorf("User.ComparePassword() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	user, err := NewUser("test@example.com", "password123", "Test User", RoleEmployee)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	oldUpdatedAt := user.UpdatedAt

	tests := []struct {
		name        string
		newPassword string
		wantErr     bool
	}{
		{"valid new password", "newpassword123", false},
		{"short password", "12345", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.UpdatePassword(tt.newPassword)

			if tt.wantErr {
				if err == nil {
					t.Error("UpdatePassword() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdatePassword() unexpected error = %v", err)
				return
			}

			if !user.ComparePassword(tt.newPassword) {
				t.Error("UpdatePassword() new password should match")
			}
			if !user.UpdatedAt.After(oldUpdatedAt) {
				t.Error("UpdatePassword() UpdatedAt should be updated")
			}
		})
	}
}

func TestUser_Activate_Deactivate(t *testing.T) {
	user, err := NewUser("test@example.com", "password123", "Test User", RoleEmployee)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if !user.IsActive {
		t.Error("User should be active by default")
	}

	user.Deactivate()
	if user.IsActive {
		t.Error("User should be inactive after Deactivate()")
	}

	user.Activate()
	if !user.IsActive {
		t.Error("User should be active after Activate()")
	}
}
