package main

import (
	"diary/internal/configs"
	"diary/internal/handlers"
	"log"
	"net/http"
)

func main() {

	handlers.LoadNotes()

	config, err := configs.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	address := config.Host + config.Port
	router := handlers.InitRoutes()

	//srv := http.Server{
	//	Addr:    address,
	//	Handler: router,
	//}
	//err = srv.ListenAndServe()
	log.Println("start")
	err = http.ListenAndServe(address, router)
	if err != nil {
		log.Println("listen and serv error:", err)
	}

}
