package service

import (
	"authentication-service/internal/data"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type TokenService struct {
	RepoManager *data.RepoManager
}

func NewTokenService(repoManager *data.RepoManager) *TokenService {
	return &TokenService{RepoManager: repoManager}
}

func (s *TokenService) CreateAccessToken(userID int64, scope string, ttl time.Duration, secret data.TokenScope) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"scope": scope,
		"exp":   time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(secret)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *TokenService) ValidateToken(tokenString string, secret string) (bool, error) {
	secretKey := []byte(secret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["exp"].(int64) < time.Now().Unix() {
			return false, errors.New("token has expired")
		}
		return true, nil
	}

	return false, errors.New("invalid token")
}

func (s *TokenService) GetTokensForUser(userID int64) ([]data.Token, error) {

	return s.RepoManager.TokenRepo.GetByUserID(userID)
}

func (s *TokenService) GetTokensForUserAndScope(userID int64, scope string) ([]data.Token, error) {

	return s.RepoManager.TokenRepo.GetByUserIDAndScope(userID, scope)
}
func (s *TokenService) DeleteToken(tokenHash []byte) error {
	return s.RepoManager.TokenRepo.Delete(tokenHash)
}
