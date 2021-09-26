package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"main/database"
	"main/settings"
	"main/utils"
	"mime/multipart"
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

func (prodData *ProductData) GetMyProducts(DB *gorm.DB, meta utils.MetaData) ([]database.Product, int, error) {
	var result []database.Product
	err := DB.Model(database.Product{}).
		Where(database.Product{UserId: prodData.UserID}).
		Scopes(utils.Paginate(&meta)).Find(&result).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, http.StatusInternalServerError, errors.New("somethings error on server")
	}
	return result, http.StatusOK, nil
}

func (prodData *ProductData) GetMyProduct(DB *gorm.DB) (database.Product, int, error) {
	var result database.Product
	err := DB.Model(database.Product{}).
		Where(database.Product{UserId: prodData.UserID, Id: prodData.ID}).
		First(&result).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			err = errors.New(fmt.Sprintf("product dengan id '%d' tidak ada", prodData.ID))
			return database.Product{}, http.StatusBadRequest, err
		}
		return database.Product{}, http.StatusInternalServerError, errors.New("somethings error on server")
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

func (prodData *ProductData) PrepareUpdate(DB *gorm.DB, body io.ReadCloser) (int, error) {
	dataEdit := map[string]interface{}{}
	var dataFromClient map[string]interface{}
	var total int

	// check if data is exist
	err := DB.Model(database.Product{}).
		Select("COUNT(*) as total").
		Where(database.Product{UserId: prodData.UserID, Id: prodData.ID}).
		First(&total).Error
	if total == 0 {
		err = errors.New(fmt.Sprintf("product dengan id '%d' tidak ada", prodData.ID))
		return http.StatusBadRequest, err
	}

	if body == http.NoBody {
		return http.StatusBadRequest, errors.New("body json di perlukan")
	}
	err = json.NewDecoder(body).Decode(&dataFromClient)
	if err != nil {
		return http.StatusBadRequest, errors.New("format json tidak benar")
	}

	name, err := utils.RequestValidator("name", dataFromClient, utils.RequestDataValidator{
		Type: reflect.String, Null: true,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	price, err := utils.RequestValidator("price", dataFromClient, utils.RequestDataValidator{
		Type: reflect.Float64, Null: true,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	description, err := utils.RequestValidator("description", dataFromClient, utils.RequestDataValidator{
		Type: reflect.String, Null: true,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	if name != nil {
		dataEdit["name"] = name.(string)
		var check database.Product
		err = DB.Model(database.Product{}).
			Where(database.Product{UserId: prodData.UserID, Name: name.(string)}).
			Where(prodData.EditData).First(&check).Error

		if check.Id != 0 {
			return http.StatusBadRequest, errors.New("nama ini sudah ada di product anda")
		}
	}
	if price != nil {
		dataEdit["price"] = price.(float64)
	}
	if description != nil {
		dataEdit["description"] = description.(string)
	}

	if dataEdit != nil {
		prodData.EditData = dataEdit
	}

	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return http.StatusInternalServerError, errors.New("something error from server")
	}
	return http.StatusCreated, nil
}

func (prodData *ProductData) Update(DB *gorm.DB) (int, error) {
	err := DB.Model(database.Product{}).
		Where(database.Product{Id: prodData.ID, UserId: prodData.UserID}).
		Updates(prodData.EditData).Error
	if err != nil {
		return http.StatusInternalServerError, errors.New("something error from server")
	}
	return http.StatusCreated, nil
}

func (prodData *ProductData) Delete(DB *gorm.DB) (int, error) {
	var check database.Product
	err := DB.Model(database.Product{}).
		Where(database.Product{Id: prodData.ID, UserId: prodData.UserID}).
		First(&check).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			err = errors.New(fmt.Sprintf("data dengan ada id '%d' tidak ada", prodData.ID))
			return http.StatusBadRequest, err
		}
		return http.StatusInternalServerError, errors.New("something error from server")
	}

	if check.Id == 0 {
		err = errors.New(fmt.Sprintf("data dengan ada id '%d' tidak ada", prodData.ID))
		return http.StatusBadRequest, err
	}

	err = DB.Delete(database.Product{UserId: prodData.UserID, Id: prodData.ID}).Error
	if err != nil {
		return http.StatusInternalServerError, errors.New("something error from server")
	}
	return http.StatusNoContent, nil
}

func (photoData *ProductPhotoData) UploadPhoto(DB *gorm.DB, file multipart.File, fileInfo *multipart.FileHeader) (utils.FileInfo, int, error) {
	var result database.Product
	var total int
	err := DB.Model(database.Product{}).
		Where(database.Product{UserId: photoData.UserID, Id: photoData.ProductID}).
		First(&result).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			err = errors.New(fmt.Sprintf("product dengan id '%d' tidak ada", photoData.ProductID))
			return utils.FileInfo{}, http.StatusBadRequest, err
		}
		return utils.FileInfo{}, http.StatusInternalServerError, errors.New("somethings error on server")
	}

	err = DB.Model(database.ProductPhoto{}).
		Select("COUNT(*) as total").
		Where(database.ProductPhoto{ProductId: photoData.ProductID}).
		Find(&total).Error
	if err != nil {
		return utils.FileInfo{}, http.StatusInternalServerError, errors.New("somethings error on server")
	}

	if total > 10 {
		return utils.FileInfo{}, http.StatusBadRequest, errors.New("maksimal photo produk hanya sampai 10 saja")
	}

	//var err error
	if file != nil {
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(file)

		allowed := []string{"png", "jpeg", "jpg"}

		filename := fmt.Sprintf("%d_%d_%s", photoData.UserID, photoData.ProductID, utils.RandomStringWithCharset(6, "1234567890ABC"))
		resultFileInfo, status, err := utils.SaveMultipartFile(file, fileInfo, settings.DataSettings.ProductPhotosFolder+filename, allowed)
		if err != nil {
			return utils.FileInfo{}, status, err
		}
		resultFileInfo.Name = filename + "." + resultFileInfo.Extension

		return resultFileInfo, http.StatusCreated, nil
	} else {
		return utils.FileInfo{}, http.StatusBadRequest, errors.New("file tidak ada")
	}
}

func (photoData *ProductPhotoData) GetPhotos(DB *gorm.DB) ([]ProductPhotoResponse, int, error) {
	var check database.Product
	var resultQuery []database.ProductPhoto
	var result []ProductPhotoResponse
	err := DB.Model(database.Product{}).
		Where(database.Product{UserId: photoData.UserID, Id: photoData.ProductID}).
		First(&check).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			err = errors.New(fmt.Sprintf("product dengan id '%d' tidak ada", photoData.ProductID))
			return nil, http.StatusBadRequest, err
		}
		return nil, http.StatusInternalServerError, errors.New("somethings error on server")
	}

	err = DB.Model(database.ProductPhoto{}).
		Where(database.ProductPhoto{ProductId: photoData.ProductID}).
		Find(&resultQuery).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, http.StatusInternalServerError, errors.New("somethings error on server")
	}
	for _, item := range resultQuery {
		result = append(result, ProductPhotoResponse{
			Id:         item.Id,
			Address:    settings.ProductPhotosPrefixUrl + item.Filename,
			CreateTime: item.CreateTime,
			UpdateTime: item.UpdateTime,
		})
	}
	return result, http.StatusOK, err
}

