package services

import (
	"banking/pkg/errors"
	"banking/pkg/models"
	"banking/pkg/types"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type userService interface {
	CreateUser(login string, password string,
		phoneNumber types.PhoneNumber, firstName string, lastName string,
		middleName string) (user *models.User, err error)
	GetUserByPhoneNumber(phoneNumber types.PhoneNumber) (user *models.User, err error)
	GetUserByLogin(login string) (user *models.User, err error)
}

func (s *ServiceManager) CreateUser(login string, password string,
	phoneNumber types.PhoneNumber, firstName string, lastName string,
	middleName string) (user *models.User, err error) {

	// Проверка на наличие всех обязательных (!) полей
	if login == "" || password == "" || lastName == "" || firstName == "" {
		return nil, fmt.Errorf("login, password, firstName, lastName are required fields and cannot be "+
			"empty. Values provided: %v, %v, %v, %v", login, password, firstName, lastName)
	}

	// Генерация хэша из сырого пароля
	passwordBytes := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		log.Println("cannot generate hash from password. Error: ", passwordHash)
		return nil, err
	}

	// Проверяем, существует ли юзер с таким логином
	loginExists, err := s.GetUserByLogin(login)
	if loginExists != nil {
		log.Println("User with the given login already exists. Login: ", login)
		return nil, err
	}

	// + доп. проверка на то, занят ли указанный номер телефона
	if phoneNumber != "" {
		phoneExists, err := s.GetUserByPhoneNumber(phoneNumber)
		if phoneExists != nil {
			log.Println("User with the given phone number already exists. Phone number: ", phoneNumber)
			return nil, err
		}
	}

	// Создание пользователя и сохранение записи в БД
	user = &models.User{
		Login:       login,
		Password:    passwordHash,
		Name:        firstName,
		LastName:    lastName,
		MiddleMame:  middleName,
		PhoneNumber: phoneNumber,
	}
	result := s.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (s *ServiceManager) GetUserByLogin(login string) (user *models.User, err error) {
	result := s.db.First(&user, "login = ?", login)
	if result.Error != nil {
		log.Println("Cannot get user by the given login. Error: ", result.Error)
		return nil, errors.UserDoesNotExist
	}
	return user, nil
}

func (s *ServiceManager) GetUserByPhoneNumber(phoneNumber types.PhoneNumber) (user *models.User, err error) {
	result := s.db.First(&user, "phone_number = ?", phoneNumber)
	if result.Error != nil {
		log.Println("Cannot get user by the given phone number. Error: ", result.Error)
		return nil, errors.UserDoesNotExist
	}
	return user, nil
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
2. Login(Authorization) (login, password) bcrypt (access token - 15min)
	- password hash
3. Middleware
JWT(access token 3 min, refresh token 15 - 20)



4. files
5. config
6. Metanit


7. PreCheck, PostCheck, Payment

*/
