package app

import (
	"encoding/base64"

	"github.com/passwall/passwall-server/internal/storage"
	"github.com/passwall/passwall-server/model"
	"github.com/spf13/viper"
)

// CreateLogin creates a login and saves it to the store
func CreateLogin(s storage.Store, dto *model.LoginDTO, schema string) (*model.Login, error) {

	rawPass := dto.Password
	dto.Password = base64.StdEncoding.EncodeToString(Encrypt(dto.Password, viper.GetString("server.passphrase")))

	createdLogin, err := s.Logins().Save(model.ToLogin(dto), schema)
	if err != nil {
		return nil, err
	}

	createdLogin.Password = rawPass
	return createdLogin, nil

}

// UpdateLogin updates the login with the dto and applies the changes in the store
func UpdateLogin(s storage.Store, login *model.Login, dto *model.LoginDTO, schema string) (*model.Login, error) {

	rawPass := dto.Password
	dto.Password = base64.StdEncoding.EncodeToString(Encrypt(dto.Password, viper.GetString("server.passphrase")))

	login.Title = dto.Title
	login.URL = dto.URL
	login.Username = dto.Username
	updatedLogin, err := s.Logins().Save(login, schema)
	if err != nil {
		return nil, err
	}
	updatedLogin.Password = rawPass
	return updatedLogin, nil
}

// DecryptLoginPassword decrypts password
func DecryptLoginPassword(s storage.Store, login *model.Login) (*model.Login, error) {
	passByte, _ := base64.StdEncoding.DecodeString(login.Password)
	login.Password = string(Decrypt(string(passByte[:]), viper.GetString("server.passphrase")))

	return login, nil
}

// DecryptLoginPasswords ...
func DecryptLoginPasswords(loginList []model.Login) {
	for i := range loginList {
		if loginList[i].Password == "" {
			continue
		}
		passByte, _ := base64.StdEncoding.DecodeString(loginList[i].Password)
		passB64 := string(Decrypt(string(passByte[:]), viper.GetString("server.passphrase")))
		loginList[i].Password = passB64
	}
}
