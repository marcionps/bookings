package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n, err := fmt.Fprintf(w, "Something")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%d", n)
	})

	_ = http.ListenAndServe(":8080", nil)
}
