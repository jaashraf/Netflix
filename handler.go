package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var netflixData []NetflixData
var err error

func addShowOnNetflix(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var netflixRequestData NetflixData
	json.NewDecoder(r.Body).Decode(&netflixRequestData)
	fmt.Println(netflixRequestData)
}

func getTVShowsApiHandler(w http.ResponseWriter, r *http.Request) {
	exeStartTime := time.Now()

	if r.URL.Query().Has("n") {
		netflixData, err = getTVShowsAsPerCountHandler(r)
	} else if r.URL.Query().Has("movieType") {
		netflixData, err = getTVShowsByMovieTypeHandler(r)
	} else if r.URL.Query().Has("country") {
		netflixData, err = getTVShowsByCountryHandler(r)
	} else if r.URL.Query().Has("startDate") && r.URL.Query().Has("endDate") {
		netflixData, err = getTVShowsBetweenDatesHandler(r)
	} else {
		json.NewEncoder(w).Encode("Invalid Input")
	}
	exeEndTime := time.Since(exeStartTime)
	fmt.Println(exeEndTime)
	w.Header().Set("X-Time-Execute", exeEndTime.String())
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(errGlobal)
	} else {
		json.NewEncoder(w).Encode(netflixData)
	}
}

func getTVShowsAsPerCountHandler(r *http.Request) ([]NetflixData, error) {
	count, _ := strconv.ParseInt(r.URL.Query().Get("n"), 10, 64)
	netflixData, _ := filterByTypeAndCount("TV Show", int(count))
	log.Default().Print("TV Shows as per the given count ", count, " : \n\n\n", netflixData)
	if netflixData != nil {
		return netflixData[0:count], nil
	}
	return nil, errors.New("cannot fetch TV shows by count")
}

func getTVShowsByMovieTypeHandler(r *http.Request) ([]NetflixData, error) {
	movieType := r.URL.Query().Get("movieType")
	//netflixData := filterByListedIn(movieType, filterBYType("TV Show", netflixData))
	netflixData := filterByTypeAndMovieType(movieType)
	log.Default().Print("TV Shows of Movie Type ", movieType, " : \n\n\n", netflixData)
	if netflixData != nil {
		return netflixData, nil
	}
	return nil, errors.New("Cannot fetch TV shows by Movie Type")
}

func getTVShowsByCountryHandler(r *http.Request) ([]NetflixData, error) {
	countryName := r.URL.Query().Get("country")
	netflixData, err := readCSVToObject(filepath)
	if err == nil {
		netflixData := filterByCountry(countryName, netflixData)
		log.Default().Print("TV Shows of Country ", countryName, " : \n\n\n", netflixData)
		if netflixData != nil {
			return netflixData, nil
		}
	}
	return nil, errors.New("Cannot fetch TV shows by country")
}

func getTVShowsBetweenDatesHandler(r *http.Request) ([]NetflixData, error) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	netflixData, err := readCSVToObject(filepath)
	if err == nil {
		netflixData, _ := filterByAddedDate(startDate, endDate, netflixData)
		log.Default().Print("TV Shows between ", startDate, "and ", endDate, " : ", netflixData)
		if netflixData != nil {
			return netflixData, nil
		}
	}
	return nil, errors.New("Cannot fetch TV Shows between given dates")
}
