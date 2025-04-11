package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zleodai/FSAESoftware/GolangServerUnity/handlers"
	"github.com/zleodai/FSAESoftware/GolangServerUnity/helpers"
)

func main() {
	db, err := sql.Open("sqlite3", "./database.db") // Open DB once
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Close DB when main function exits

	err = helpers.SetupDatabaseSchema(db) // Setup schema once at startup
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/test", handlers.Test).Methods("GET")

	router.HandleFunc("/sqliteQuery", func(w http.ResponseWriter, r *http.Request) {
		handlers.SqliteQuery(w, r, db) // Pass db here
	}).Methods("GET")
	router.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetSessionList(w, r, db) // Pass db here
	}).Methods("GET")
	router.HandleFunc("/sessions/{session_id}/laps", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetLapList(w, r, db) // Pass db here
	}).Methods("GET")
	router.HandleFunc("/sessions/{session_id}/fastest_lap", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetFastestLap(w, r, db) // Pass db here
	}).Methods("GET")
	router.HandleFunc("/laps/compare", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetLapComparisonData(w, r, db) // Pass db here
	}).Methods("GET")
	router.HandleFunc("/sessions/{session_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteSession(w, r, db) // Pass db here
	}).Methods("DELETE")
	router.HandleFunc("/laps/{lap_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteLap(w, r, db) // Pass db here
	}).Methods("DELETE")
	router.HandleFunc("/csvInsert", func(w http.ResponseWriter, r *http.Request) {
		handlers.CsvInsert(w, r, db) // Pass db here
	}).Methods("POST")
	router.HandleFunc("/addRow", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddRow(w, r, db) // Pass db here
	}).Methods("POST")
	router.HandleFunc("/clearDatabase", func(w http.ResponseWriter, r *http.Request) {
		handlers.ClearDatabase(w, r, db) // Pass db here
	}).Methods("DELETE")
	router.HandleFunc("/appendCSV", func(w http.ResponseWriter, r *http.Request) {
		handlers.AppendCSV(w, r, db) // Pass db here
	}).Methods("POST")
	router.HandleFunc("/databaseToCSV", func(w http.ResponseWriter, r *http.Request) {
		handlers.DatabaseToCSV(w, r, db) // Pass db here
	}).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
