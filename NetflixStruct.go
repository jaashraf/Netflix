package main

type NetflixData struct {
	showId      string   `json:"showId"`
	movieType   string   `json:"movieType"`
	title       string   `json:"title"`
	director    []string `json:"director"`
	cast        []string `json:"cast"`
	country     []string `json:"country"`
	dateAdded   string   `json:"dateAdded"`
	releaseYear int      `json:"releaseYear"`
	rating      string   `json:"rating"`
	duration    string   `json:"duration"`
	listedIn    []string `json:"listedIn"`
	description string   `json:"description"`
}
