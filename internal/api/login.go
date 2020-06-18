package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pass-wall/passwall-server/internal/app"
	"github.com/pass-wall/passwall-server/internal/storage"
	"github.com/pass-wall/passwall-server/model"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

const (
	LoginDeleteSuccess = "Login deleted successfully!"
)

// FindAll ...
func FindAllLogins(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var loginList []model.Login

		fields := []string{"id", "created_at", "updated_at", "url", "username"}
		argsStr, argsInt := SetArgs(r, fields)

		schema := r.Context().Value("schema").(string)
		loginList, err = s.Logins().FindAll(argsStr, argsInt, schema)

		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		// loginList = app.DecryptLoginPasswords(loginList)
		RespondWithJSON(w, http.StatusOK, loginList)
	}
}

// FindByID ...
func FindLoginsByID(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		schema := r.Context().Value("schema").(string)
		login, err := s.Logins().FindByID(uint(id), schema)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		// uLogin, err := app.DecryptLoginPassword(s, login)
		// if err != nil {
		// 	RespondWithError(w, http.StatusInternalServerError, err.Error())
		// 	return
		// }

		RespondWithJSON(w, http.StatusOK, model.ToLoginDTO(login))
	}
}

// Create ...
func CreateLogin(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginDTO model.LoginDTO

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&loginDTO); err != nil {
			RespondWithError(w, http.StatusBadRequest, InvalidRequestPayload)
			return
		}
		defer r.Body.Close()

		if loginDTO.Password == "" {
			generatedPass, err := app.GenerateSecureKey(viper.GetInt("server.generatedPasswordLength"))
			if err != nil {
				RespondWithError(w, http.StatusSeeOther, err.Error())
			}
			loginDTO.Password = generatedPass
		}

		schema := r.Context().Value("schema").(string)
		createdLogin, err := app.CreateLogin(s, &loginDTO, schema)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, model.ToLoginDTO(createdLogin))
	}
}

// Update ...
func UpdateLogin(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		var loginDTO model.LoginDTO
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&loginDTO); err != nil {
			RespondWithError(w, http.StatusBadRequest, InvalidRequestPayload)
			return
		}
		defer r.Body.Close()
		schema := r.Context().Value("schema").(string)
		login, err := s.Logins().FindByID(uint(id), schema)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		updatedLogin, err := app.UpdateLogin(s, login, &loginDTO, schema)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, model.ToLoginDTO(updatedLogin))
	}
}

// Delete ...
func DeleteLogin(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		schema := r.Context().Value("schema").(string)
		login, err := s.Logins().FindByID(uint(id), schema)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		err = s.Logins().Delete(login.ID, schema)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		response := model.Response{
			Code:    http.StatusOK,
			Status:  Success,
			Message: LoginDeleteSuccess,
		}
		RespondWithJSON(w, http.StatusOK, response)
	}
}
