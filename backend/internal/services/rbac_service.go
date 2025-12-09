package services

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type RBACService interface {
	CheckPermission(userID uint, resource, action string) (bool, error)
	GetUserPermissions(userID uint) ([]*models.Permission, error)
	HasRole(userID uint, roleName string) (bool, error)
	GetUserRole(userID uint) (*models.Role, error)
	CanManageUser(managerID, targetUserID uint) (bool, error)
	GetRoleHierarchyLevel(roleName string) int
}

type RBACServiceImpl struct {
	rbacRepo repository.RBACRepository
}

func NewRBACService(rbacRepo repository.RBACRepository) RBACService {
	return &RBACServiceImpl{
		rbacRepo: rbacRepo,
	}
}

// CheckPermission - Check if user has specific permission
func (s *RBACServiceImpl) CheckPermission(userID uint, resource, action string) (bool, error) {
	user, err := s.rbacRepo.FindUserWithRoleAndPermissions(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("user not found")
		}
		return false, err
	}

	// Build required permission name
	requiredPermission := fmt.Sprintf("%s.%s", resource, action)

	// Check if user has the permission
	for _, permission := range user.Role.Permissions {
		// Check by full permission name (e.g., "user.create")
		if permission.Name == requiredPermission {
			return true, nil
		}

		// Check by resource and action separately
		if permission.Resource == resource && permission.Action == action {
			return true, nil
		}
	}

	return false, nil
}

// GetUserPermissions - Get all permissions for a user
func (s *RBACServiceImpl) GetUserPermissions(userID uint) ([]*models.Permission, error) {
	user, err := s.rbacRepo.FindUserWithRoleAndPermissions(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user.Role.Permissions, nil
}

func (s *RBACServiceImpl) HasRole(userID uint, roleName string) (bool, error) {
	user, err := s.rbacRepo.FindUserWithRole(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("user not found")
		}
		return false, err
	}

	return user.Role.Name == roleName, nil
}

func (s *RBACServiceImpl) GetUserRole(userID uint) (*models.Role, error) {
	user, err := s.rbacRepo.FindUserWithRoleAndPermissions(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user.Role, nil
}

func (s *RBACServiceImpl) CanManageUser(managerID, targetUserID uint) (bool, error) {

	managerRole, err := s.GetUserRole(managerID)
	if err != nil {
		return false, err
	}

	targetRole, err := s.GetUserRole(targetUserID)
	if err != nil {
		return false, err
	}

	managerLevel := s.GetRoleHierarchyLevel(managerRole.Name)
	targetLevel := s.GetRoleHierarchyLevel(targetRole.Name)

	return managerLevel > targetLevel, nil
}

func (s *RBACServiceImpl) GetRoleHierarchyLevel(roleName string) int {
	roleHierarchy := map[string]int{
		"Super Admin": 5,
		"Admin":       4,
		"Manager":     3,
		"Seller":      2,
		"User":        1,
		"Guest":       0,
	}

	level, exists := roleHierarchy[roleName]
	if !exists {
		return 0 
	}

	return level
}
