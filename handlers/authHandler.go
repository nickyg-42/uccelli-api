package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"nest/db"
	"nest/models"
	"nest/utils"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var userDto models.UserDTO

	if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
		log.Printf("ERROR: Failed to decode registration request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.ValidateNewUser(r, userDto); err != nil {
		log.Printf("ERROR: User validation failed for %s: %v", userDto.Email, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR: Failed to hash password for user %s: %v", userDto.Email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	validEmailsStr := os.Getenv("VALID_EMAILS")
	validEmails := strings.Split(validEmailsStr, ",")

	if !slices.Contains(validEmails, strings.ToLower(userDto.Email)) {
		log.Printf("ERROR: Registration attempt with non-whitelisted email: %s", userDto.Email)
		http.Error(w, "Email not whitelisted", http.StatusUnauthorized)
		return
	}

	user := models.User{
		FirstName:    strings.ToLower(userDto.FirstName),
		LastName:     strings.ToLower(userDto.LastName),
		Email:        strings.ToLower(userDto.Email),
		Username:     strings.ToLower(userDto.Username),
		PasswordHash: hashedPassword,
	}

	createdUser, err := db.CreateUser(r.Context(), &user)
	if err != nil {
		log.Printf("ERROR: Failed to create user in database - Email: %s, Username: %s: %v",
			userDto.Email, userDto.Username, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Successfully registered new user - ID: %d, Email: %s, Username: %s",
		createdUser.ID, createdUser.Email, createdUser.Username)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func Login(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("ERROR: Failed to decode login request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByUsername(r.Context(), strings.ToLower(credentials.Username))
	if err != nil {
		log.Printf("ERROR: Failed to find user during login - Username: %s: %v", credentials.Username, err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(credentials.Password)); err != nil {
		log.Printf("ERROR: Invalid password attempt for user %s from IP %s",
			credentials.Username, r.RemoteAddr)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 504).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("ERROR: Failed to generate JWT token for user %s: %v",
			user.Username, err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Successful login - User: %s, ID: %d", user.Username, user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func GeneratePasswordResetCode(w http.ResponseWriter, r *http.Request) {
	var email struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		log.Printf("ERROR: Failed to decode email for password reset: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resetCode, err := db.GeneratePasswordResetCode(r.Context(), email.Email)
	if err != nil {
		log.Printf("ERROR: Failed to generate password reset code for email %s: %v", email.Email, err)
		http.Error(w, "Failed to generate password reset code", http.StatusInternalServerError)
		return
	}

	emailBody := fmt.Sprintf(`Here is your password reset code, if you did not request this please contact an Admin: %s`, resetCode)
	utils.NotifyUser(email.Email, "Password Reset Code", emailBody)

	log.Printf("INFO: Password reset code generated for email %s, %s", email.Email, resetCode)
	w.WriteHeader(http.StatusOK)
}

func VerifyPasswordResetCode(w http.ResponseWriter, r *http.Request) {
	var code struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&code); err != nil {
		log.Printf("ERROR: Failed to decode password reset code: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.VerifyPasswordResetCode(r.Context(), code.Code, code.Email)
	if err != nil {
		log.Printf("ERROR: Failed to verify password reset code %s: %v", code.Code, err)
		http.Error(w, "Failed to verify password reset code", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Password reset code %s successfully verified", code.Code)
	w.WriteHeader(http.StatusOK)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var reset struct {
		Email    string `json:"email"`
		Code     string `json:"code"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reset); err != nil {
		log.Printf("ERROR: Failed to decode password reset request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash new password and commence with update
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reset.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR: Failed to hash password for user %s: %v", reset.Email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = db.ResetPassword(r.Context(), reset.Email, reset.Code, hashedPassword)
	if err != nil {
		log.Printf("ERROR: Failed to reset password for code %s: %v", reset.Code, err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Password successfully reset for code %s", reset.Code)
	w.WriteHeader(http.StatusOK)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
