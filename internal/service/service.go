package service

import (
	"authentication-service/internal/domain"
	"net/http"
)

type UserServiceInterface interface {
	RegisterUser(r *http.Request) error
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(r *http.Request) error
}

type ServiceManager struct {
	UserService UserServiceInterface
}

func NewServiceManager(userService UserServiceInterface) *ServiceManager {
	return &ServiceManager{UserService: userService}
}
