package domain

import (
	"authentication-service/internal/data"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
const nameRegex = `^[a-zA-Z\s]+$`

type UserDomainModel struct {
	ID        int64
	CreatedAt time.Time
	Name      name
	Email     email
	Password  Password
	Activated bool
	Version   int
}

func (um *UserDomainModel) IntoUserModel() data.UserModel {
	userResponse := data.UserModel{
		ID:        um.ID,
		CreatedAt: um.CreatedAt,
		Name:      um.Name.Inner_value,
		Email:     um.Email.Inner_value,
		Password:  um.Password.PasswordHash,
	}
	return userResponse
}

type name struct {
	Inner_value string
}

// Set method for setting the name with validation.
func (n *name) Set(name string, ve *OperationErrors) {
	// Check if the name is empty
	if name == "" {
		ve.AddValidationError("name", "name must not be empty")
		return
	}

	// Check if the name length is reasonable (between 1 and 100 characters)
	if len(name) < 1 || len(name) > 100 {
		ve.AddValidationError("name", "name length must be between 1 and 100 characters")
		return
	}

	// Optionally, check for any invalid characters (e.g., numbers, special symbols)
	if !isValidName(name) {
		ve.AddValidationError("name", "name contains invalid characters")
		return
	}

	// Set the name if all checks pass
	n.Inner_value = name
}

// Helper function to check if the name contains valid characters.
func isValidName(name string) bool {
	// Simple regex to allow only alphabets, spaces, and basic punctuation (you can adjust this)
	namePattern := nameRegex
	re := regexp.MustCompile(namePattern)
	return re.MatchString(name)
}

type email struct {
	Inner_value string
}

func (e *email) Set(email string, ve *OperationErrors) {
	// Check if email is empty
	if email == "" {
		ve.AddValidationError("email", "email must not be empty")
		return
	}

	// Check if email format is valid
	if !isValidEmail(email) {
		ve.AddValidationError("email", "invalid email format")
		return
	}

	// Check if email length is reasonable (between 5 and 255 characters)
	if len(email) < 5 || len(email) > 255 {
		ve.AddValidationError("email", "email length must be between 5 and 255 characters")
		return
	}

	// If all checks pass, set the email
	e.Inner_value = email
}

func isValidEmail(email string) bool {
	emailPattern := emailRegex
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}

type Password struct {
	passwordText *string
	PasswordHash []byte
}

func (p *Password) Set(plainTextPassword string, ve *OperationErrors) {
	if len(plainTextPassword) < 8 {
		ve.AddValidationError("password", "password must be at least 8 characters long")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		ve.AddValidationError("password", "error hashing the password")
		return
	}
	p.passwordText = &plainTextPassword
	p.PasswordHash = hash
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.PasswordHash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
