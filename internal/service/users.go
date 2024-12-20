package service

import (
	"authentication-service/internal/data"
	"authentication-service/internal/domain"
	"context"
	"time"
)

type UserRegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserService implements the UserServiceInterface and contains the business logic
type UserService struct {
	RepoManager data.RepoManager
}

// NewUserService creates a new instance of UserService
func NewUserService(repoManager data.RepoManager) *UserService {
	return &UserService{RepoManager: repoManager}
}

// RegisterUser registers a new user in the system
func (s *UserService) RegisterUser(ctx context.Context, input *UserRegisterInput) (*UserResponse, *domain.OperationErrors) {

	// Create a new user
	user := &domain.User{
		CreatedAt: time.Now(),
	}
	validateUser, operationError := domain.FromServiceUserRegisterInput(*input)

	if len(operationError.Validation) > 0 {
		return nil, operationError
	}
	userModel := &data.UserModel{}
	userModel.FromUserDomain(validateUser)

	err := s.RepoManager.UserRepo.Insert(userModel)
	if err != nil {
		operationError.Database = make(map[string][]string)
		operationError.AddDatabaseError("Database", err.Error())
		return nil, operationError
	}
	output := UserResponse{
		Name:  user.Name.Inner_value,
		Email: user.Email.Inner_value,
	}
	return &output, nil
}

// GetUserByEmail retrieves a user by email from the system
func (s *UserService) GetUserByEmail(email string) (*UserResponse, *domain.OperationErrors) {
	operationError := domain.OperationErrors{
		Database:   make(map[string][]string),
		Validation: make(map[string][]string),
	}
	userModel, err := s.RepoManager.UserRepo.GetByEmail(email)
	if err != nil {
		operationError.Database = make(map[string][]string)
		operationError.AddDatabaseError("Database", err.Error())
		return nil, &operationError
	}
	output := UserResponse{
		Name:  userModel.Name,
		Email: userModel.Email,
	}
	return &output, nil
}

// UpdateUser updates an existing user in the system
func (s *UserService) UpdateUser(input *UserRegisterInput) (*UserResponse, *domain.OperationErrors) {
	user := &domain.User{
		CreatedAt: time.Now(),
	}
	validateUser, operationError := domain.FromServiceUserRegisterInput(*input)

	if len(operationError.Validation) > 0 {
		return nil, operationError
	}
	userModel := &data.UserModel{}
	userModel.FromUserDomain(validateUser)

	err := s.RepoManager.UserRepo.Update(userModel)
	if err != nil {
		operationError.Database = make(map[string][]string)
		operationError.AddDatabaseError("Database", err.Error())
		return nil, operationError
	}
	output := UserResponse{
		Name:  user.Name.Inner_value,
		Email: user.Email.Inner_value,
	}
	return &output, nil
}
