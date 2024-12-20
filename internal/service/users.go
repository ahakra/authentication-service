package service

import (
	"authentication-service/internal/data"
	"authentication-service/internal/domain"
	"context"
)

type UserRegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterOutput struct {
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
func (s *UserService) RegisterUser(ctx context.Context, input *UserRegisterInput) (*UserRegisterOutput, *domain.OperationErrors) {

	user, operationError := domain.CreateUserFromUserRegisterInput(*input)

	if len(operationError.Validation) > 0 {
		return nil, operationError
	}
	userModel := data.CreateUserModelFromUserDomain(user)

	err := s.RepoManager.UserRepo.Insert(userModel)
	if err != nil {
		operationError.Database = make(map[string][]string)
		operationError.AddDatabaseError("Database", err.Error())
		return nil, operationError
	}
	output := UserRegisterOutput{
		Name:  user.Name.Inner_value,
		Email: user.Email.Inner_value,
	}
	return &output, nil
}

// GetUserByEmail retrieves a user by email from the system
func (s *UserService) GetUserByEmail(email string) (*data.UserModel, error) {
	// You could add additional logic here, e.g., check if the user is activated
	return s.RepoManager.UserRepo.GetByEmail(email)
}

// UpdateUser updates an existing user in the system
func (s *UserService) UpdateUser(user *data.User) error {
	// Add any business logic here before updating the user, such as validation or logging
	return s.RepoManager.UserRepo.Update(user)
}
