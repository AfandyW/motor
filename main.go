package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
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

// var motors []Motor

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

func remove(slice []Motor, s int) []Motor {
	return append(slice[:1], slice[s+1:]...)
}

func main() {

	db, err := sql.Open("postgres", "postgres://root:root@localhost/motors?sslmode=disable")

	if err != nil {
		panic(err.Error())
	}

	if err = db.Ping(); err != nil {
		panic(err.Error())
	}

	defer db.Close()

	http.HandleFunc("/api/v1/motorcycles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			rows, err := db.Query("select name, price from motors")
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, get motors", nil, w)
			}

			var motors []Motor
			for rows.Next() {
				var motor Motor

				err = rows.Scan(
					&motor.Name,
					&motor.Price,
				)
				if err != nil {
					sendResponse(http.StatusInternalServerError, "internal server error, get motors, scan data", nil, w)
				}
				motors = append(motors, motor)
			}

			sendResponse(http.StatusOK, "success", motors, w)
			return
		}

		if r.Method == http.MethodPost {
			dataByte, err := io.ReadAll(r.Body)
			if err != nil {
				sendResponse(http.StatusBadRequest, "bad request", nil, w)
			}
			defer r.Body.Close()

			var motor Motor
			err = json.Unmarshal(dataByte, &motor)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error", nil, w)
			}

			_, err = db.Exec("insert into motors(name,price) values($1,$2)", motor.Name, motor.Price)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, get motors", nil, w)
			}

			sendResponse(http.StatusCreated, "success", nil, w)
			return
		}

		if r.Method == http.MethodPut {
			// get query param
			id := r.URL.Query().Get("id")

			if id == "" {
				sendResponse(http.StatusBadRequest, "bad request, data id param is null", nil, w)
				return
			}

			// cek id tersedia atau tdk
			idInt, err := strconv.Atoi(id)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
				return
			}

			// found := idInt <= len(motors)
			// if !found {
			// 	sendResponse(http.StatusNotFound, "data motors not found", nil, w)
			// 	return
			// }

			idInt -= 1

			dataByte, err := io.ReadAll(r.Body)
			if err != nil {
				sendResponse(http.StatusBadRequest, "bad request", nil, w)
			}
			defer r.Body.Close()

			var motor Motor
			err = json.Unmarshal(dataByte, &motor)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error", nil, w)
			}

			// motors[idInt].Name = motor.Name
			// motors[idInt].Price = motor.Price
			sendResponse(http.StatusOK, "success updated", nil, w)
			return
		}

		if r.Method == http.MethodDelete {
			// get query param
			id := r.URL.Query().Get("id")

			if id == "" {
				sendResponse(http.StatusBadRequest, "bad request, data id param is null", nil, w)
				return
			}

			// cek id tersedia atau tdk
			idInt, err := strconv.Atoi(id)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
				return
			}

			// found := idInt <= len(motors)
			// if !found {
			// 	sendResponse(http.StatusNotFound, "data motors not found", nil, w)
			// 	return
			// }

			idInt -= 1
			// motors = remove(motors, idInt)

			sendResponse(http.StatusOK, "success deleted", nil, w)

			return
		}

		w.Write([]byte("wrong method"))
	})

	port := "8000"
	fmt.Println("server run on port ", port)
	http.ListenAndServe(":"+port, nil)
}
