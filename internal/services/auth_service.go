package services

import (
	"errors"
	"time"

	"starter-kit-restapi-gofiber/internal/dto"
	"starter-kit-restapi-gofiber/internal/models"
	"starter-kit-restapi-gofiber/pkg/utils"
)

type AuthService struct {
	UserDir  *UserService
	TokenDir *TokenService
	EmailDir *EmailService
}

func NewAuthService(u *UserService, t *TokenService, e *EmailService) *AuthService {
	return &AuthService{UserDir: u, TokenDir: t, EmailDir: e}
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.UserDir.GetUserByEmail(req.Email)
	if err != nil || !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("incorrect email or password")
	}

	return s.generateResponse(user)
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	createReq := &dto.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     utils.RoleUser,
	}
	user, err := s.UserDir.CreateUser(createReq)
	if err != nil {
		return nil, err
	}
	return s.generateResponse(user)
}

func (s *AuthService) Logout(refreshToken string) error {
	tokenDoc, err := s.TokenDir.VerifyToken(refreshToken, utils.TokenTypeRefresh)
	if err != nil {
		return errors.New("token not found")
	}
	// Blacklist or Delete
	// Here we delete
	return s.TokenDir.DB.Delete(&tokenDoc).Error
}

func (s *AuthService) RefreshAuth(refreshToken string) (*dto.AuthResponse, error) {
	tokenDoc, err := s.TokenDir.VerifyToken(refreshToken, utils.TokenTypeRefresh)
	if err != nil {
		return nil, errors.New("please authenticate")
	}

	user, err := s.UserDir.GetUserById(tokenDoc.UserID.String())
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Remove old refresh token
	s.TokenDir.DB.Delete(&tokenDoc)

	return s.generateResponse(user)
}

func (s *AuthService) ForgotPassword(email string) error {
	token, err := s.TokenDir.GenerateResetPasswordToken(email)
	if err != nil {
		return err // Do not reveal if email exists or not in production usually, but sticking to logic
	}
	return s.EmailDir.SendResetPasswordEmail(email, token)
}

func (s *AuthService) ResetPassword(token string, newPassword string) error {
	tokenDoc, err := s.TokenDir.VerifyToken(token, utils.TokenTypeResetPassword)
	if err != nil {
		return errors.New("password reset failed")
	}
	
	user, _ := s.UserDir.GetUserById(tokenDoc.UserID.String())
	req := &dto.UpdateUserRequest{Password: newPassword}
	
	if _, err := s.UserDir.UpdateUserById(user.ID.String(), req); err != nil {
		return err
	}
	
	// Consume token
	s.TokenDir.DB.Delete(&tokenDoc)
	return nil
}

func (s *AuthService) SendVerificationEmail(userID string) error {
	user, err := s.UserDir.GetUserById(userID)
	if err != nil {
		return err
	}
	token, err := s.TokenDir.GenerateVerifyEmailToken(user)
	if err != nil {
		return err
	}
	return s.EmailDir.SendVerificationEmail(user.Email, token)
}

func (s *AuthService) VerifyEmail(token string) error {
	tokenDoc, err := s.TokenDir.VerifyToken(token, utils.TokenTypeVerifyEmail)
	if err != nil {
		return errors.New("email verification failed")
	}

	user, _ := s.UserDir.GetUserById(tokenDoc.UserID.String())
	user.IsEmailVerified = true
	s.UserDir.DB.Save(&user)

	s.TokenDir.DB.Delete(&tokenDoc)
	return nil
}

func (s *AuthService) generateResponse(user *models.User) (*dto.AuthResponse, error) {
	access, refresh, accExp, refExp, err := s.TokenDir.GenerateAuthTokens(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:              user.ID,
			Name:            user.Name,
			Email:           user.Email,
			Role:            user.Role,
			IsEmailVerified: user.IsEmailVerified,
			CreatedAt:       user.CreatedAt,
		},
		Tokens: dto.TokensResponse{
			Access:  dto.TokenDetails{Token: access, Expires: accExp.Format(time.RFC3339)},
			Refresh: dto.TokenDetails{Token: refresh, Expires: refExp.Format(time.RFC3339)},
		},
	}, nil
}