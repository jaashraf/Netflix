package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "welcome to netflix binge watching")
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}
