package adminuser

import (
	"time"

	"pemira-api/internal/shared/constants"
)

type User struct {
	ID          int64          `json:"id"`
	Username    string         `json:"username"`
	Email       string         `json:"email"`
	FullName    string         `json:"full_name"`
	Role        constants.Role `json:"role"`
	VoterID     *int64         `json:"voter_id,omitempty"`
	TPSID       *int64         `json:"tps_id,omitempty"`
	LecturerID  *int64         `json:"lecturer_id,omitempty"`
	StaffID     *int64         `json:"staff_id,omitempty"`
	IsActive    bool           `json:"is_active"`
	LastLoginAt *time.Time     `json:"last_login_at,omitempty"`
	LoginCount  *int           `json:"login_count,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ListFilter struct {
	Search string
	Role   string
	Active *bool
}

type CreateInput struct {
	Username   string         `json:"username"`
	Email      string         `json:"email"`
	FullName   string         `json:"full_name"`
	Role       constants.Role `json:"role"`
	Password   string         `json:"password"`
	VoterID    *int64         `json:"voter_id,omitempty"`
	TPSID      *int64         `json:"tps_id,omitempty"`
	LecturerID *int64         `json:"lecturer_id,omitempty"`
	StaffID    *int64         `json:"staff_id,omitempty"`
	IsActive   *bool          `json:"is_active,omitempty"`
}

type UpdateInput struct {
	Email      *string         `json:"email,omitempty"`
	FullName   *string         `json:"full_name,omitempty"`
	Role       *constants.Role `json:"role,omitempty"`
	VoterID    *int64          `json:"voter_id,omitempty"`
	TPSID      *int64          `json:"tps_id,omitempty"`
	LecturerID *int64          `json:"lecturer_id,omitempty"`
	StaffID    *int64          `json:"staff_id,omitempty"`
	IsActive   *bool           `json:"is_active,omitempty"`
}

type ResetPasswordInput struct {
	NewPassword string `json:"new_password"`
}
