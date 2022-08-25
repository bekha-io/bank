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
	BaseModel
	ID          string            `gorm:"primaryKey" json:"id"`
	Login       string            `gorm:"notNull" json:"login"`
	Password    []byte            `gorm:"notNull" json:"-"`
	Name        string            `gorm:"notNull" json:"name"`
	LastName    string            `gorm:"notNull" json:"lastName"`
	MiddleName  string            `json:"middleName"`
	PhoneNumber types.PhoneNumber `gorm:"unique;" json:"phoneNumber"`
	Role        types.Role        `gorm:"default:0" json:"-"`
	AccessToken string            `json:"-"`
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
	return u.LastName + " " + u.Name + " " + u.MiddleName
}

func (u *User) GetAccounts() []Account {
	var accounts []Account
	Db.Find(&accounts, "user_id = ?", u.ID)
	return accounts
}

func (u *User) GetCards() []Card {
	var accounts = u.GetAccounts()
	var cards []Card

	for _, acc := range accounts {
		for _, card := range acc.GetCards() {
			cards = append(cards, card)
		}
	}

	return cards
}

func (u *User) CheckPassword(passwordInput []byte) (ok bool) {
	err := bcrypt.CompareHashAndPassword(u.Password, passwordInput)
	if err != nil {
		fmt.Println("Error occured while comparing hashed and unhashed passwords. Error: ", err)
		return false
	}
	return true
}

func (u *User) IsAdminRole() bool {
	return u.Role == types.AdminRole
}

func (u *User) IsSupportRole() bool {
	return u.Role == types.SupportRole
}

func (u *User) IsUserRole() bool {
	return u.Role == types.UserRole
}
