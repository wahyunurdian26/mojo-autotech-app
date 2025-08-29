package user_authentication

import (
	"context"
	"errors"

	"mojo-autotech/constant"
	"mojo-autotech/utils"

	"mojo-autotech/model/user_authentication"
	entity "mojo-autotech/model/user_authentication"
)

type (
	LoginReq    = entity.LoginReq
	LoginRes    = entity.LoginRes
	RegisterReq = entity.RegisterReq
	User        = entity.User
)

type IAuthService interface {
	Login(ctx context.Context, req LoginReq) (LoginRes, error)
	CreateAccount(ctx context.Context, req RegisterReq) (res User, err error)
}

type AuthService struct {
	user_authentication user_authentication.IAuthRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		user_authentication: user_authentication.NewAuthRepository(),
	}
}

func (s *AuthService) Login(ctx context.Context, req LoginReq) (LoginRes, error) {
	user, err := s.user_authentication.Login(ctx, req)
	if err != nil {
		return LoginRes{}, err
	}

	// Status aktif?
	if !user.IsActive {
		return LoginRes{}, errors.New("akun tidak aktif")
	}

	accessToken, expiresIn, err := utils.GenerateAccessToken(user.ID, user.Role, constant.AccessTTL)
	if err != nil {
		return LoginRes{}, err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID, constant.RefreshTTL)
	if err != nil {
		return LoginRes{}, err
	}

	return LoginRes{
		Username:     user.Username,
		Role:         user.Role,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}

func (s *AuthService) CreateAccount(ctx context.Context, req RegisterReq) (User, error) {

	if err := utils.ValidateCreateAccount(req); err != nil {
		return User{}, err
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return User{}, err
	}

	req.Password = hash

	created, err := s.user_authentication.CreateUser(ctx, req)
	if err != nil {
		return User{
			Username: req.Username,
			Email:    req.Email,
		}, err
	}
	return created, nil
}
