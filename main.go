package main

import (
	"fmt"
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

func main() {
	http.HandleFunc("/api/v1/motorcycles", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sucess"))
	})

	port := "8000"
	fmt.Println("server run on port ", port)
	http.ListenAndServe(":"+port, nil)
}