func (photoData *ProductPhotoData) SaveToDB(DB *gorm.DB, fileInfo utils.FileInfo) (int, error) {
	err := DB.Create(&database.ProductPhoto{
		Filename:  fileInfo.Name,
		ProductId: photoData.ProductID,
	}).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return http.StatusBadRequest, errors.New("this email is already registered")
		}
		return http.StatusInternalServerError, errors.New("something error from server")
	}
	return http.StatusCreated, nil
}

func (photoData *ProductPhotoData) Delete(DB *gorm.DB) (int, error) {
	var check database.Product
	var photo database.ProductPhoto
	err := DB.Model(database.Product{}).
		Where(database.Product{Id: photoData.ProductID, UserId: photoData.UserID}).
		First(&check).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			err = errors.New(fmt.Sprintf("produk dengan id '%d' tidak ada", photoData.ProductID))
			return http.StatusBadRequest, err
		}
		return http.StatusInternalServerError, errors.New("something error from server")
	}

	if check.Id == 0 {
		err = errors.New(fmt.Sprintf("produk dengan id '%d' tidak ada", photoData.ProductID))
		return http.StatusBadRequest, err
	}

	//var check database.Product
	err = DB.Model(database.ProductPhoto{}).
		Where(database.ProductPhoto{Id: photoData.PhotoID, ProductId: photoData.ProductID}).
		First(&photo).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			err = errors.New(fmt.Sprintf("foto dengan ada id '%d' tidak ada", photoData.PhotoID))
			return http.StatusBadRequest, err
		}
		return http.StatusInternalServerError, errors.New("something error from server")
	}

	err = DB.Delete(database.ProductPhoto{Id: photoData.PhotoID, ProductId: photoData.ProductID}).Error
	if err != nil {
		return http.StatusInternalServerError, errors.New("something error from server")
	}
	return http.StatusNoContent, nil
}
