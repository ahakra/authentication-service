package service

type UserServiceInterface interface {
	RegisterUser(input *UserRegisterInput) (*UserResponse, *OperationErrors)
	GetUserByEmail(email string) (*UserResponse, *OperationErrors)
	UpdateUser(input *UserRegisterInput) *OperationErrors
}

type ServiceManager struct {
	UserService UserServiceInterface
}

func NewServiceManager(userService UserServiceInterface) *ServiceManager {
	return &ServiceManager{UserService: userService}
}
