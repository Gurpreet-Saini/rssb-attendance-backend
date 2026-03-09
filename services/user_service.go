package services

import (
	"attendance-system/backend/config"
	"attendance-system/backend/models"
	"attendance-system/backend/repositories"
	"attendance-system/backend/utils"
	"fmt"
)

type UserService interface {
	List(requesterRole string, requesterCenterID *uint, role string) ([]models.User, error)
	Create(name, email, password, role string, centerID *uint) (*models.User, error)
}

type userService struct {
	users  repositories.UserRepository
	audits repositories.AuditRepository
	email  EmailService
	cfg    config.Config
}

func NewUserService(users repositories.UserRepository, audits repositories.AuditRepository, email EmailService, cfg config.Config) UserService {
	return &userService{users: users, audits: audits, email: email, cfg: cfg}
}

func (s *userService) List(requesterRole string, requesterCenterID *uint, role string) ([]models.User, error) {
	if requesterRole == models.RoleCenterAdmin {
		return s.users.List(requesterCenterID, role)
	}
	return s.users.List(nil, role)
}

func (s *userService) Create(name, email, password, role string, centerID *uint) (*models.User, error) {
	if !ValidateRole(role) {
		return nil, fmt.Errorf("invalid role")
	}
	if role != models.RoleSuperAdmin && centerID == nil {
		return nil, fmt.Errorf("center is required")
	}
	if role == models.RoleSuperAdmin {
		centerID = nil
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hash,
		Role:     role,
		CenterID: centerID,
	}
	if err := s.users.Create(user); err != nil {
		return nil, err
	}
	if role == models.RoleCenterAdmin || role == models.RoleOperator {
		portalURL := s.cfg.FrontendURL
		body := fmt.Sprintf(
			"<p>Hello %s,</p><p>Your %s account has been created for the attendance system.</p><p>Email: %s</p><p>Temporary password: %s</p><p>Sign in at <a href=\"%s\">%s</a>.</p>",
			name,
			role,
			email,
			password,
			portalURL,
			portalURL,
		)
		_ = s.email.Send(email, "Your attendance system account", body)
	}
	_ = s.audits.Create(&models.AuditLog{
		Action:     "user_created",
		EntityType: "user",
		EntityID:   fmt.Sprintf("%d", user.ID),
		Details:    email,
	})
	return user, nil
}
