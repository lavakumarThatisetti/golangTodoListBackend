package router

import (
	"todolistbackend/middleware"

	"github.com/gorilla/mux"
)

//Router Configuration
func Router() *mux.Router {
	//Init Router
	router := mux.NewRouter()
	// Route Handlers /Endpoints
	router.HandleFunc("/api/todos", middleware.GetAllTodos).Methods("GET")
	router.HandleFunc("/api/todos/{id}", middleware.GetTodo).Methods("GET")
	router.HandleFunc("/api/todos", middleware.CreateTodo).Methods("POST")
	router.HandleFunc("/api/todos/{id}", middleware.UpdateTodo).Methods("PUT")
	router.HandleFunc("/api/todos/{id}", middleware.DeleteTodo).Methods("DELETE")
	return router
}
