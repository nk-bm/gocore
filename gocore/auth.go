package gocore

import (
	"github.com/gin-gonic/gin"
	"github.com/nk-bm/gocore/gincore/ginmw"
	"github.com/nk-bm/gocore/goutils"
)

type AuthManager struct {
	secretKey string
	idKey     string
	authType  string
}

func NewAuthManager(secretKey string, idKey, authType string) *AuthManager {
	return &AuthManager{
		secretKey: secretKey,
		idKey:     idKey,
		authType:  authType,
	}
}

func (m *AuthManager) GinMiddleware() gin.HandlerFunc {
	return ginmw.AuthMW(m.secretKey, m.idKey, m.authType)
}

func (m *AuthManager) ExtractIDFromToken(token string) (int64, error) {
	return goutils.ExtractIDFromJWT(m.secretKey, token, m.idKey, m.authType)
}

func (m *AuthManager) GenerateToken(id int64) (string, error) {
	return goutils.GenerateJWT(m.secretKey, m.authType, m.idKey, id)
}
