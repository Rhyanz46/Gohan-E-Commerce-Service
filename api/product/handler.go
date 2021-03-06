package product

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"main/utils"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	DB *gorm.DB
}

func Routes(admin *Product) *Product {
	return admin
}

func (product *Product) HandleProduct(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method == http.MethodPost {
		var data ProductData
		status, userData := utils.GetTokenData(utils.GetBearerToken(r))
		if status != http.StatusOK {
			utils.ResponseJson(w, status, utils.MessageResponse{
				Message: "you cannot access endpoint",
			})
			return
		}

		status, err = data.Validation(r.Body, userData)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}

		status, err = data.DBValidation(product.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}

		status, err = data.Insert(product.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}
		utils.ResponseJson(w, status, utils.MessageResponse{Message: "success"})
		return
	} else if r.Method == http.MethodGet {
		meta := utils.MetaData{}
		status, userData := utils.GetTokenData(utils.GetBearerToken(r))
		if status != http.StatusOK {
			utils.ResponseJson(w, status, utils.MessageResponse{
				Message: "you cannot access endpoint",
			})
			return
		}

		data := ProductData{UserID: userData.Id}
		result, status, err := data.GetMyProducts(product.DB, meta)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
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

func (product *Product) HandleProductDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.MessageResponse{
			Message: "internal server error",
		})
		return
	}

	status, userData := utils.GetTokenData(utils.GetBearerToken(r))
	if status != http.StatusOK {
		utils.ResponseJson(w, status, utils.MessageResponse{
			Message: "you cannot access endpoint",
		})
		return
	}

	if r.Method == http.MethodGet {
		data := ProductData{
			UserID: userData.Id,
			ID:     uint(productId),
		}
		result, status, err := data.GetMyProduct(product.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}

		utils.ResponseJson(w, status, utils.DataResponse{
			Message: "success",
			Data:    result,
		})
		return
	} else if r.Method == http.MethodPut {
		data := ProductData{
			UserID: userData.Id,
			ID:     uint(productId),
		}

		status, err := data.PrepareUpdate(product.DB, r.Body)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}

		if data.EditData == nil {
			utils.ResponseJson(w, http.StatusOK, utils.MessageResponse{Message: "tidak ada data yang di edit"})
			return
		}

		status, err = data.Update(product.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	} else if r.Method == http.MethodDelete {
		data := ProductData{
			UserID: userData.Id,
			ID:     uint(productId),
		}

		status, err := data.Delete(product.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (product *Product) HandleProductDetailPhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.MessageResponse{
			Message: "internal server error",
		})
		return
	}

	status, userData := utils.GetTokenData(utils.GetBearerToken(r))
	if status != http.StatusOK {
		utils.ResponseJson(w, status, utils.MessageResponse{
			Message: "you cannot access endpoint",
		})
		return
	}

	if r.Method == http.MethodPost && strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		file, fileInfo, err := r.FormFile("photo")
		data := ProductPhotoData{
			UserID:    userData.Id,
			ProductID: uint(productId),
		}

		photoDetail, status, err := data.UploadPhoto(product.DB, file, fileInfo)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}

		status, err = data.SaveToDB(product.DB, photoDetail)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}
		utils.ResponseJson(w, status, utils.MessageResponse{
			Message: "success to upload!",
		})
		return
	} else if r.Method == http.MethodGet {
		data := ProductPhotoData{
			UserID:    userData.Id,
			ProductID: uint(productId),
		}
		photos, status, err := data.GetPhotos(product.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}
		utils.ResponseJson(w, status, utils.DataResponse{
			Message: "success",
			Data:    photos,
		})
		return
	}
}

func (product *Product) HandleProductDetailPhotoDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["product_id"])
	photoId, err := strconv.Atoi(vars["photo_id"])
	if err != nil {
		utils.ResponseJson(w, http.StatusInternalServerError, utils.MessageResponse{
			Message: "internal server error",
		})
		return
	}

	status, userData := utils.GetTokenData(utils.GetBearerToken(r))
	if status != http.StatusOK {
		utils.ResponseJson(w, status, utils.MessageResponse{
			Message: "you cannot access endpoint",
		})
		return
	}

	if r.Method == http.MethodDelete {
		data := ProductPhotoData{
			UserID:    userData.Id,
			ProductID: uint(productId),
			PhotoID:   uint(photoId),
		}
		status, err := data.Delete(product.DB)
		if err != nil {
			utils.ResponseJson(w, status, utils.MessageResponse{Message: err.Error()})
			return
		}
		utils.ResponseJson(w, status, utils.MessageResponse{
			Message: "deleted!",
		})
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
