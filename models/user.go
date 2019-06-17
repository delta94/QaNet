package models

import (
	"time"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/jinzhu/gorm"
)

// Roles - different user roles
var Roles = map[string]string{
	"admin_user":  "ADMIN_USER",
	"normal_user": "NORMAL_USER",
}

// User is the gorm model representation of table users.
type User struct {
	ID           string    `json:"id"`
	DisplayName  string    `json:"display_name"`
	About        string    `json:"about"`
	Points       int       `json:"points"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	RoleID       string    `json:"role_id"`
	IsActive     bool      `json:"is_active"`
	LastLoggedIn time.Time `json:"last_logged_in"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt" gorm:"not null"`

	Password string `json:"password" gorm:"-"`
}

// Validate - validates the user details.
func (u *User) Validate(tx *gorm.DB) *validate.Errors {
	verrs := validate.Validate(
		&validators.StringIsPresent{Field: u.ID, Name: "Id"},
		&validators.StringLengthInRange{Field: u.ID, Name: "Id", Min: 6},

		&validators.StringIsPresent{Field: u.DisplayName, Name: "DisplayName"},
		&validators.StringLengthInRange{Field: u.DisplayName, Name: "DisplayName", Min: 6},

		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
	)

	var existingUser User

	if len(verrs.Get("id")) == 0 || len(verrs.Get("id")) == 0 {
		tx.Table("users").Where("id = ? or email=?", u.ID, u.Email).First(&existingUser)

		// check to see if the ID is already taken:
		if len(verrs.Get("id")) == 0 {
			if existingUser.ID == u.ID {
				verrs.Add("id", "UserID already taken.")
			}
		}

		// check to see if the email address is already taken:
		if len(verrs.Get("email")) == 0 {
			if existingUser.Email == u.Email {
				verrs.Add("email", "Email already taken.")
			}
		}
	}

	return verrs
}
