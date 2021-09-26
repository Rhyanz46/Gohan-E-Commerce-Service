package account

import (
	"gorm.io/gorm"
	"main/utils"
	"net/http"
)

type Account struct {
	DB *gorm.DB
}

func Routes(account *Account) *Account {
	return account
}

func (acc *Account) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var err error
	var status int

	if r.Method == "POST" {

		// data validation
		data := LoginData{}
		status, err = data.Validation(r.Body)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}

		result, status, err := data.GetAccount(acc.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}

		token, err := data.CreateToken()
		if err != nil {
			utils.ResponseJson(w, http.StatusInternalServerError, utils.MessageResponse{Message: err.Error()})
			return
		}
		utils.ResponseJson(w, status, utils.DataResponse{
			Message: "success to login",
			Data: map[string]interface{}{
				"id":          result.Id,
				"email":       result.Email,
				"fullname":    result.FullName,
				"username":    result.Username,
				"create_time": result.CreateTime,
				"token":       token,
			},
		})
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (acc *Account) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var err error
	status := http.StatusCreated

	if r.Method == "POST" {
		// acc data validation
		data := RegisterData{}
		status, err = data.Validation(r.Body)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}

		// acc data validate to database
		status, err = data.DBValidation()
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}

		status, err = data.Create(acc.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}
		utils.ResponseJson(w, status, utils.MessageResponse{Message: "your account is success created"})
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (acc *Account) HandleDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		status, tokenData := utils.GetTokenData(utils.GetBearerToken(r))
		if status != http.StatusOK {
			utils.ResponseJson(w, status, utils.DataResponse{
				Message: "you cannot access endpoint",
			})
			return
		}

		userDetail, status, err := tokenData.GetUser(acc.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{
				Message: err.Error(),
			})
			return
		}

		utils.ResponseJson(w, status, utils.DataResponse{
			Message: "you cannot access endpoint",
			Data: map[string]interface{}{
				"id":       tokenData.Id,
				"email":    tokenData.Email,
				"username": tokenData.Username,
				"fullname": userDetail.FullName,
			},
		})
		return
	} else if r.Method == "PUT" {
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
