package user_authentication

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserId       uint           `gorm:"primaryKey" json:"user_id"`
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:120;uniqueIndex;not null" json:"email"`
	FullName     string         `gorm:"size:120" json:"full_name"`
	Phone        string         `gorm:"size:20" json:"phone"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         string         `gorm:"size:20;default:EMPLOYEE;index" json:"role"` // ADMIN / EMPLOYEE
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	FailedLogin  uint           `gorm:"default:0" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type RegisterReq struct {
	UserId   uint   `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required,alphanum,min=3"`
	Email    string `json:"email"    binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone"     binding:"omitempty"`
	Password string `json:"password"  binding:"required,min=8"`
	Role     string `json:"role"      binding:"omitempty,oneof=ADMIN EMPLOYEE"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRes struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // detik
	TokenType    string `json:"token_type"` // "Bearer"
}
