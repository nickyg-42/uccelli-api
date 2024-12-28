package utils

import (
	"nest/db"
	"nest/models"
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

func IsSelfOrSA(r *http.Request, userID int) bool {
	role := r.Context().Value("role").(string)
	authenticatedUserID := r.Context().Value("user_id").(int)

	// If not self or SA, deny
	if userID != authenticatedUserID && role != string(models.SuperAdmin) {
		return false
	}

	return true
}

func IsGroupOwnerOrSA(r *http.Request, groupID int) bool {
	role := r.Context().Value("role").(string)
	authenticatedUserID := r.Context().Value("user_id").(int)
	isGroupAdmin, err := db.IsUserGroupAdmin(r.Context(), authenticatedUserID, groupID)
	if err != nil {
		return false
	}

	// If not group owner or SA, deny
	if role != string(models.SuperAdmin) && !isGroupAdmin {
		return false
	}

	return true
}
