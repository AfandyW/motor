package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/AfandyW/motor/models"
	"github.com/AfandyW/motor/service"
	"github.com/gorilla/mux"
)

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

type Controller struct {
	service service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		service: *service,
	}
}

func (ctrl *Controller) Create(w http.ResponseWriter, r *http.Request) {
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

	err = ctrl.service.Create(motor)
	if err != nil {
		sendResponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
	}

	sendResponse(http.StatusCreated, "success", nil, w)
}

func (ctrl *Controller) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		sendResponse(http.StatusBadRequest, "bad request, data id param is null", nil, w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendResponse(http.StatusBadRequest, "cannot convert id param from string to int", nil, w)
		return
	}

	dataByte, err := io.ReadAll(r.Body)
	if err != nil {
		sendResponse(http.StatusBadRequest, "bad request", nil, w)
		return
	}
	defer r.Body.Close()

	var motor models.Motor
	err = json.Unmarshal(dataByte, &motor)
	if err != nil {
		sendResponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
		return
	}

	motor.ID = idInt

	err = ctrl.service.Update(motor)
	if err != nil {
		sendResponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
		return
	}

	sendResponse(http.StatusOK, "success updated", nil, w)
}

func (ctrl *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		sendResponse(http.StatusBadRequest, "bad request, data id param is null", nil, w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendResponse(http.StatusBadRequest, "cannot convert id param from string to int", nil, w)
		return
	}

	err = ctrl.service.Delete(idInt)
	if err != nil {
		sendResponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
		return
	}

	sendResponse(http.StatusOK, "success delete", nil, w)
}

func (ctrl *Controller) List(w http.ResponseWriter, r *http.Request) {
	motors, err := ctrl.service.List()
	if err != nil {
		sendResponse(http.StatusBadRequest, "internal server error, get motors", err.Error(), w)
		return
	}
	sendResponse(http.StatusOK, "success", motors, w)
}

func (ctrl *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		sendResponse(http.StatusBadRequest, "bad request, data id param is null", nil, w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendResponse(http.StatusBadRequest, "cannot convert id param from string to int", nil, w)
		return
	}

	motor, err := ctrl.service.Get(idInt)
	if err != nil {
		sendResponse(http.StatusBadRequest, "internal server error, get motors", err.Error(), w)
		return
	}
	sendResponse(http.StatusOK, "success", motor, w)
}
