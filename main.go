package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Motor struct {
	ID    int     `json:"id"`
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

	r.HandleFunc("/api/v1/motorcycles", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("select id, name, price from motors")
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error, get motors", nil, w)
		}

		var motors []Motor
		for rows.Next() {
			var motor Motor

			err = rows.Scan(
				&motor.ID,
				&motor.Name,
				&motor.Price,
			)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, get motors, scan data", nil, w)
			}
			motors = append(motors, motor)
		}

		sendResponse(http.StatusOK, "success", motors, w)
	}).Methods(http.MethodGet)

	r.HandleFunc("/api/v1/motorcycles/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		if id == "" {
			sendResponse(http.StatusBadRequest, "bad request, data id param is null", nil, w)
			return
		}

		rows, err := db.Query("select id, name, price from motors where id = $1", id)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error, get motors", nil, w)
			return
		}

		var motor Motor

		if rows.Next() {
			err = rows.Scan(
				&motor.ID,
				&motor.Name,
				&motor.Price,
			)

			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, scan data return err", nil, w)
				return
			}
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

		var newMotor Motor
		err = json.Unmarshal(dataByte, &newMotor)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
			return
		}

		motor.Name = newMotor.Name
		motor.Price = newMotor.Price

		_, err = db.Exec("update motors set name=$2, price=$3 where id=$1", motor.ID, motor.Name, motor.Price)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error, update motors", err.Error(), w)
			return
		}

		sendResponse(http.StatusOK, "success updated", nil, w)
	}).Methods(http.MethodPut)

	r.HandleFunc("/api/v1/motorcycles", func(w http.ResponseWriter, r *http.Request) {
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
	}).Methods(http.MethodPost)

	

	http.HandleFunc("/api/v1/motorcycles", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodDelete {
			// get query param
			id := r.URL.Query().Get("id")

			if id == "" {
				sendResponse(http.StatusBadRequest, "bad request, data id param is null", nil, w)
				return
			}

			rows, err := db.Query("select id, name, price from motors where id = $1", id)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, get motors", nil, w)
				return
			}

			var motor Motor

			if rows.Next() {
				err = rows.Scan(
					&motor.ID,
					&motor.Name,
					&motor.Price,
				)

				if err != nil {
					sendResponse(http.StatusInternalServerError, "internal server error, scan data return err", nil, w)
					return
				}
			}

			if motor.ID == 0 {
				if err != nil {
					sendResponse(http.StatusNotFound, "data not found", nil, w)
					return
				}
			}

			_, err = db.Exec("delete from motors where id = $1", motor.ID)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, delete motors return err", nil, w)
				return
			}

			sendResponse(http.StatusOK, "success deleted", nil, w)

			return
		}

		w.Write([]byte("wrong method"))
	})

	port := "8000"
	fmt.Println("server run on port ", port)
	http.ListenAndServe(":"+port, r)
}
