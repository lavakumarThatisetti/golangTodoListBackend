package main

import (
	"log"
	"net/http"
	"todolistbackend/router"

	"github.com/gorilla/handlers"
)

func main() {
	//Init Router
	router := router.Router()
	log.Fatal(http.ListenAndServe(":8000",
		handlers.CORS(
			handlers.AllowedHeaders(
				[]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods(
				[]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(router)))

}
