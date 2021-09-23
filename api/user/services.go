package user

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"io"
	"main/database"
	"main/settings"
	"main/utils"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func (user *LoginData) Validation(body io.ReadCloser) (int, error) {
	var dataFromClient map[string]interface{}
	if body == http.NoBody {
		return http.StatusBadRequest, errors.New("body json di perlukan")
	}
	err := json.NewDecoder(body).Decode(&dataFromClient)
	if err != nil {
		return http.StatusBadRequest, errors.New("format json tidak benar")
	}

	email, err := utils.RequestValidator("email", dataFromClient, utils.RequestDataValidator{
		Type: reflect.String, Max: 10, Min: 5, Null: false, Email: true,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	password, err := utils.RequestValidator("password", dataFromClient, utils.RequestDataValidator{
		Type: reflect.String, Max: 10, Min: 5, Null: false,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	user.Email = email.(string)
	user.Password = utils.ToMD5(password.(string))

	return http.StatusOK, nil
}

func (user *LoginData) GetAccount(DB *gorm.DB) (database.User, int, error) {
	var result database.User
	err := DB.Model(database.User{}).
		Where("email = ? AND password = ?", user.Email, user.Password).
		First(&result).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return result, http.StatusBadRequest, errors.New("username or password is wrong")
		}
		return result, http.StatusInternalServerError, errors.New("something error")
	}
	user.Id = result.Id
	user.Username = result.Username
	return result, http.StatusOK, nil
}

func (user *LoginData) CreateToken() (string, error) {
	token := jwt.NewWithClaims(
		settings.JwtSigningMethod,
		jwt.MapClaims{
			"id":       user.Id,
			"email":    user.Email,
			"username": user.Username,
			"exp":      time.Now().Add(settings.LoginExpirationDuration).Unix(),
		},
	)
	signedToken, err := token.SignedString(settings.JwtSignatureKey)
	if err != nil {
		return "", errors.New("error create token")
	}
	return signedToken, nil
}

func (register *RegisterData) Validation(body io.ReadCloser) (int, error) {
	var dataFromClient map[string]interface{}
	if body == http.NoBody {
		return http.StatusBadRequest, errors.New("body json di perlukan")
	}
	err := json.NewDecoder(body).Decode(&dataFromClient)
	if err != nil {
		return http.StatusBadRequest, errors.New("format json tidak benar")
	}

	username, err := utils.RequestValidator("username", dataFromClient, utils.RequestDataValidator{
		Type: reflect.String, Max: 10, Min: 5, Null: false,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	password, err := utils.RequestValidator("password", dataFromClient, utils.RequestDataValidator{
		Type: reflect.String, Max: 10, Min: 5, Null: false,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	email, err := utils.RequestValidator("email", dataFromClient, utils.RequestDataValidator{
		Type: reflect.String, Max: 10, Min: 5, Null: false, Email: true,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	fullname, err := utils.RequestValidator("fullname", dataFromClient, utils.RequestDataValidator{
		Type: reflect.String, Max: 10, Min: 5, Null: false,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	register.FullName = fullname.(string)
	register.Username = username.(string)
	register.Password = password.(string)
	register.Email = email.(string)

	return http.StatusOK, nil
}

func (register *RegisterData) DBValidation() (int, error) {
	return http.StatusOK, nil
}

func (register *RegisterData) Create(DB *gorm.DB) (int, error) {
	err := DB.Create(&database.User{
		Username: register.Username,
		FullName: register.FullName,
		Email:    register.Email,
		Password: utils.ToMD5(register.Password),
	}).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return http.StatusBadRequest, errors.New("this email is already registered")
		}
		return http.StatusInternalServerError, errors.New("something error from server")
	}
	return http.StatusCreated, nil
}
