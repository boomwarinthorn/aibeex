package entity

import "time"

// User represents the core user entity in the domain
type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password,omitempty"`
	FullName    string    `json:"fullName"`
	PhoneNumber string    `json:"phoneNumber"`
	Birthday    string    `json:"birthday"`
	CreatedAt   time.Time `json:"createdAt"`
}

// NewUser creates a new user entity
func NewUser(email, password, fullName, phoneNumber, birthday string) *User {
	return &User{
		Email:       email,
		Password:    password,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Birthday:    birthday,
		CreatedAt:   time.Now(),
	}
}

// IsValidEmail checks if the email format is valid
func (u *User) IsValidEmail() bool {
	return u.Email != "" && len(u.Email) > 5
}

// WithoutPassword returns user without password field for security
func (u *User) WithoutPassword() *User {
	userCopy := *u
	userCopy.Password = ""
	return &userCopy
}
