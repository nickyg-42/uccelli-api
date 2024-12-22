package utils

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your_secret_key")

func GenerateToken(userID int) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func HasAccess(r *http.Request, requiredRole string) bool {
	role := r.Context().Value("role")
	if role == nil {
		return false
	}

	roleHierarchy := map[string]int{
		"sa":    3,
		"admin": 2,
		"user":  1,
	}

	userRoleLevel, userRoleExists := roleHierarchy[role.(string)]
	requiredRoleLevel, requiredRoleExists := roleHierarchy[requiredRole]

	return userRoleExists && requiredRoleExists && userRoleLevel >= requiredRoleLevel
}
