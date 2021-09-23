package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"main/database"
	"main/settings"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type RequestDataValidator struct {
	Null  bool
	Type  reflect.Kind
	Max   int
	Min   int
	Email bool
}

type Auth struct {
	GetToken bool
	Id       uint
	Username string
	Email    string
}

func (auth *Auth) GetUser(DB *gorm.DB) (result database.User, status int, err error) {
	err = DB.Model(database.User{}).Where("id = ?", auth.Id).First(&result).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			status = http.StatusUnauthorized
			err = errors.New("token is not valid")
			return
		}
		status = http.StatusInternalServerError
		err = errors.New("server have a error")
		return
	}
	status = http.StatusOK
	return
}

func RequestValidator(field string, dataTemp map[string]interface{}, validator RequestDataValidator) (interface{}, error) {
	data := dataTemp[field]

	// null validation
	if !validator.Null && data == nil {
		return nil, errors.New("field " + field + " harus ada")
	} else if validator.Null && data == nil {
		return nil, nil
	}

	// data type validation
	if validator.Type == reflect.Int {
		validator.Type = reflect.Float64
	}

	if reflect.TypeOf(data).Kind() != validator.Type {
		return nil, errors.New("field " + field + " harusnya bertipe " + validator.Type.String())
	}

	if validator.Email {
		if reflect.TypeOf(data).Kind() != reflect.String || !IsEmail(data.(string)) {
			return nil, errors.New("field " + field + " merupakan tidak email yang valid")
		}
	}

	return data, nil
}

func GetTokenData(token string) (int, Auth) {
	if token == "" {
		return http.StatusUnauthorized, Auth{}
	}
	bearerToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return settings.JwtSignatureKey, nil
	})
	if err != nil || !bearerToken.Valid {
		return http.StatusUnauthorized, Auth{}
	}
	dat := bearerToken.Claims.(jwt.MapClaims)
	useridString := fmt.Sprintf("%v", dat["id"])
	userId, _ := strconv.Atoi(useridString)
	return http.StatusOK, Auth{
		GetToken: true,
		Id:       uint(userId),
		Username: dat["username"].(string),
		Email:    dat["email"].(string),
	}
}

func GetBearerToken(r *http.Request) string {
	var token string
	token = r.Header.Get("Authorization")
	if token != "" {
		splitToken := strings.Split(token, " ")
		if len(splitToken) == 2 {
			if strings.ToLower(splitToken[0]) != "bearer" {
				return ""
			}
			return splitToken[1]
		}
	}
	token = r.Header.Get("authorization")
	if token != "" {
		splitToken := strings.Split(token, " ")
		if len(splitToken) == 2 {
			if strings.ToLower(splitToken[0]) != "bearer" {
				return ""
			}
			return splitToken[1]
		}
	}
	return ""
}
