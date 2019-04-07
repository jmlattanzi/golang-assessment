package main

import (
	"fmt"
	"net/http"

	"github.com/jmlattanzi/assessment/controller"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("[*] Starting application....")
	router := mux.NewRouter()

	router.HandleFunc("/person", controller.GetAllPeopleHandler).Methods("GET")
	router.HandleFunc("/person", controller.AddPersonHandler).Methods("POST")
	router.HandleFunc("/person/{id}", controller.GetSinglePersonHandler).Methods("GET")
	router.HandleFunc("/person/{id}", controller.EditPersonHandler).Methods("PUT")
	router.HandleFunc("/person/{id}", controller.DeletePersonHandler).Methods("DELETE")

	// export to csv
	router.HandleFunc("/export", controller.ExportToCSVHandler).Methods("GET")
	// import csv
	router.HandleFunc("/import", controller.ImportCSVHandler).Methods("POST")

	http.ListenAndServe(":8000", router)
}
