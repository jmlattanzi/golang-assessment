package main

import (
	"fmt"
	"net/http"

	"github.com/jmlattanzi/golang-assessment/personcontroller"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("[*] Starting application....")
	router := mux.NewRouter()

	router.HandleFunc("/person", personcontroller.GetAllPeopleHandler).Methods("GET")
	router.HandleFunc("/person", personcontroller.AddPersonHandler).Methods("POST")
	router.HandleFunc("/person/{id}", personcontroller.GetSinglePersonHandler).Methods("GET")
	router.HandleFunc("/person/{id}", personcontroller.EditPersonHandler).Methods("PUT")
	router.HandleFunc("/person/{id}", personcontroller.DeletePersonHandler).Methods("DELETE")

	// export to csv
	router.HandleFunc("/export", personcontroller.ExportToCSVHandler).Methods("GET")
	// import csv
	router.HandleFunc("/import", personcontroller.ImportCSVHandler).Methods("POST")

	http.ListenAndServe(":8000", router)
}
