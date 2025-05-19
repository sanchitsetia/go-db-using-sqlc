package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sanchitsetia/go-db-using-sqlc/db"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	var u db.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		http.Error(w, "error while decoding the JSON to struct", http.StatusBadRequest)
	}
	_, queries := db.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Important to release resources
	userCreated, err := queries.CreateUser(ctx, db.CreateUserParams{Username: u.Username, Email: u.Email, Password: u.Password})
	if err != nil {
		http.Error(w, "error while creating user", http.StatusInternalServerError)

	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userCreated)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Error reading the path parameter - id", http.StatusBadRequest)
	}
	_, queries := db.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Important to release resources
	id_converted, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "error while converting string to int- id", http.StatusBadRequest)
	}
	user, err := queries.GetUser(ctx, int32(id_converted))
	if err != nil {
		http.Error(w, "error while querying the user", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(user)
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	_, queries := db.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	users, err := queries.ListUsers(ctx)
	if err != nil {
		http.Error(w, "error while querying for list users", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(users)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Error reading the path parameter - id", http.StatusBadRequest)
	}
	_, queries := db.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Important to release resources
	id_converted, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "error while converting string to int- id", http.StatusBadRequest)
	}
	err = queries.DeleteUser(ctx, int32(id_converted))
	if err != nil {
		http.Error(w, "error while deleting user in DB", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "user deleted successfully",
	})
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", listUsers).Methods("GET")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	db.InitDb()
	dbConn, _ := db.GetDB()
	defer dbConn.Close()
	log.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
