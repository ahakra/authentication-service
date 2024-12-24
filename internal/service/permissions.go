package service

import (
	"authentication-service/internal/data"
	"fmt"
)

type PermissionsService struct {
	RepoManager *data.RepoManager
}

func NewPermissionsService(repoManager *data.RepoManager) *TokenService {
	return &TokenService{RepoManager: repoManager}
}

func (s *PermissionsService) AddPermission(userID int64, permission string) error {
	err := s.RepoManager.PermissionsRepo.InsertPermissions(permission)
	if err != nil {
		return fmt.Errorf("could not add permission: %w", err)
	}

	permissionID, err := s.RepoManager.PermissionsRepo.GetPermissionIDByName(permission)
	if err != nil {
		return fmt.Errorf("could not retrieve permission ID: %w", err)
	}

	err = s.RepoManager.PermissionsRepo.InsertUserPermissions(userID, permissionID)
	if err != nil {
		return fmt.Errorf("could not assign permission to user: %w", err)
	}

	return nil
}

func (s *PermissionsService) RemovePermission(userID int64, permission string) error {
	permissionID, err := s.RepoManager.PermissionsRepo.GetPermissionIDByName(permission)
	if err != nil {
		return fmt.Errorf("could not retrieve permission ID: %w", err)
	}

	err = s.RepoManager.PermissionsRepo.DeleteUserPermissions(userID, permissionID)
	if err != nil {
		return fmt.Errorf("could not remove permission from user: %w", err)
	}

	return nil
}

func (s *PermissionsService) GetPermissionsForUser(userID int64) ([]string, error) {
	permissions, err := s.RepoManager.PermissionsRepo.GetAllForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve user permissions: %w", err)
	}

	return permissions, nil
}
