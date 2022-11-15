package services

import (
	"banking/pkg/errors"
	"banking/pkg/utils"
	"fmt"
	"log"
)

type authInterface interface {
	RefreshAccessToken(login string, password []byte) (accessToken string, err error)
}

func (s *ServiceManager) RefreshAccessToken(login string, password []byte) (accessToken string, err error) {
	user, err := s.GetUserByLogin(login)
	if err != nil {
		log.Println("Cannot refresh access token. Error: ", err)
		return "", err
	}

	// In case when password is incorrect
	if ok := user.CheckPassword(password); ok != true {
		return "", errors.InvalidPassword
	}

	if accessToken, err = utils.GenerateAccessToken(user.Login); err != nil {
		fmt.Println("Cannot refresh access token. Error: ", err)
		return "", err
	}

	// Here we return new access token and update user record in DB
	user.AccessToken = accessToken
	s.DB.Save(&user)
	return accessToken, nil
}
