package utils

import (
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

// func IsSelfOrSA(r http.Request) bool {
// 	role := r.Context().Value("role").(string)
// 	authenticatedUserID := r.Context().Value("user_id")

// 	// If not self or SA, deny
// 	if id != authenticatedUserID && role != string(models.SuperAdmin) {
// 		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
// 		return
// 	}
// }
