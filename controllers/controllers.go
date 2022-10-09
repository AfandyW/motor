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
	write(resp, code, w)
}

func sendResponseError(err error, w http.ResponseWriter) {
	if v, ok := err.(*models.Errors); ok {
		write(v, v.Code, w)
		return
	}

	data := models.NewInternalServerError(err.Error())
	write(data, http.StatusInternalServerError, w)
}

func write(data interface{}, code int, w http.ResponseWriter) {
	dataByte, _ := json.Marshal(data)
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
		sendResponseError(models.NewInternalServerError(err.Error()), w)
		return
	}
	defer r.Body.Close()

	var motor models.Motor
	err = json.Unmarshal(dataByte, &motor)
	if err != nil {
		sendResponseError(models.NewInternalServerError(err.Error()), w)
		return
	}

	if err = motor.Validate(); err != nil {
		sendResponseError(err, w)
		return
	}

	err = ctrl.service.Create(motor)
	if err != nil {
		sendResponseError(err, w)
		return
	}

	sendResponse(http.StatusCreated, "success created data", nil, w)
}

func (ctrl *Controller) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		sendResponseError(models.NewBadRequestParameterID(), w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendResponseError(models.NewFailedToConvertData(), w)
		return
	}

	dataByte, err := io.ReadAll(r.Body)
	if err != nil {
		sendResponseError(models.NewInternalServerError(err.Error()), w)
		return
	}
	defer r.Body.Close()

	var motor models.Motor
	err = json.Unmarshal(dataByte, &motor)
	if err != nil {
		sendResponseError(models.NewInternalServerError(err.Error()), w)
		return
	}

	motor.ID = idInt

	err = ctrl.service.Update(motor)
	if err != nil {
		sendResponseError(err, w)
		return
	}

	sendResponse(http.StatusOK, "success updated", nil, w)
}

func (ctrl *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		sendResponseError(models.NewBadRequestParameterID(), w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendResponseError(models.NewFailedToConvertData(), w)
		return
	}

	err = ctrl.service.Delete(idInt)
	if err != nil {
		sendResponseError(err, w)
		return
	}

	sendResponse(http.StatusOK, "success delete", nil, w)
}

func (ctrl *Controller) List(w http.ResponseWriter, r *http.Request) {
	motors, err := ctrl.service.List()
	if err != nil {
		sendResponseError(err, w)
		return
	}
	sendResponse(http.StatusOK, "success", motors, w)
}

func (ctrl *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		sendResponseError(models.NewBadRequestParameterID(), w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendResponseError(models.NewFailedToConvertData(), w)
		return
	}

	motor, err := ctrl.service.Get(idInt)
	if err != nil {
		sendResponseError(err, w)
		return
	}
	sendResponse(http.StatusOK, "success", motor, w)
}
