package model

import (
	"errors"

	"todo-service/common"
)

const (
	EntityName = "User"
)

type User struct {
	common.SQLModel
	Email     string   `json:"email" gorm:"column:email;"`
	Password  string   `json:"password" gorm:"column:password;"`
	Salt      string   `json:"salt" gorm:"column:salt;"`
	LastName  string   `json:"last_name" gorm:"column:last_name;"`
	FirstName string   `json:"first_name" gorm:"column:first_name;"`
	Phone     string   `json:"phone" gorm:"column:phone;"`
	Status    int      `json:"status" gorm:"column:status;"`
	Role      UserRole `json:"role" gorm:"column:role;"`
}

func (u *User) GetUserId() int {
	return u.Id
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRole() string {
	return u.Role.String()
}

func (User) TableName() string {
	return "users"
}

type UserCreation struct {
	common.SQLModel
	Email     string `json:"email" gorm:"column:email;"`
	Password  string `json:"password" gorm:"column:password;"`
	LastName  string `json:"last_name" gorm:"column:last_name;"`
	FirstName string `json:"first_name" gorm:"column:first_name;"`
	Phone     string `json:"phone" gorm:"column:phone;"`
	Role      string `json:"-" gorm:"column:role;"`
	Salt      string `json:"-" gorm:"column:salt;"`
}

func (UserCreation) TableName() string {
	return User{}.TableName()
}

type UserLogin struct {
	Email    string `json:"email" gorm:"column:email;"`
	Password string `json:"password" gorm:"column:password;"`
}

func (UserLogin) TableName() string {
	return User{}.TableName()
}

var (
	ErrEmailOrPasswordInvalid = common.NewCustomError(
		errors.New("email or password invalid"),
		"email or password invalid",
		"ErrEmailOrPasswordInvalid",
	)

	ErrEmailExisted = common.NewCustomError(
		errors.New("email has already existed"),
		"email has already existed",
		"ErrEmailExisted",
	)
)
