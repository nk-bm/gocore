package goutils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(secretKey, authType, idKey string, id int64) (string, error) {
	return GenerateJWTWithExpiration(secretKey, authType, idKey, id, time.Hour*24)
}

func GenerateJWTWithExpiration(secretKey, authType, idKey string, id int64, expDuration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims[idKey] = id
	claims["auth_type"] = authType

	claims["exp"] = time.Now().Add(expDuration).Unix()
	claims["iat"] = time.Now().Unix()
	// sign token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractIDFromJWT(secretKey, jwtToken, idKey, authType string) (int64, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token claims")
	}

	expTime, err := claims.GetExpirationTime()
	if err != nil || expTime == nil {
		return 0, fmt.Errorf("bad token format")
	}

	if time.Now().After(expTime.Time) {
		return 0, fmt.Errorf("token expired")
	}

	if at, ok := claims["auth_type"].(string); !ok || at != authType {
		return 0, fmt.Errorf("auth type not valid")
	}

	id, ok := claims[idKey].(float64)
	if !ok {
		return 0, fmt.Errorf("user id not valid")
	}

	return int64(id), nil
}

type JWTManager struct {
	secretKey string
	idKey     string
	authType  string
}

func NewJWTManager(secretKey string, idKey string, authType string) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
		idKey:     idKey,
		authType:  authType,
	}
}

func (m *JWTManager) GenerateJWT(id int64) (string, error) {
	return GenerateJWTWithExpiration(m.secretKey, m.authType, m.idKey, id, time.Hour*24)
}

func (m *JWTManager) Validate(token string) bool {
	_, err := ExtractIDFromJWT(m.secretKey, token, m.idKey, m.authType)
	return err == nil
}

func (m *JWTManager) ExtractID(jwtToken string) (int64, error) {
	return ExtractIDFromJWT(m.secretKey, jwtToken, m.idKey, m.authType)
}
