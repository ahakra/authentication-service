package service

import (
	"authentication-service/internal/data"
	"authentication-service/internal/domain"
	"fmt"
	"time"
)

// NewUserService creates a new instance of UserService
func NewUserService(repoManager *data.RepoManager) *UserService {
	return &UserService{RepoManager: repoManager}
}

type UserRegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uri *UserRegisterInput) IntoUserDomainModel() (*domain.UserDomainModel, *domain.OperationErrors) {
	validationErrors := &domain.OperationErrors{
		Validation: make(map[string][]string),
		Database:   make(map[string][]string),
	}

	user := &domain.UserDomainModel{}
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
	return user, validationErrors
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// RegisterUser registers a new user in the system
type UserService struct {
	RepoManager *data.RepoManager
}

func (s *UserService) RegisterUser(input *UserRegisterInput) (*UserResponse, *domain.OperationErrors) {

	validateUser, operationError := input.IntoUserDomainModel()

	fmt.Printf("User Domain Model:%v\n", validateUser)
	if len(operationError.Validation) > 0 {
		return nil, operationError
	}
	validateUser.Activated = false
	validateUser.Version = 1
	fmt.Printf("User Domain Model after adding fields:%v\n", validateUser)

	userModel := validateUser.IntoUserModel()
	fmt.Printf("User Model To Insert :%v\n", userModel)
	output, err := s.RepoManager.UserRepo.Insert(&userModel)
	if err != nil {
		operationError.AddDatabaseError("database", err.Error())
		return nil, operationError
	}
	fmt.Printf("Returned User Model from Insert :%v\n", output)
	res := &UserResponse{
		ID:    output.ID,
		Name:  output.Name,
		Email: output.Email,
	}
	return res, nil
}

// GetUserByEmail retrieves a user by email from the system
func (s *UserService) GetUserByEmail(email string) (*UserResponse, *domain.OperationErrors) {
	operationError := domain.OperationErrors{
		Database:   make(map[string][]string),
		Validation: make(map[string][]string),
	}
	output, err := s.RepoManager.UserRepo.GetByEmail(email)
	if err != nil {
		operationError.Database = make(map[string][]string)
		operationError.AddDatabaseError("Database", err.Error())
		return nil, &operationError
	}

	res := &UserResponse{
		ID:    output.ID,
		Name:  output.Name,
		Email: output.Email,
	}
	return res, nil
}

// UpdateUser updates an existing user in the system
func (s *UserService) UpdateUser(input *UserRegisterInput) *domain.OperationErrors {
	validateUser, operationError := input.IntoUserDomainModel()

	if len(operationError.Validation) > 0 {
		return operationError
	}

	userModel := validateUser.IntoUserModel()
	err := s.RepoManager.UserRepo.Update(&userModel)
	if err != nil {
		operationError.AddDatabaseError("Database", err.Error())
		return operationError
	}
	validateUser.Version = validateUser.Version + 1
	return nil
}
