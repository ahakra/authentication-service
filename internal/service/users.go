package service

import (
	"authentication-service/internal/data"
	"time"
)

type UserService struct {
	RepoManager *data.RepoManager
}

// NewUserService creates a new instance of UserService
func NewUserService(repoManager *data.RepoManager) *UserService {
	return &UserService{RepoManager: repoManager}
}

type UserRegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uri *UserRegisterInput) IntoUserModel() (*data.UserModel, *OperationErrors) {
	validationErrors := &OperationErrors{
		Validation: make(map[string][]string),
		Database:   make(map[string][]string),
	}

	user := &data.UserModel{}
	// Validate and set the name
	user.Name.Set(uri.Name, validationErrors)

	// Validate and set the email
	user.Email.Set(uri.Email, validationErrors)

	// Validate and set the password
	user.Password.Set(uri.Password, validationErrors)
	user.CreatedAt = time.Now()
	// If there are validation errors, return nil and the errors
	if len(validationErrors.Validation) > 0 {
		return nil, validationErrors
	}

	// Return the created user if no validation errors
	return user, nil
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// RegisterUser registers a new user in the system
func (s *UserService) RegisterUser(input *UserRegisterInput) (*UserResponse, *OperationErrors) {

	validateUser, operationError := input.IntoUserModel()

	if len(operationError.Validation) > 0 {
		return nil, operationError
	}
	validateUser.Activated = false
	output, err := s.RepoManager.UserRepo.Insert(validateUser)
	if err != nil {
		operationError.AddDatabaseError("database", err.Error())
		return nil, operationError
	}

	return output.IntoUserResponse(), nil
}

// GetUserByEmail retrieves a user by email from the system
func (s *UserService) GetUserByEmail(email string) (*UserResponse, *OperationErrors) {
	operationError := OperationErrors{
		Database:   make(map[string][]string),
		Validation: make(map[string][]string),
	}
	output, err := s.RepoManager.UserRepo.GetByEmail(email)
	if err != nil {
		operationError.Database = make(map[string][]string)
		operationError.AddDatabaseError("Database", err.Error())
		return nil, &operationError
	}

	return output.IntoUserResponse(), nil
}

// UpdateUser updates an existing user in the system
func (s *UserService) UpdateUser(input *UserRegisterInput) *OperationErrors {
	validateUser, operationError := input.IntoUserModel()

	if len(operationError.Validation) > 0 {
		return operationError
	}

	err := s.RepoManager.UserRepo.Update(validateUser)
	if err != nil {
		operationError.AddDatabaseError("Database", err.Error())
		return operationError
	}
	validateUser.Version = validateUser.Version + 1
	return nil
}
