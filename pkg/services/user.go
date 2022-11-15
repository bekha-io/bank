package services

import (
	"banking/pkg/errors"
	"banking/pkg/models"
	"banking/pkg/types"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

type userService interface {
	CreateUser(login string, password string, phoneNumber types.PhoneNumber, firstName string, lastName string,
		middleName string) (user *models.User, err error)
	GetUserById(userID string) (user *models.User, err error)
	GetUserByPhoneNumber(phoneNumber types.PhoneNumber) (user *models.User, err error)
	GetUserByLogin(login string) (user *models.User, err error)
	CountUsersBy(fieldName string, value interface{}) (count int64, err error)
	CountUsersByLogin(login string) (count int64, err error)
	CountUsersByPhoneNumber(phoneNumber types.PhoneNumber) (count int64, err error)
}

func validateUserMandatoryFields(login, password string) error {
	// Проверка на наличие всех обязательных (!) полей
	if login == "" || password == "" {
		return fmt.Errorf("login, password, firstName, lastName are required fields and cannot be "+
			"empty. Values provided: %v, %v", login, password)
	}
	return nil
}

func generatePasswordHash(rawPassword string) ([]byte, error) {
	// Генерация хэша из сырого пароля
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("cannot generate hash from password. Error: ", passwordHash)
		return nil, err
	}
	return passwordHash, err
}

func (s *ServiceManager) CreateUserForLegalEntity(login, password string, phoneNumber types.PhoneNumber,
	officialName string, taxIdNumber string) (user *models.User, err error) {

	// Проверка на наличие обязательных аутентификационных (!) полей
	if err = validateUserMandatoryFields(login, password); err != nil {
		return nil, err
	}

	// Number should be in international format (without '+' sign)
	if err = phoneNumber.Validate(); err != nil {
		return nil, errors.InvalidPhoneNumberError
	}

	// Генерация хэша пароля
	passwordHash, err := generatePasswordHash(password)
	if err != nil {
		return nil, err
	}

	// Проверяем, существует ли юзер с таким логином
	usersWithLogin, err := s.CountUsersByLogin(login)
	if err != nil || usersWithLogin > 0 {
		log.Println("User with the given login already exists. Login: ", login)
		log.Println("Error: ", err)
		return nil, errors.LoginOccupiedError
	}

	// + доп. проверка на то, занят ли указанный номер телефона
	usersWithPhN, err := s.CountUsersByPhoneNumber(phoneNumber)
	if err != nil || usersWithPhN > 0 {
		log.Println("User with the given phone number already exists. Phone number: ", phoneNumber)
		log.Println("Error: ", err)
		return nil, errors.PhoneNumberOccupiedError
	}

	var customerInfo = make(map[string]interface{})
	customerInfo["official_name"] = officialName
	customerInfo["tax_id_number"] = taxIdNumber

	user = &models.User{
		Login:        login,
		Password:     passwordHash,
		PhoneNumber:  phoneNumber,
		CustomerInfo: customerInfo,
		LegalStatus:  types.LegalEntityStatus,
	}

	result := s.DB.Create(user)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return user, nil
}

func (s *ServiceManager) CreateUser(login string, password string,
	phoneNumber types.PhoneNumber, firstName string, lastName string,
	middleName string) (user *models.User, err error) {

	// Проверка на наличие обязательных аутентификационных (!) полей
	if err = validateUserMandatoryFields(login, password); err != nil {
		return nil, err
	}

	// Number should be in international format (without '+' sign)
	if err = phoneNumber.Validate(); err != nil {
		return nil, errors.InvalidPhoneNumberError
	}

	// Генерация хэша пароля
	passwordHash, err := generatePasswordHash(password)
	if err != nil {
		return nil, err
	}

	// Проверяем, существует ли юзер с таким логином
	usersWithLogin, err := s.CountUsersByLogin(login)
	if err != nil || usersWithLogin > 0 {
		log.Println("User with the given login already exists. Login: ", login)
		log.Println("Error: ", err)
		return nil, errors.LoginOccupiedError
	}

	// + доп. проверка на то, занят ли указанный номер телефона
	usersWithPhN, err := s.CountUsersByPhoneNumber(phoneNumber)
	if err != nil || usersWithPhN > 0 {
		log.Println("User with the given phone number already exists. Phone number: ", phoneNumber)
		log.Println("Error: ", err)
		return nil, errors.PhoneNumberOccupiedError
	}

	var customerInfo = make(map[string]interface{})
	customerInfo["first_name"] = firstName
	customerInfo["last_name"] = lastName
	customerInfo["middle_name"] = middleName

	// Создание пользователя и сохранение записи в БД
	user = &models.User{
		Login:        login,
		Password:     passwordHash,
		PhoneNumber:  phoneNumber,
		CustomerInfo: customerInfo,
		LegalStatus:  types.IndividualStatus,
	}

	result := s.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s *ServiceManager) GetUserById(userID string) (user *models.User, err error) {
	if err = s.DB.First(&user, "id = ?", userID).Error; err != nil {

		if err == gorm.ErrRecordNotFound || user == nil {
			return nil, errors.UserDoesNotExist
		}

		log.Println("Cannot get user by the given id. Error: ", err)
		return nil, err
	}
	return user, nil
}

func (s *ServiceManager) GetUserByLogin(login string) (user *models.User, err error) {
	if err = s.DB.First(&user, "login = ?", login).Error; err != nil {

		if err == gorm.ErrRecordNotFound || user == nil {
			return nil, errors.UserDoesNotExist
		}

		log.Println("Cannot get user by the given login. Error: ", err)
		return nil, err
	}
	return user, nil
}

func (s *ServiceManager) GetUserByPhoneNumber(phoneNumber types.PhoneNumber) (user *models.User, err error) {
	if err = s.DB.First(&user, "phone_number = ?", phoneNumber).Error; err != nil {

		if err == gorm.ErrRecordNotFound || user == nil {
			return nil, errors.UserDoesNotExist
		}

		log.Println("Cannot get user by the given phone number. Error: ", err)
		return nil, err
	}
	return user, nil
}

func (s *ServiceManager) CountUsersBy(fieldName string, value interface{}) (count int64, err error) {
	if res := s.DB.Table("users").Where(fieldName+" = ?", value).Count(&count); res.Error != nil {
		return count, res.Error
	}
	return count, err
}

func (s *ServiceManager) CountUsersByLogin(login string) (count int64, err error) {
	count, err = s.CountUsersBy("login", login)
	return count, err
}

func (s *ServiceManager) CountUsersByPhoneNumber(phoneNumber types.PhoneNumber) (count int64, err error) {
	count, err = s.CountUsersBy("phone_number", string(phoneNumber))
	return count, err
}

//TODO
/*
1.1) Roles
	- admin
	- support


Postgres
pgadmin4 --- data grip
postman



1. Create User (delete, get, update)
	- add card (валидационные проверки)
	-
[X] 2. Login(Authorization) (login, password) bcrypt (access token - 15min)
[X]	- password hash

[] 3. Middleware
[] JWT(access token 3 min, refresh token 15 - 20)



4. files
5. config
6. Metanit


7. PreCheck, PostCheck, Payment

*/
