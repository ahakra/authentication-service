package service

import (
	"authentication-service/internal/data"
	"authentication-service/internal/domain"
	"time"
)

type UserServiceInterface interface {
	RegisterUser(input *UserRegisterInput) (*UserResponse, *domain.OperationErrors)
	GetUserByEmail(email string) (*UserResponse, *domain.OperationErrors)
	UpdateUser(input *UserRegisterInput) *domain.OperationErrors
}

type TokenServiceInterface interface {
	CreateAccessToken(userID int64, scope string, ttl time.Duration, secret string) (string, error)
	ValidateToken(tokenString string, secret string) (bool, error)
	GetTokensForUser(userID int64) ([]data.Token, error)
	GetTokensForUserAndScope(userID int64, scope string) ([]data.Token, error)
	DeleteToken(tokenHash []byte) error
}
type PermissionsServiceInterface interface {
}
type ServiceManager struct {
	UserService        UserServiceInterface
	TokenService       TokenServiceInterface
	PermissionsService PermissionsServiceInterface
}

func NewServiceManager(userService UserServiceInterface, tokenService TokenServiceInterface, permissionsService PermissionsServiceInterface) *ServiceManager {
	return &ServiceManager{
		UserService:        userService,
		TokenService:       tokenService,
		PermissionsService: permissionsService,
	}
}
