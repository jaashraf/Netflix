package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const filepath = "Resources/netflix_titles.csv"

var netflixDataGlobal = []NetflixData{}
var errGlobal error

func main() {

	excelToDBSyncScheduler()
	router := mux.NewRouter()
	router.HandleFunc("/tvshows", getTVShowsApiHandler).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:8000", router))

}
