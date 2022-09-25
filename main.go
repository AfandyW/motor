package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Motor struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var motors []Motor

func sendResponse(code int, message string, data interface{}, w http.ResponseWriter) {
	resp := Response{
		Code:    code,
		Data:    data,
		Message: message,
	}
	dataByte, err := json.Marshal(resp)
	if err != nil {
		resp := Response{
			Code:    http.StatusInternalServerError,
			Data:    nil,
			Message: "Internal Server Error",
		}
		dataByte, _ = json.Marshal(resp)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(dataByte)
}

func main() {
	http.HandleFunc("/api/v1/motorcycles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sendResponse(http.StatusOK, "success", motors, w)
			return
		}

		if r.Method == http.MethodPost {
			dataByte, err := io.ReadAll(r.Body)
			if err != nil{
				sendResponse(http.StatusBadRequest, "bad request", nil, w)
			}
			defer r.Body.Close()

			var motor Motor
			err = json.Unmarshal(dataByte, &motor)
			if err != nil{
				sendResponse(http.StatusInternalServerError, "internal server error", nil, w)
			}

			motors = append(motors, motor)
			sendResponse(http.StatusCreated, "success", nil, w)
			return
		}

		if r.Method == http.MethodPut {
			w.Write([]byte("ini put"))
			return
		}

		if r.Method == http.MethodDelete {
			w.Write([]byte("ini delete"))
			return
		}

		w.Write([]byte("wrong method"))
	})

	port := "8000"
	fmt.Println("server run on port ", port)
	http.ListenAndServe(":"+port, nil)
}
