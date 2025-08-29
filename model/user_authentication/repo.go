package user_authentication

import (
	"context"
	"errors"
	"mojo-autotech/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Interface
type IAuthRepository interface {
	Login(ctx context.Context, req LoginReq) (user User, err error)
	CreateUser(ctx context.Context, req RegisterReq) (User, error)
}

// Impl
type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository() IAuthRepository {
	return &AuthRepository{db: config.NewDB()}
}

const (
	SelectUserByUsername = `
SELECT id, username, email, full_name, phone, password_hash, role, is_active,
       last_login_at, failed_login, created_at, updated_at
FROM users
WHERE LOWER(username) = ? 
LIMIT 1;
`

	FailedLogin = `
UPDATE users SET failed_login = failed_login + 1
WHERE id = ?;
`

	CheckDuplicateUser = `
SELECT COUNT(1)
FROM users
WHERE LOWER(username) = ? OR LOWER(email) = ?;
`
	InsertUser = `
INSERT INTO users
  (user_id,username, email, full_name, phone, password_hash, role, is_active, created_at, updated_at)
VALUES
  (?,?, ?, ?, ?, ?, COALESCE(NULLIF(?, ''), 'EMPLOYEE'), TRUE, NOW(), NOW())
RETURNING
  id,
  user_id,
  username,
  email,
  full_name,
  phone,
  role,
  is_active,
  created_at,
  updated_at;
`
)

func (r *AuthRepository) Login(ctx context.Context, req LoginReq) (user User, err error) {

	tx := r.db.WithContext(ctx).Raw(SelectUserByUsername, req.Username).Scan(&user)
	if tx.Error != nil {
		return User{}, errors.New("gagal mengambil data user")
	}
	if tx.RowsAffected == 0 {
		return User{}, errors.New("username tidak ditemukan")
	}

	// Verifikasi password (bcrypt)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		_ = r.db.WithContext(ctx).Exec(FailedLogin, user.ID).Error
		return User{}, errors.New("password salah")
	}

	return
}

func (a *AuthRepository) CreateUser(ctx context.Context, req RegisterReq) (res User, err error) {
	err = a.db.Raw(InsertUser,
		req.UserId,
		req.Username,
		req.Email,
		req.FullName,
		req.Phone,
		req.Password,
		req.Role,
	).Scan(&res).Error
	return
}
