package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Type    string   `json:"type"` // "text", "radio", "select"
	Options []string `json:"options,omitempty"`
}

type Answer struct {
	QuestionID int    `json:"question_id"`
	Value      string `json:"value"`
}

type SubmissionRequest struct {
	Answers []Answer `json:"answers"`
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		frontendDir = "../frontend"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/survey.db"
	}

	store, err := openStore(dbPath)
	if err != nil {
		log.Fatalf("init db: %v", err)
	}
	defer store.Close()
	log.Printf("SQLite: %s", dbPath)

	mux := http.NewServeMux()

	mux.HandleFunc("/questions", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		questions, err := store.ListQuestions()
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, questions)
	}))

	mux.HandleFunc("/answers", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req SubmissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if len(req.Answers) == 0 {
			http.Error(w, "answers required", http.StatusBadRequest)
			return
		}

		id, err := store.CreateSubmission(req.Answers)
		if err != nil {
			log.Printf("create submission error: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		log.Printf("Получен ответ #%d: %+v", id, req.Answers)
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "id": id})
	}))

	mux.HandleFunc("/submissions", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		items, err := store.ListSubmissions()
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, items)
	}))

	mux.HandleFunc("/submissions/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/submissions/")
		idStr = strings.Trim(idStr, "/")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		item, err := store.GetSubmission(id)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if item == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, item)
	}))

	mux.HandleFunc("/stats", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		stats, err := store.GetStats()
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, stats)
	}))

	fs := http.FileServer(http.Dir(frontendDir))
	mux.Handle("/", fs)

	log.Printf("Фронтенд раздаётся из: %s", frontendDir)
	log.Printf("Сервер запущен на http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
