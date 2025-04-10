package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type HealthResponse struct {
	Nama      string `json:"nama"`
	NRP       string `json:"nrp"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Uptime    int64  `json:"uptime"`
}

var startTime = time.Now()

func healthHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	uptime := int64(now.Sub(startTime).Seconds())

	response := HealthResponse{
		Nama:      "Rogelio Kenny Arisandi",
		NRP:       "5025231074",
		Status:    "UP",
		Timestamp: now.Unix(),
		Uptime:    uptime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	uptime := int64(now.Sub(startTime).Seconds())

	response := HealthResponse{
		Nama:      "Rogelio Kenny Arisandi",
		NRP:       "5025231074",
		Status:    "UP",
		Timestamp: now.Unix(),
		Uptime:    uptime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/check", checkHandler)
	port := ":8080"
	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
