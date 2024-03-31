package utils

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Status     bool   `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Response   any    `json:"data"`
}
type Empty struct {
}

func (api *ApiResponse) ToJson() []byte {
	bytes, err := json.Marshal(api)
	if err != nil {
		return []byte("{\"status\": false, \"response\": null, \"statuscode\": -1}")
	} else {
		return bytes
	}
}

func CreateErrorResponse(w http.ResponseWriter, code int, reason string, err error) {
	temp := ""
	if err != nil {
		temp = err.Error()
	}

	apiResponse := ApiResponse{
		Status:     false,
		Message:    reason,
		StatusCode: code,
		Response:   temp,
	}
	w.WriteHeader(code)
	w.Write(apiResponse.ToJson())
}

func CreateResponse[T any](w http.ResponseWriter, reason string, data *T) {
	apiResponse := ApiResponse{
		Status:     true,
		Message:    reason,
		StatusCode: 200,
		Response:   data,
	}

	w.Write(apiResponse.ToJson())
}

func CreateDeleteResponse[T any](w http.ResponseWriter, reason string, data *T) {
	apiResponse := ApiResponse{
		Status:     true,
		Message:    reason,
		StatusCode: 204,
		Response:   data,
	}

	w.Write(apiResponse.ToJson())
}
