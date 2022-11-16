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
	errGlobal = nil
	netflixDataGlobal, errGlobal = readCSVToObject(filepath)
	if errGlobal != nil {
		log.Default().Println("Cannot fetch data from Netflix CSV")
	} else {
		syncDB(netflixDataGlobal[1:])
	}

	router := mux.NewRouter()
	router.HandleFunc("/tvshows", getTVShowsApiHandler).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:8000", router))

}
