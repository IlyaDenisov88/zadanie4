package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// APIResponse представляет структуру ответа сервера
type APIResponse struct {
	Message  string `json:"message"`
	Header   string `json:"x-header-value"`
	Body     string `json:"request_body"`
}

// Middleware для обработки CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем CORS заголовки
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "x-test, ngrok-skip-browser-warning, Content-Type, Accept")
		
		// Обрабатываем предварительный запрос OPTIONS
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// Обработчик запросов
func apiHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Формируем ответ
	response := APIResponse{
		Message:  "ilya-denisov",
		Header:   r.Header.Get("x-test"),
		Body:     string(reqBody),
	}

	// Устанавливаем заголовок и кодируем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func main() {
	// Создаем маршрутизатор
	mux := http.NewServeMux()
	mux.HandleFunc("/result4/", apiHandler)

	// Обертываем маршрутизатор в middleware
	wrappedMux := corsMiddleware(mux)

	// Запускаем сервер
	log.Println("Starting server on :5000")
	if err := http.ListenAndServe(":5000", wrappedMux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}