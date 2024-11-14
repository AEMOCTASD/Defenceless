package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type StringData struct {
	Value string `json:"value"`
}

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	err = createTableIfNotExists()
	if err != nil {
		log.Fatal("Error setting up the database schema:", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/add", AddStringHandler).Methods("POST")
	router.HandleFunc("/get/{id}", GetStringHandler).Methods("GET")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createTableIfNotExists() error {
	query := `
    CREATE TABLE IF NOT EXISTS strings (
        id SERIAL PRIMARY KEY,
        value TEXT NOT NULL
    );`
	_, err := db.Exec(query)
	return err
}

func AddStringHandler(w http.ResponseWriter, r *http.Request) {
	var stringData StringData
	err := json.NewDecoder(r.Body).Decode(&stringData)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var id int
	err = db.QueryRow("INSERT INTO strings(value) VALUES($1) RETURNING id", stringData.Value).Scan(&id)
	if err != nil {
		http.Error(w, "Error inserting into database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "String added with ID: %d", id)
}

func GetStringHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var value string
	err := db.QueryRow(fmt.Sprintf("SELECT value FROM strings WHERE id = %s", id)).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "String not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StringData{Value: value})
}
