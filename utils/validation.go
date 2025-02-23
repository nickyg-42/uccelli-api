package utils

import (
	"errors"
	"log"
	"nest/db"
	"nest/models"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

// **************************************
// USER VALIDATION
// **************************************
func ValidateNewUser(r *http.Request, userDTO models.UserDTO) error {
	if !ValidateEmail(userDTO.Email) {
		return errors.New("invalid email")
	}

	if !ValidateName(userDTO.FirstName) {
		return errors.New("invalid first name")
	}

	if !ValidateName(userDTO.LastName) {
		return errors.New("invalid last name")
	}

	if !ValidateUsername(r, userDTO.Username) {
		return errors.New("invalid username")
	}

	if !ValidatePassword(userDTO.Password) {
		return errors.New("invalid password")
	}

	return nil
}

func ValidateEmail(email string) bool {
	if len(email) < 1 {
		return false
	}

	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	emailRegex := regexp.MustCompile(emailRegexPattern)

	return emailRegex.MatchString(email)
}

// Alphabetic, < 50, no special chars (besides dash/space)
func ValidateName(name string) bool {
	log.Println("name received: ", name)

	if len(name) < 1 {
		return false
	}

	log.Println("passed len check", len(name))

	const nameRegexPattern = `^[a-zA-Z]+(?:[ '-][a-zA-Z]+)*$`

	nameRegex := regexp.MustCompile(nameRegexPattern)

	log.Println("nameRegex: ", nameRegex)
	log.Println("Passes regex?: ", nameRegex.MatchString(name))

	return len(name) > 0 && len(name) <= 50 && nameRegex.MatchString(name)
}

// Alphanumeric, underscores/dots
func ValidateUsername(r *http.Request, username string) bool {
	if len(username) < 1 {
		return false
	}

	const usernameRegexPattern = `^[a-zA-Z][a-zA-Z0-9._]{2,29}$`

	usernameRegex := regexp.MustCompile(usernameRegexPattern)

	if !usernameRegex.MatchString(username) {
		return false
	}

	isTaken, err := db.IsUsernameTaken(r.Context(), username)
	if err != nil {
		log.Println("Error checking username", err)
		return false
	}

	return !isTaken
}

// At least one uppercase letter, one lowercase letter, one number, one special character, at least 8 characters long
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// **************************************
// EVENT VALIDATION
// **************************************
func ValidateNewEvent(event models.EventDTO) error {
	if len(strings.TrimSpace(event.Name)) == 0 {
		return errors.New("event name cannot be empty")
	}
	if len(event.Name) > 255 {
		return errors.New("event name cannot exceed 255 characters")
	}

	if len(event.Description) > 1000 {
		return errors.New("event description cannot exceed 1000 characters")
	}

	if event.StartTime.IsZero() || event.EndTime.IsZero() {
		return errors.New("start time and end time must be provided")
	}
	if event.StartTime.After(event.EndTime) {
		return errors.New("start time cannot be after end time")
	}

	return nil
}

// **************************************
// GROUP VALIDATION
// **************************************
func ValidateNewGroup(group models.GroupDTO) error {
	if len(strings.TrimSpace(group.Name)) == 0 {
		return errors.New("group name cannot be empty")
	}
	if len(group.Name) > 255 {
		return errors.New("group name cannot exceed 255 characters")
	}

	return nil
}
