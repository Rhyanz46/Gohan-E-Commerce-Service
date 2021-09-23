package product

import (
	"gorm.io/gorm"
	"main/utils"
	"net/http"
)

type Product struct {
	DB *gorm.DB
}

func Routes(admin *Product) *Product {
	return admin
}

func (product *Product) HandleProduct(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method == "POST" {
		var data ProductData
		status, userData := utils.GetTokenData(utils.GetBearerToken(r))
		if status != http.StatusOK {
			utils.ResponseJson(w, status, utils.DataResponse{
				Message: "you cannot access endpoint",
			})
			return
		}

		status, err = data.Validation(r.Body, userData)
		if err != nil {
			utils.ResponseJson(w, status, utils.DataResponse{Message: err.Error()})
			return
		}

		status, err = data.DBValidation(product.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.DataResponse{Message: err.Error()})
			return
		}

		status, err = data.Insert(product.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.DataResponse{Message: err.Error()})
			return
		}
		utils.ResponseJson(w, status, utils.DataResponse{Message: "success"})
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (product *Product) HandleProductList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var data ProductData
		meta := utils.MetaData{}
		status, userData := utils.GetTokenData(utils.GetBearerToken(r))
		if status != http.StatusOK {
			utils.ResponseJson(w, status, utils.DataResponse{
				Message: "you cannot access endpoint",
			})
			return
		}

		result, status, err := data.GetMyProducts(product.DB, meta, userData)
		if err != nil {
			utils.ResponseJson(w, status, utils.DataResponse{Message: err.Error()})
			return
		}

		utils.ResponseJson(w, status, utils.MetaDataResponse{
			Message: "success",
			Data:    result,
			Meta:    meta,
		})
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (product *Product) HandleUserDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	} else if r.Method == "PUT" {

	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
