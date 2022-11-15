package models

import (
	"banking/pkg/types"
	"banking/pkg/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"unicode/utf8"
)

type User struct {
	BaseModel
	ID         string `gorm:"primaryKey" json:"id"`
	TelegramID string `gorm:"unique" json:"telegram_id"`

	Login       string            `gorm:"notNull" json:"login"`
	Password    []byte            `gorm:"notNull" json:"-"`
	PhoneNumber types.PhoneNumber `gorm:"unique;" json:"phone_number"`

	CustomerInfo map[string]interface{} `json:"customer_info" gorm:"serializer:json"`

	LegalStatus types.LegalStatus `json:"legal_status" gorm:"notNull"`

	// API access params
	Role        types.Role `gorm:"default:0" json:"-"`
	AccessToken string     `json:"access_token,omitempty"`
}

type IndividualCustomerRaw struct {
	// Individual status related attributes (omitempty all)
	Name       string `json:"first_name,omitempty" binding:"required"`
	LastName   string `json:"last_name,omitempty" binding:"required"`
	MiddleName string `json:"middle_name,omitempty"`
}

type LegalEntityCustomerRaw struct {
	// Legal entity status related attributes (omitempty all)
	OfficialName string `json:"official_name,omitempty" binding:"required"`
	TaxIdNumber  string `json:"tax_id_number,omitempty" binding:"required"`
}

func (u *User) Mask() {
	u.AccessToken = ""
}

func generateAccessToken(login string) (accessToken string) {
	accessToken, _ = utils.GenerateAccessToken(login)
	if accessToken != "" {
		return accessToken
	}
	return ""
}

func (u *User) isValidOnCreate() (err error) {
	if u.IsIndividualStatus() {
		name, ok := u.CustomerInfo["first_name"]
		if !ok || utf8.RuneCountInString(name.(string)) == 0 {
			return errors.New("user should have first name")
		}
		lastName, ok := u.CustomerInfo["last_name"]
		if !ok || utf8.RuneCountInString(lastName.(string)) == 0 {
			return errors.New("user should have last name")
		}

	} else if u.IsLegalEntityStatus() {
		officialName, ok := u.CustomerInfo["official_name"]
		if !ok || utf8.RuneCountInString(officialName.(string)) == 0 {
			return errors.New("legal entity user should have a full official name")
		}

		taxIdNum, ok := u.CustomerInfo["tax_id_number"]
		if !ok || utf8.RuneCountInString(taxIdNum.(string)) == 0 {
			return errors.New("legal entity user should have a valid tax id number")
		}

		//var asModel LegalEntityCustomerRaw
		//if err = json.Unmarshal(u.CustomerInfo, &asModel); err == nil {
		//	if asModel.OfficialName == "" || asModel.TaxIdNumber == "" {
		//		return err
		//	}
		//}
	}

	return err
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

func (u *User) IsLegalEntityStatus() bool {
	return u.LegalStatus == types.LegalEntityStatus
}

func (u *User) IsIndividualStatus() bool {
	return u.LegalStatus == types.IndividualStatus
}
