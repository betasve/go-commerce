// The main package for our inventory
package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		fmt.Fprint(w, "Inventory Service is running.")
	})

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

	fmt.Printf("Inventory Service started on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
