package data

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func GenerateHashToken() ([]byte, error) {
	token, err := GenerateRandomToken()
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash token: %w", err)
	}
	return hash, nil
}

func GenerateRandomToken() (string, error) {
	tokenBytes := make([]byte, 16)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate random bytes: %w", err)
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token, nil
}

type TokenRepository struct {
	DB *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{DB: db}
}

// Insert inserts a new token into the database
func (r *TokenRepository) Insert(token *Token) (*Token, error) {
	// Prepare the SQL query to insert a new token
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) 
	VALUES (?, ?, ?, ?)`

	// Execute the query
	_, err := r.DB.Exec(query, token.Hash, token.UserID, token.Expiry, string(token.Scope))
	if err != nil {
		return nil, fmt.Errorf("could not insert token: %w", err)
	}

	// Return the token object after insertion
	return token, nil
}

func (r *TokenRepository) Delete(hash []byte) error {
	query := `DELETE FROM tokens WHERE hash = ?`

	_, err := r.DB.Exec(query, hash)
	if err != nil {
		return fmt.Errorf("could not delete token: %w", err)
	}

	return nil
}
func (r *TokenRepository) DeleteTokensForUser(userID int64, scope TokenScope) error {
	query := `DELETE FROM tokens WHERE user_id = ? and scope = ?`

	_, err := r.DB.Exec(query, userID, string(scope))
	if err != nil {
		return fmt.Errorf("could not delete token: %w", err)
	}

	return nil
}

func (repo *TokenRepository) GetByUserID(userID int64) ([]Token, error) {
	query := `SELECT hash, user_id, expiry, scope FROM tokens WHERE user_id = ?`

	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying tokens: %w", err)
	}
	defer rows.Close()

	var tokens []Token

	for rows.Next() {
		var token Token
		if err := rows.Scan(&token.Hash, &token.UserID, &token.Expiry, &token.Scope); err != nil {
			return nil, fmt.Errorf("error scanning token: %w", err)
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tokens, nil
}

func (repo *TokenRepository) GetByUserIDAndScope(userID int64, scope TokenScope) ([]Token, error) {
	query := `SELECT hash, user_id, expiry, scope FROM tokens WHERE user_id = ? and scope = ?`

	rows, err := repo.DB.Query(query, userID, string(scope))
	if err != nil {
		return nil, fmt.Errorf("error querying tokens: %w", err)
	}
	defer rows.Close()

	var tokens []Token
	var expiry string

	for rows.Next() {
		var token Token
		if err := rows.Scan(&token.Hash, &token.UserID, &expiry, &token.Scope); err != nil {
			return nil, fmt.Errorf("error scanning token: %w", err)
		}
		parsedExpiry, err := time.Parse("2006-01-02 15:04:05.999999-07:00", expiry)
		if err != nil {
			return nil, fmt.Errorf("error parsing expiry: %w", err)
		}
		token.Expiry = parsedExpiry
		if !isValidTokenScope(string(token.Scope)) {
			return nil, fmt.Errorf("invalid token scope: %s", token.Scope)
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tokens, nil
}

func isValidTokenScope(scope string) bool {
	switch TokenScope(scope) {
	case ActivateEmailToken, UserAccessToken:
		return true
	default:
		return false
	}
}
