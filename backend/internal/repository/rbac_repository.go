package repository

import (
	"e-commerce/backend/internal/models"

	"gorm.io/gorm"
)

type RBACRepository interface {
	FindUserWithRole(userID uint) (*models.User, error)
	FindUserWithRoleAndPermissions(userID uint) (*models.User, error)
	FindRoleByName(roleName string) (*models.Role, error)
	FindRoleByID(roleID uint) (*models.Role, error)
	FindPermissionByName(permissionName string) (*models.Permission, error)
	FindPermissionsByRoleID(roleID uint) ([]*models.Permission, error)
}

type RBACRepositoryImpl struct {
	db *gorm.DB
}

func NewRBACRepository(db *gorm.DB) RBACRepository {
	return &RBACRepositoryImpl{db: db}
}

// FindUserWithRole - Get user with role data
func (r *RBACRepositoryImpl) FindUserWithRole(userID uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserWithRoleAndPermissions - Get user with role and permissions
func (r *RBACRepositoryImpl) FindUserWithRoleAndPermissions(userID uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role.Permissions").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindRoleByName - Find role by name
func (r *RBACRepositoryImpl) FindRoleByName(roleName string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", roleName).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindRoleByID - Find role by ID with permissions
func (r *RBACRepositoryImpl) FindRoleByID(roleID uint) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, roleID).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindPermissionByName - Find permission by name
func (r *RBACRepositoryImpl) FindPermissionByName(permissionName string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("name = ?", permissionName).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindPermissionsByRoleID - Get all permissions for a role
func (r *RBACRepositoryImpl) FindPermissionsByRoleID(roleID uint) ([]*models.Permission, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, roleID).Error
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}
