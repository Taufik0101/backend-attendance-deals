package models

import (
	"backend-attendance-deals/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseModel
	Username string         `json:"username" gorm:"column:username;type:varchar(255);not null"`
	Password string         `json:"-" gorm:"column:password;type:text;not null"`
	Salary   int            `json:"salary" gorm:"column:salary;type:integer;not null"`
	Role     utils.UserType `json:"role" gorm:"column:role;type:role_types;not null;default: 'employee'"`
}

func (*User) TableName() string {
	return "users"
}

// ComparePasswords Compare user password and payload
func (u *User) ComparePasswords(password string) error {
	if u.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
			return err
		}
	}

	return nil
}
