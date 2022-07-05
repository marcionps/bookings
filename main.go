package main

import (
	"fmt"
	"net/http"
)

const portNumber = ":8080"

func Home(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "This is the home page")
	if err != nil {
		fmt.Println(err)
	}
}

func About(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "This is the about page")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Printf("Starting application on port %s\n", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}
