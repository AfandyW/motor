package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AfandyW/motor/models"
	"github.com/AfandyW/motor/repository"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// var motors []Motor

func sendResponse(code int, message string, data interface{}, w http.ResponseWriter) {
	resp := models.Response{
		Code:    code,
		Data:    data,
		Message: message,
	}
	dataByte, err := json.Marshal(resp)
	if err != nil {
		resp := models.Response{
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

var db *sql.DB

func main() {

	db, err := sql.Open("postgres", "postgres://root:root@localhost/motors?sslmode=disable")

	if err != nil {
		panic(err.Error())
	}

	if err = db.Ping(); err != nil {
		panic(err.Error())
	}

	defer db.Close()

	r := mux.NewRouter()

	// get
	r.HandleFunc("/api/v1/motorcycles", func(w http.ResponseWriter, r *http.Request) {
		motors, err := repository.GetMotors(db)
		if err != nil {
			sendResponse(http.StatusBadRequest, "internal server error, get motors", err.Error(), w)
			return
		}
		sendResponse(http.StatusOK, "success", motors, w)
	}).Methods(http.MethodGet)

	// update
	r.HandleFunc("/api/v1/motorcycles/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		if id == "" {
			sendResponse(http.StatusBadRequest, "bad request, data id param is null", nil, w)
			return
		}

		motor, err := repository.GetMotor(db, id)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error, get motor return error", err.Error(), w)
			return
		}

		if motor.ID == 0 {
			if err != nil {
				sendResponse(http.StatusNotFound, "data not found", nil, w)
				return
			}
		}

		dataByte, err := io.ReadAll(r.Body)
		if err != nil {
			sendResponse(http.StatusBadRequest, "bad request", nil, w)
			return
		}
		defer r.Body.Close()

		var newMotor models.Motor
		err = json.Unmarshal(dataByte, &newMotor)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
			return
		}

		motor.Name = newMotor.Name
		motor.Price = newMotor.Price

		err = repository.UpdateMotors(db, motor)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
			return
		}

		sendResponse(http.StatusOK, "success updated", nil, w)
	}).Methods(http.MethodPut)

	//create
	r.HandleFunc("/api/v1/motorcycles", func(w http.ResponseWriter, r *http.Request) {
		dataByte, err := io.ReadAll(r.Body)
		if err != nil {
			sendResponse(http.StatusBadRequest, "bad request", nil, w)
		}
		defer r.Body.Close()

		var motor models.Motor
		err = json.Unmarshal(dataByte, &motor)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error", nil, w)
		}

		err = repository.CreateMotors(db, motor)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
		}

		sendResponse(http.StatusCreated, "success", nil, w)
		return
	}).Methods(http.MethodPost)

	// delete
	r.HandleFunc("/api/v1/motorcycles/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		if id == "" {
			sendResponse(http.StatusBadRequest, "bad request, data id param is null", nil, w)
			return
		}

		motor, err := repository.GetMotor(db, id)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error, get motor return err", err.Error(), w)
			return
		}

		if motor.ID == 0 {
			if err != nil {
				sendResponse(http.StatusNotFound, "data not found", nil, w)
				return
			}
		}

		err = repository.Delete(db, id)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error, delete motor return err", err.Error(), w)
			return
		}

		sendResponse(http.StatusOK, "success deleted", nil, w)
	}).Methods(http.MethodDelete)

	port := "8000"
	fmt.Println("server run on port: ", port)
	http.ListenAndServe(":"+port, r)
}
