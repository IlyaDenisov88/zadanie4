package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Message  string `json:"message"`
	XResult  string `json:"x-result"`
	XBody    string `json:"x-body"`
}

func enableCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "x-test,ngrok-skip-browser-warning,Content-Type,Accept,Access-Control-Allow-Headers")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		handler(w, r)
	}
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	response := Response{
		Message: "ilya_denisov",
		XResult: r.Header.Get("x-test"),
		XBody:   string(body),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/result4/", enableCORS(resultHandler))
	log.Printf("Server starting on http://0.0.0.0:5000")
	log.Fatal(http.ListenAndServe("0.0.0.0:5000", nil))
}