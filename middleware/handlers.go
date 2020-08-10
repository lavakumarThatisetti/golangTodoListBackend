package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"log"
	"net/http" // used to access the request and response object of the api

	// used to read the environment variable
	// package used to covert string into int type
	"todolistbackend/models" // models package where User schema is defined

	"github.com/google/uuid"
	"github.com/gorilla/mux" // used to get the params from the route

	// package used to read the .env file
	_ "github.com/lib/pq" // postgres golang driver
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "****"
	dbname   = "postgres"
)

// response format
type response struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// check the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db
}

// CreateTodo create a Todo in the postgres db
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// create an empty todo of type models.todo
	var todo models.Todo
	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	todo.UUID = uuid.New().String()
	// call insert user function and pass the user
	insertID := insertTodo(todo)
	// format a response object
	res := response{
		UUID:    insertID,
		Message: "Todo created successfully",
	}
	// send the response
	json.NewEncoder(w).Encode(res)
}

// GetTodo will return a single todo by its id
func GetTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)
	// convert the id type from string to int
	id := params["id"]
	// call the getUser function with user id to retrieve a single user
	todo, err := getTodo(id)
	if err != nil {
		log.Fatalf("Unable to get todo. %v", err)
	}
	// send the response
	json.NewEncoder(w).Encode(todo)
}

// GetAllTodos will return all the users
func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the users in the db
	todos, err := getAllTodos()
	if err != nil {
		log.Fatalf("Unable to get all todos. %v", err)
	}
	// send all the users as response
	json.NewEncoder(w).Encode(todos)
}

// UpdateTodo update tod's detail in the postgres db
func UpdateTodo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)
	// convert the id type from string to int
	id := params["id"]
	// create an empty user of type models.User
	var todo models.Todo
	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	// call update user to update the todo
	updatedRows := updateTodo(id, todo)
	// format the message string
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)
	// format the response message
	res := response{
		UUID:    id,
		Message: msg,
	}
	// send the response
	json.NewEncoder(w).Encode(res)
}

// DeleteTodo delete user's detail in the postgres db
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)
	// convert the id in string to int
	id := params["id"]
	// call the deleteUser, convert the int to int64
	deletedRows := deleteTodo(id)
	// format the message string
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", deletedRows)
	// format the reponse message
	res := response{
		UUID:    id,
		Message: msg,
	}
	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------
// insert one user in the DB
func insertTodo(todo models.Todo) string {

	// create the postgres db connection
	db := createConnection()
	// close the db connection
	defer db.Close()
	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO todolist (uuid, title, text, completed) VALUES ($1, $2, $3,$4) RETURNING uuid`
	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, todo.UUID, todo.Title, todo.Text, todo.Completed).Scan(&todo.UUID)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	fmt.Printf("Inserted a single record  with uuid%v \n", todo.UUID)
	// return the inserted id
	return todo.UUID
}

// get one user from the DB by its userid
func getTodo(id string) (models.Todo, error) {
	// create the postgres db connection
	db := createConnection()
	// close the db connection
	defer db.Close()
	// create a user of models.User type
	var todo models.Todo
	// create the select sql query
	sqlStatement := `SELECT * FROM todolist WHERE uuid=$1`
	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)
	// unmarshal the row object to user
	err := row.Scan(&todo.ID, &todo.Title, &todo.Text, &todo.Completed, &todo.UUID)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return todo, nil
	case nil:
		return todo, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}
	// return empty user on error
	return todo, err
}

// get one user from the DB by its userid
func getAllTodos() ([]models.Todo, error) {
	// create the postgres db connection
	db := createConnection()
	// close the db connection
	defer db.Close()
	var todos []models.Todo
	// create the select sql query
	sqlStatement := `SELECT * FROM todolist`
	// execute the sql statement
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	// close the statement
	defer rows.Close()
	for rows.Next() {
		var todo models.Todo
		err = rows.Scan(&todo.ID, &todo.Title, &todo.Text, &todo.Completed, &todo.UUID)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		// append the user in the users slice
		todos = append(todos, todo)
	}
	// return empty user on error
	return todos, err
}

// update user in the DB
func updateTodo(id string, todo models.Todo) int64 {
	// create the postgres db connection
	db := createConnection()
	// close the db connection
	defer db.Close()
	// create the update sql query
	sqlStatement := `UPDATE todolist SET title=$2, text=$3, completed=$4 WHERE uuid=$1`
	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, todo.Title, todo.Text, todo.Completed)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	// check how many rows affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}

// delete user in the DB
func deleteTodo(id string) int64 {
	// create the postgres db connection
	db := createConnection()
	// close the db connection
	defer db.Close()
	// create the delete sql query
	sqlStatement := `DELETE FROM todolist WHERE uuid=$1`
	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	// check how many rows affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}
