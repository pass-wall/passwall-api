package app

import (
	"github.com/passwall/passwall-server/internal/storage"
	"github.com/passwall/passwall-server/model"
)

// FindSamePassword ...
func FindSamePassword(s storage.Store, password model.Password, schema string) (model.URLs, error) {

	loginList, err := s.Logins().All(schema)
	if err != nil {
		return *new(model.URLs), nil
	}
	DecryptLoginPasswords(loginList)
	newUrls := model.URLs{Items: []string{}}

	for _, login := range loginList {
		if login.Password == password.Password {
			newUrls.AddItem(login.URL)
		}
	}

	return newUrls, err
}
