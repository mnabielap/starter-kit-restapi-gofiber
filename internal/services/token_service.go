package services

import (
	"errors"
	"time"

	"starter-kit-restapi-gofiber/internal/config"
	"starter-kit-restapi-gofiber/internal/models"
	"starter-kit-restapi-gofiber/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenService struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewTokenService(db *gorm.DB, cfg *config.Config) *TokenService {
	return &TokenService{DB: db, Config: cfg}
}

// GenerateAuthTokens (Access + Refresh)
func (s *TokenService) GenerateAuthTokens(userID uuid.UUID) (string, string, time.Time, time.Time, error) {
	// Access Token
	accessExp := time.Duration(s.Config.JWTAccessExpirationMinutes) * time.Minute
	accessToken, err := utils.GenerateToken(userID, accessExp, s.Config.JWTSecret, utils.TokenTypeAccess)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}
	accessExpTime := time.Now().Add(accessExp)

	// Refresh Token
	refreshExp := time.Duration(s.Config.JWTRefreshExpirationDays) * 24 * time.Hour
	refreshToken, err := utils.GenerateToken(userID, refreshExp, s.Config.JWTSecret, utils.TokenTypeRefresh)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}
	refreshExpTime := time.Now().Add(refreshExp)

	// Save Refresh Token
	err = s.SaveToken(refreshToken, userID, refreshExpTime, utils.TokenTypeRefresh)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return accessToken, refreshToken, accessExpTime, refreshExpTime, nil
}

func (s *TokenService) SaveToken(token string, userID uuid.UUID, expiresAt time.Time, tokenType string) error {
	tokenDoc := models.Token{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Type:      tokenType,
		Blacklisted: false,
	}
	return s.DB.Create(&tokenDoc).Error
}

func (s *TokenService) VerifyToken(token string, tokenType string) (*models.Token, error) {
	payload, err := utils.ValidateToken(token, s.Config.JWTSecret)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	
	sub, _ := payload.Claims.GetSubject()
	userID, _ := uuid.Parse(sub)

	var tokenDoc models.Token
	err = s.DB.Where("token = ? AND type = ? AND user_id = ? AND blacklisted = ?", token, tokenType, userID, false).First(&tokenDoc).Error
	if err != nil {
		return nil, errors.New("token not found or blacklisted")
	}

	return &tokenDoc, nil
}

func (s *TokenService) GenerateResetPasswordToken(email string) (string, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("no user found with this email")
	}

	expires := time.Duration(s.Config.JWTResetPwdExpirationMin) * time.Minute
	token, err := utils.GenerateToken(user.ID, expires, s.Config.JWTSecret, utils.TokenTypeResetPassword)
	if err != nil {
		return "", err
	}
	
	s.SaveToken(token, user.ID, time.Now().Add(expires), utils.TokenTypeResetPassword)
	return token, nil
}

func (s *TokenService) GenerateVerifyEmailToken(user *models.User) (string, error) {
	expires := time.Duration(s.Config.JWTVerifyEmailExpirationMin) * time.Minute
	token, err := utils.GenerateToken(user.ID, expires, s.Config.JWTSecret, utils.TokenTypeVerifyEmail)
	if err != nil {
		return "", err
	}

	s.SaveToken(token, user.ID, time.Now().Add(expires), utils.TokenTypeVerifyEmail)
	return token, nil
}