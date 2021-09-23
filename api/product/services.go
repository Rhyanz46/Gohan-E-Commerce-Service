package product

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"io"
	"main/database"
	"main/utils"
	"net/http"
	"reflect"
	"strings"
)

func (prodData *ProductData) Validation(body io.ReadCloser, user utils.Auth) (int, error) {
	var dataFromClient map[string]interface{}
	if body == http.NoBody {
		return http.StatusBadRequest, errors.New("body json di perlukan")
	}
	err := json.NewDecoder(body).Decode(&dataFromClient)
	if err != nil {
		return http.StatusBadRequest, errors.New("format json tidak benar")
	}

	name, err := utils.RequestValidator("name", dataFromClient, utils.RequestDataValidator{
		Type: reflect.String, Max: 10, Min: 5, Null: false,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	price, err := utils.RequestValidator("price", dataFromClient, utils.RequestDataValidator{
		Type: reflect.Float64, Null: false,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	description, err := utils.RequestValidator("description", dataFromClient, utils.RequestDataValidator{
		Type: reflect.String, Max: 10, Min: 5, Null: false,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	prodData.Name = name.(string)
	prodData.Price = price.(float64)
	prodData.Description = description.(string)
	prodData.UserID = user.Id
	return http.StatusOK, nil
}

func (prodData *ProductData) DBValidation(DB *gorm.DB) (int, error) {
	var total int
	DB.Model(database.Product{}).Select("COUNT(*)").Where(database.Product{UserId: prodData.UserID, Name: prodData.Name}).Find(&total)
	if total > 0 {
		return http.StatusBadRequest, errors.New("product with this name is exist on your product list")
	}
	return http.StatusOK, nil
}

func (prodData *ProductData) GetMyProducts(DB *gorm.DB, meta utils.MetaData, user utils.Auth) ([]database.Product, int, error) {
	var result []database.Product
	err := DB.Debug().Model(database.Product{}).
		Where(database.Product{UserId: user.Id}).
		Scopes(utils.Paginate(&meta)).Find(&result).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, http.StatusInternalServerError, errors.New("somethings error on server")
	}
	return result, http.StatusOK, nil
}

func (prodData *ProductData) Insert(DB *gorm.DB) (int, error) {
	err := DB.Create(&database.Product{
		Name:        prodData.Name,
		Price:       prodData.Price,
		Description: prodData.Description,
		UserId:      prodData.UserID,
	}).Error
	if err != nil {
		return http.StatusInternalServerError, errors.New("something error from server")
	}
	return http.StatusCreated, nil
}
