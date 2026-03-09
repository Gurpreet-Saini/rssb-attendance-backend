package services

import (
	"attendance-system/backend/config"
	"attendance-system/backend/models"
	"attendance-system/backend/repositories"
	"attendance-system/backend/utils"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	Login(email, password string) (map[string]interface{}, error)
	ForgotPassword(email string) error
	ResetPassword(token, password string) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
}

type authService struct {
	users repositories.UserRepository
	email EmailService
	cfg   config.Config
}

func NewAuthService(users repositories.UserRepository, email EmailService, cfg config.Config) AuthService {
	return &authService{users: users, email: email, cfg: cfg}
}

func (s *authService) Login(email, password string) (map[string]interface{}, error) {
	user, err := s.users.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return nil, err
	}
	token, err := utils.GenerateToken(s.cfg.JWTSecret, user.ID, user.Role, user.CenterID, s.cfg.TokenDuration)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"role":      user.Role,
			"center_id": user.CenterID,
			"center":    user.Center,
		},
	}, nil
}

func (s *authService) ForgotPassword(email string) error {
	user, err := s.users.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	token, err := generateResetToken()
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(s.cfg.PasswordResetTTL)
	user.PasswordResetToken = &token
	user.PasswordResetExpires = &expiresAt
	if err := s.users.Save(user); err != nil {
		return err
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.FrontendURL, token)
	body := fmt.Sprintf(
		"<p>Hello %s,</p><p>Use the link below to reset your password.</p><p><a href=\"%s\">Reset password</a></p><p>This link expires in %d minutes.</p>",
		user.Name,
		resetURL,
		int(s.cfg.PasswordResetTTL.Minutes()),
	)
	return s.email.Send(user.Email, "Reset your password", body)
}

func (s *authService) ResetPassword(token, password string) error {
	user, err := s.users.GetByResetToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("invalid reset link")
		}
		return err
	}
	if user.PasswordResetExpires == nil || time.Now().After(*user.PasswordResetExpires) {
		return fmt.Errorf("reset link has expired")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = hash
	user.PasswordResetToken = nil
	user.PasswordResetExpires = nil
	return s.users.Save(user)
}

func (s *authService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.users.GetByID(userID)
	if err != nil {
		return err
	}
	if err := utils.CheckPasswordHash(oldPassword, user.Password); err != nil {
		return fmt.Errorf("old password is incorrect")
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hash
	user.PasswordResetToken = nil
	user.PasswordResetExpires = nil
	return s.users.Save(user)
}

func ValidateRole(role string) bool {
	return role == models.RoleSuperAdmin || role == models.RoleCenterAdmin || role == models.RoleOperator
}

func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
