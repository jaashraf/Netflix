package main

type NetflixData struct {
	ShowId      string   `json:"showId"`
	MovieType   string   `json:"movieType"`
	Title       string   `json:"title"`
	Director    []string `json:"director"`
	Cast        []string `json:"cast"`
	Country     []string `json:"country"`
	DateAdded   string   `json:"dateAdded"`
	ReleaseYear int      `json:"releaseYear"`
	Rating      string   `json:"rating"`
	Duration    string   `json:"duration"`
	ListedIn    []string `json:"listedIn"`
	Description string   `json:"description"`
}
