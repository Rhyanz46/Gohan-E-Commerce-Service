package user

import (
	"gorm.io/gorm"
	"main/utils"
	"net/http"
)

type User struct {
	DB *gorm.DB
}

func Routes(user *User) *User {
	return user
}

func (user *User) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var err error
	var status int

	if r.Method == "POST" {

		// data validation
		data := LoginData{}
		status, err = data.Validation(r.Body)
		if err != nil {
			utils.ResponseJson(w, status, utils.DataResponse{Message: err.Error()})
			return
		}

		result, status, err := data.GetAccount(user.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.DataResponse{Message: err.Error()})
			return
		}

		token, err := data.CreateToken()
		if err != nil {
			utils.ResponseJson(w, http.StatusInternalServerError, utils.DataResponse{Message: err.Error()})
			return
		}
		utils.ResponseJson(w, status, utils.DataResponse{
			Message: "success to fetch data",
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

func (user *User) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	var err error
	status := http.StatusCreated

	if r.Method == "POST" {
		// user data validation
		data := RegisterData{}
		status, err = data.Validation(r.Body)
		if err != nil {
			utils.ResponseJson(w, status, utils.DataResponse{Message: err.Error()})
			return
		}

		// user data validate to database
		status, err = data.DBValidation()
		if err != nil {
			utils.ResponseJson(w, status, utils.DataResponse{Message: err.Error()})
			return
		}

		status, err = data.Create(user.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.DataResponse{Message: err.Error()})
			return
		}
		utils.ResponseJson(w, status, utils.DataResponse{Message: "your account is success created"})
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (user *User) HandleUserDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		status, tokenData := utils.GetTokenData(utils.GetBearerToken(r))
		if status != http.StatusOK {
			utils.ResponseJson(w, status, utils.DataResponse{
				Message: "you cannot access endpoint",
			})
			return
		}

		userDetail, status, err := tokenData.GetUser(user.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.DataResponse{
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
