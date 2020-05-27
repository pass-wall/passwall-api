package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pass-wall/passwall-server/internal/app"
	"github.com/pass-wall/passwall-server/internal/common"
	"github.com/pass-wall/passwall-server/internal/encryption"
	"github.com/pass-wall/passwall-server/model"
	"github.com/spf13/viper"
)

const (
	LoginDeleteSuccess = "Login deleted successfully!"
)

// LoginAPI ...
type LoginAPI struct {
	LoginService app.LoginService
}

// NewLoginAPI ...
func NewLoginAPI(p app.LoginService) LoginAPI {
	return LoginAPI{LoginService: p}
}

// GetHandler ...
func (p *LoginAPI) GetHandler(w http.ResponseWriter, r *http.Request) {
	action := mux.Vars(r)["action"]

	switch action {
	case "backup":
		app.ListBackup(w, r)
	default:
		common.RespondWithError(w, http.StatusNotFound, InvalidRequestPayload)
		return
	}
}

// PostHandler ...
func (p *LoginAPI) PostHandler(w http.ResponseWriter, r *http.Request) {
	action := mux.Vars(r)["action"]

	switch action {
	case "import":
		app.Import(w, r)
	case "export":
		app.Export(w, r)
	case "backup":
		app.Backup(w, r)
	case "restore":
		app.Restore(w, r)
	case "generate-password":
		app.GeneratePassword(w, r)
	case "check-password":
		p.FindSamePassword(w, r)
	default:
		common.RespondWithError(w, http.StatusNotFound, InvalidRequestPayload)
		return
	}
}

// FindAll ...
func (p *LoginAPI) FindAll(w http.ResponseWriter, r *http.Request) {
	var err error
	var loginList []model.Login

	fields := []string{"id", "created_at", "updated_at", "url", "username"}
	argsStr, argsInt := SetArgs(r, fields)

	loginList, err = p.LoginService.FindAll(argsStr, argsInt)

	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	loginList = app.DecryptLoginPasswords(loginList)
	common.RespondWithJSON(w, http.StatusOK, loginList)
}

// FindByID ...
func (p *LoginAPI) FindByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	login, err := p.LoginService.FindByID(uint(id))
	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	passByte, _ := base64.StdEncoding.DecodeString(login.Password)
	login.Password = string(encryption.Decrypt(string(passByte[:]), viper.GetString("server.passphrase")))

	common.RespondWithJSON(w, http.StatusOK, model.ToLoginDTO(login))
}

// Create ...
func (p *LoginAPI) Create(w http.ResponseWriter, r *http.Request) {
	var loginDTO model.LoginDTO

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginDTO); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, InvalidRequestPayload)
		return
	}
	defer r.Body.Close()

	if loginDTO.Password == "" {
		loginDTO.Password = encryption.Password()
	}

	rawPass := loginDTO.Password
	loginDTO.Password = base64.StdEncoding.EncodeToString(encryption.Encrypt(loginDTO.Password, viper.GetString("server.passphrase")))

	createdLogin, err := p.LoginService.Save(model.ToLogin(loginDTO))
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdLogin.Password = rawPass

	common.RespondWithJSON(w, http.StatusOK, model.ToLoginDTO(createdLogin))
}

// Update ...
func (p *LoginAPI) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var loginDTO model.LoginDTO
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginDTO); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	login, err := p.LoginService.FindByID(uint(id))
	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if loginDTO.Password == "" {
		loginDTO.Password = encryption.Password()
	}
	rawPass := loginDTO.Password
	loginDTO.Password = base64.StdEncoding.EncodeToString(encryption.Encrypt(loginDTO.Password, viper.GetString("server.passphrase")))

	login.URL = loginDTO.URL
	login.Username = loginDTO.Username
	login.Password = loginDTO.Password
	updatedLogin, err := p.LoginService.Save(login)
	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	updatedLogin.Password = rawPass
	common.RespondWithJSON(w, http.StatusOK, model.ToLoginDTO(updatedLogin))
}

// Delete ...
func (p *LoginAPI) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	login, err := p.LoginService.FindByID(uint(id))
	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	err = p.LoginService.Delete(login.ID)
	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	response := model.Response{Code: http.StatusOK, Status: Success, Message: LoginDeleteSuccess}
	common.RespondWithJSON(w, http.StatusOK, response)
}

// FindSamePassword ...
func (p *LoginAPI) FindSamePassword(w http.ResponseWriter, r *http.Request) {
	var password model.Password

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&password); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, InvalidRequestPayload)
		return
	}
	defer r.Body.Close()

	urls, err := app.FindSamePassword(&p.LoginService, password)

	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	common.RespondWithJSON(w, http.StatusOK, urls)
}

// Migrate ...
func (p *LoginAPI) Migrate() {
	p.LoginService.Migrate()
}
