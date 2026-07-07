package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
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

type Submission struct {
	Answers   []Answer  `json:"answers"`
	CreatedAt time.Time `json:"created_at"`
}

var questions = []Question{
	{ID: 1, Text: "Как вас зовут?", Type: "text"},
	{ID: 2, Text: "Какой ваш любимый язык программирования?", Type: "radio",
		Options: []string{"Go", "Python", "JavaScript", "Java", "Другой"}},
	{ID: 3, Text: "Какой у вас опыт в разработке?", Type: "select",
		Options: []string{"Менее 1 года", "1–3 года", "3–5 лет", "Более 5 лет"}},
	{ID: 4, Text: "Что вы хотите изучить в этом курсе?", Type: "text"},
	{ID: 5, Text: "Откуда вы узнали о курсе?", Type: "radio",
		Options: []string{"Друзья / коллеги", "Соцсети", "Поисковик", "Реклама"}},
}

var (
	mu          sync.Mutex
	submissions []Submission
)

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

func handleQuestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func handleAnswers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var sub Submission
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	sub.CreatedAt = time.Now()

	mu.Lock()
	submissions = append(submissions, sub)
	total := len(submissions)
	mu.Unlock()

	log.Printf("Получен ответ #%d: %+v", total, sub.Answers)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
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

	mux := http.NewServeMux()

	mux.HandleFunc("/questions", corsMiddleware(handleQuestions))
	mux.HandleFunc("/answers", corsMiddleware(handleAnswers))

	fs := http.FileServer(http.Dir(frontendDir))
	mux.Handle("/", fs)

	log.Printf("Фронтенд раздаётся из: %s", frontendDir)
	log.Printf("Сервер запущен на http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
