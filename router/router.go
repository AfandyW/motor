package router

import (
	"net/http"

	"github.com/AfandyW/motor/controllers"
	"github.com/gorilla/mux"
)

func NewRouter(ctrl *controllers.Controller) *mux.Router {
	r := mux.NewRouter()

	// list
	r.HandleFunc("/api/v1/motorcycles", ctrl.List).Methods(http.MethodGet)

	// get
	r.HandleFunc("/api/v1/motorcycles/{id}", ctrl.Get).Methods(http.MethodGet)

	// update
	r.HandleFunc("/api/v1/motorcycles/{id}", ctrl.Update).Methods(http.MethodPut)

	//create
	r.HandleFunc("/api/v1/motorcycles", ctrl.Create).Methods(http.MethodPost)

	// delete
	r.HandleFunc("/api/v1/motorcycles/{id}", ctrl.Delete).Methods(http.MethodDelete)
	
	return r
}
