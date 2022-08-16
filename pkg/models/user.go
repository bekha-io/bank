package models

import (
	"banking/pkg/types"
	"banking/pkg/utils"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Login       string `gorm:"notNull"`
	Password    []byte `gorm:"notNull"`
	Name        string `gorm:"notNull"`
	LastName    string `gorm:"notNull"`
	MiddleMame  string
	PhoneNumber types.PhoneNumber `gorm:"unique;"`

	AccessToken string `json:"access_token"`
}

func generateAccessToken(login string) (accessToken string) {
	accessToken, _ = utils.GenerateAccessToken(login)
	if accessToken != "" {
		return accessToken
	}
	return ""
}

func (u *User) isValidOnCreate() (err error) {
	err = u.PhoneNumber.IsValid()
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	err = u.isValidOnCreate()
	if err != nil {
		return
	}

	u.ID = uuid.New().String()
	u.AccessToken = generateAccessToken(u.Login)
	return
}

func (u *User) FullName() string {
	return u.LastName + " " + u.Name + " " + u.MiddleMame
}

func (u *User) GetAccounts() []Account {
	var accounts []Account
	Db.Find(&accounts, "user_id = ?", u.ID)
	return accounts
}

func (u *User) CheckPassword(passwordInput []byte) (ok bool) {
	err := bcrypt.CompareHashAndPassword(u.Password, passwordInput)
	if err != nil {
		fmt.Println("Error occured while comparing hashed and unhashed passwords. Error: ", err)
		return false
	}
	return true
}
