package data

import (
	"authentication-service/internal/domain"
	"database/sql"
	"fmt"
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

type UserModel struct {
	ID        int64
	CreatedAt time.Time
	Name      string
	Email     string
	Password  []byte
	Activated bool
	Version   int
}

func CreateUserModelFromUserDomain(u *domain.User) *UserModel {
	return &UserModel{
		Name:      u.Name.Inner_value,
		CreatedAt: u.CreatedAt,
		Email:     u.Email.Inner_value,
		Password:  u.Password.PasswordHash,
		Activated: u.Activated,
		Version:   u.Version,
	}
}

// Insert creates a new user in the database
func (r *UserRepository) Insert(user *UserModel) error {
	query := `INSERT INTO users (name, email, password, activated, version, created_at) 
	VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.Activated, user.Version, user.CreatedAt)
	return err
}

// GetByEmail retrieves a user by email from the database
func (r *UserRepository) GetByEmail(email string) (*UserModel, error) {
	query := `SELECT id, name, email, password, activated, version, created_at 
		FROM users WHERE email = ?`
	row := r.DB.QueryRow(query, email)

	var user UserModel
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Activated, &user.Version, &user.CreatedAt)
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
	_, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.Activated, user.ID)
	return err
}
