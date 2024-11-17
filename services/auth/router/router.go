// The router package for user authentication
package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		fmt.Fprint(w, "User Auth Service is running.")
	}).Methods("GET")

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

	fmt.Printf("User Auth Service started on port %s\n", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
