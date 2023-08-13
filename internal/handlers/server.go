package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/read_all", ReadAll).Methods(http.MethodGet)
	router.HandleFunc("/create", Create).Methods(http.MethodPost)
	router.HandleFunc("/notes/{id:[0-9]+}", HandleNote).
		Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	//router.HandleFunc("/notebooks", HandleNote).Methods(http.MethodPatch)
	//router.HandleFunc("/notebooks", HandleNote).Methods(http.MethodDelete)

	return router
}
