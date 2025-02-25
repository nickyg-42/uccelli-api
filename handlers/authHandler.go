package handlers

import (
	"encoding/json"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	if err := utils.ValidateNewUser(r, userDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error validating new user:", err)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error generating password hash:", err)
		return
	}

	validEmails := []string{
		"evangelinegrimaldi8@gmail.com",
		"nicholasgrimaldi42@gmail.com",
		"fiona.tetreault@gmail.com",
		"jessicaagrimaldi@gmail.com",
		"noahgrimaldi1@gmail.com",
		"josiahgrimaldi@gmail.com",
		"john.grimaldi@gmail.com",
		"dominicgrimaldi1738@gmail.com",
		"karengrimaldi@gmail.com",
	}

	if !slices.Contains(validEmails, strings.ToLower(userDto.Email)) {
		http.Error(w, "Email not whitelisted", http.StatusInternalServerError)
		log.Println("Email not whitelisted")
		return
	}

	user := models.User{
		FirstName:    userDto.FirstName,
		LastName:     userDto.LastName,
		Email:        userDto.Email,
		Username:     userDto.Username,
		PasswordHash: hashedPassword,
	}

	createdUser, err := db.CreateUser(r.Context(), &user)
	if err != nil {
		log.Println("Error saving user:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error creating user:", err)
		return
	}

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	user, err := db.GetUserByUsername(r.Context(), credentials.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		log.Println("Error retrieving user by username:", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		log.Println("Error comparing password hash:", err)
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"user_id":  user.ID,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		log.Println("Error generating JWT token:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
