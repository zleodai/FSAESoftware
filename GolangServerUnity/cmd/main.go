package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zleodai/FSAESoftware/GolangServerUnity/handlers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/test", handlers.Test).Methods("GET")

	router.HandleFunc("/sqliteQuery", handlers.SqliteQuery).Methods("GET")
	router.HandleFunc("/sessions", handlers.GetSessionList).Methods("GET")
	router.HandleFunc("/sessions/{session_id}/laps", handlers.GetLapList).Methods("GET")
	router.HandleFunc("/sessions/{session_id}/fastest_lap", handlers.GetFastestLap).Methods("GET")
	router.HandleFunc("/laps/compare", handlers.GetLapComparisonData).Methods("GET")
	router.HandleFunc("/sessions/{session_id}", handlers.DeleteSession).Methods("DELETE")
	router.HandleFunc("/laps/{lap_id}", handlers.DeleteLap).Methods("DELETE")
	router.HandleFunc("/csvInsert", handlers.CsvInsert).Methods("POST")
	router.HandleFunc("/addRow", handlers.AddRow).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
