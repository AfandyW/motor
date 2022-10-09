package main

import (
	"fmt"
	"net/http"

	"github.com/AfandyW/motor/config"
	"github.com/AfandyW/motor/controllers"
	"github.com/AfandyW/motor/repository"
	"github.com/AfandyW/motor/router"
	"github.com/AfandyW/motor/service"
)

func main() {
	con := config.NewConfig()
	db := config.NewDatabase(&con.DB)

	defer db.Close()

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	controller := controllers.NewController(service)

	r := router.NewRouter(controller)

	fmt.Printf("server run on port %s:%s ", con.API.BaseUrl, con.API.Port)
	http.ListenAndServe(con.API.BaseUrl+":"+con.API.Port, r)
}
