package goutils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractGinToken(c *gin.Context, headerKey string, tokenType string) (string, error) {
	token := c.GetHeader(headerKey)
	if token == "" {
		return "", errors.New("token is missing")
	}

	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != tokenType {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}
