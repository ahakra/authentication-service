package service

import (
	"authentication-service/internal/data"
	"authentication-service/internal/domain"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// RegisterUser registers a new user in the system
type UserService struct {
	RepoManager *data.RepoManager
}

// NewUserService creates a new instance of UserService
func NewUserService(repoManager *data.RepoManager) *UserService {
	return &UserService{RepoManager: repoManager}
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
	user.Password.Set(string(uri.Password), validationErrors)
	user.CreatedAt = time.Now()
	// If there are validation errors, return nil and the errors
	if len(validationErrors.Validation) > 0 {
		return nil, validationErrors
	}

	// Return the created user if no validation errors
	return user, validationErrors
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
		ID:        output.ID,
		Name:      output.Name,
		Email:     output.Email,
		Activated: output.Activated,
		Password:  output.Password,
	}
	return res, nil
}

func (s *UserService) GetUserByID(userId int64) (*UserResponse, *domain.OperationErrors) {
	operationError := domain.OperationErrors{
		Database:   make(map[string][]string),
		Validation: make(map[string][]string),
	}
	output, err := s.RepoManager.UserRepo.GetById(userId)
	if err != nil {
		operationError.Database = make(map[string][]string)
		operationError.AddDatabaseError("Database", err.Error())
		return nil, &operationError
	}

	res := &UserResponse{
		ID:        output.ID,
		Name:      output.Name,
		Email:     output.Email,
		Activated: output.Activated,
		Password:  output.Password,
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

	fromDatabaseUser, err := s.RepoManager.UserRepo.GetByEmail(userModel.Email)

	if err != nil {
		operationError.Database = make(map[string][]string)
		operationError.AddDatabaseError("Database", err.Error())
		return operationError
	}
	pass := domain.Password{
		PasswordHash: fromDatabaseUser.Password,
	}
	isMatch, err := pass.Matches(input.Password)
	if err != nil {
		operationError.AddValidationError("Combination", "Invalid combination")
		return operationError
	}
	fmt.Printf("Match user:%v\n", isMatch)
	if isMatch {
		userModel.ID = fromDatabaseUser.ID
		fmt.Printf("User info to update:%v\n", userModel)
		err = s.RepoManager.UserRepo.Update(&userModel)
		if err != nil {
			operationError.AddDatabaseError("Database", err.Error())
			return operationError
		}

		return nil
	}
	return operationError
}

func (s *UserService) ValidateUser(input RegenerateEmailTokenInput) (*ReGenerateEmailTokenResponse, error) {
	fromDatabaseUser, err := s.RepoManager.UserRepo.GetByEmail(input.Email)
	var response ReGenerateEmailTokenResponse
	response.Email = input.Email

	if err != nil {
		response.IsMatch = false
		return &response, err
	}
	response.ID = fromDatabaseUser.ID

	pass := domain.Password{
		PasswordHash: fromDatabaseUser.Password,
	}
	isMatch, err := pass.Matches(input.Password)
	if err != nil && errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) {
		response.IsMatch = false

		return &response, err
	} else if err != nil {
		return &response, err
	}
	fmt.Printf("Match user:%v\n", isMatch)
	if isMatch {
		response.IsMatch = true
		return &response, nil
	}
	return &response, nil
}

func (s *UserService) UpdateUserActivationStatus(userID int64, status bool) error {
	return s.RepoManager.UserRepo.UpdateUserActivationStatus(userID, status)
}
