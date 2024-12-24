package service

import (
	"authentication-service/internal/domain"
)

type UserServiceInterface interface {
	RegisterUser(input *UserRegisterInput) (*UserResponse, *domain.OperationErrors)
	GetUserByEmail(email string) (*UserResponse, *domain.OperationErrors)
	UpdateUser(input *UserRegisterInput) *domain.OperationErrors
}

type ServiceManager struct {
	UserService UserServiceInterface
}

func NewServiceManager(userService UserServiceInterface) *ServiceManager {
	return &ServiceManager{UserService: userService}
}
