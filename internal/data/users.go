package data

import (
	"database/sql"
	"fmt"
)

// UserRepository struct holds a reference to the database connection
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository with the given DB connection
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Insert creates a new user in the database
func (r *UserRepository) Insert(user *UserModel) (*UserModel, error) {
	query := `INSERT INTO users (name, email, password_hash, activated, version, created_at) 
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
	query := `SELECT id, name, email, password_hash, activated, version, created_at 
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
func (r *UserRepository) GetById(userID int64) (*UserModel, error) {
	query := `SELECT id, name, email, password_hash, activated, version, created_at 
		FROM users WHERE id = ?`
	row := r.DB.QueryRow(query, userID)

	var user UserModel
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Activated, &user.Version, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", userID)
		}
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user's details in the database
func (r *UserRepository) Update(user *UserModel) error {

	query := `UPDATE users SET name = ?, email = ?, password_hash = ?,  version = version + 1  , activated=?
		WHERE id = ?`
	_, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.ID, user.Activated)
	return err
}
