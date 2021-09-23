package utils

import (
	"encoding/json"
	"net/http"
)

type MetaData struct {
	Limit       int `json:"limit"`
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	TotalData   int `json:"total_data"`
	NextPage    int `json:"next_page"`
	PrevPage    int `json:"prev_page"`
}

type DataResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type MetaDataResponse struct {
	Message string      `json:"message"`
	Meta    MetaData    `json:"meta"`
	Data    interface{} `json:"data"`
}

func ResponseJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}
