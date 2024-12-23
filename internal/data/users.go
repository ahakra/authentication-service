package data

import (
	"authentication-service/internal/service"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

// UserRepository struct holds a reference to the database connection
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository with the given DB connection
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$\n`
const nameRegex = `^[a-zA-Z\s]+$`

type UserModel struct {
	ID        int64
	CreatedAt time.Time
	Name      name
	Email     email
	Password  password
	Activated bool
	Version   int
}

func (um *UserModel) IntoUserResponse() *service.UserResponse {
	userResponse := &service.UserResponse{
		ID:    um.ID,
		Name:  um.Name.Inner_value,
		Email: um.Email.Inner_value,
	}
	return userResponse
}

type name struct {
	Inner_value string
}

// Set method for setting the name with validation.
func (n *name) Set(name string, ve *service.OperationErrors) {
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

func (e *email) Set(email string, ve *service.OperationErrors) {
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

type password struct {
	passwordText *string
	PasswordHash []byte
}

func (p *password) Set(plainTextPassword string, ve *service.OperationErrors) {
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

func (p *password) Matches(plaintextPassword string) (bool, error) {
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

// Insert creates a new user in the database
func (r *UserRepository) Insert(user *UserModel) (*UserModel, error) {
	query := `INSERT INTO users (name, email, password, activated, version, created_at) 
	VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.Activated, user.Version, user.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Retrieve the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

// GetByEmail retrieves a user by email from the database
func (r *UserRepository) GetByEmail(email string) (*UserModel, error) {
	query := `SELECT id, name, email, password, activated, version, created_at 
		FROM users WHERE email = ?`
	row := r.DB.QueryRow(query, email)

	var user UserModel
	err := row.Scan(&user.ID, &user.Name.Inner_value, &user.Email.Inner_value, &user.Password.PasswordHash, &user.Activated, &user.Version, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user's details in the database
func (r *UserRepository) Update(user *UserModel) error {
	query := `UPDATE users SET name = ?, email = ?, password = ?, activated = ?,  version = version + 1  
		WHERE id = ?`
	_, err := r.DB.Exec(query, user.Name.Inner_value, user.Email.Inner_value, user.Password.PasswordHash, user.Activated, user.ID)
	return err
}
