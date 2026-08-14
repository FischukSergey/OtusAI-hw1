package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func openStore(dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)

	if _, err := db.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable WAL: %w", err)
	}
	if _, err := db.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable foreign_keys: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := s.seedQuestions(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS questions (
  id INTEGER PRIMARY KEY,
  text TEXT NOT NULL,
  type TEXT NOT NULL,
  options_json TEXT NOT NULL DEFAULT '[]'
);

CREATE TABLE IF NOT EXISTS submissions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS answers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  submission_id INTEGER NOT NULL,
  question_id INTEGER NOT NULL,
  value TEXT NOT NULL,
  FOREIGN KEY (submission_id) REFERENCES submissions(id),
  FOREIGN KEY (question_id) REFERENCES questions(id)
);
`
	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func (s *Store) seedQuestions() error {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM questions`).Scan(&count); err != nil {
		return fmt.Errorf("count questions: %w", err)
	}
	if count > 0 {
		return nil
	}

	defaults := []Question{
		{ID: 1, Text: "Как вас зовут?", Type: "text"},
		{ID: 2, Text: "Какой ваш любимый язык программирования?", Type: "radio",
			Options: []string{"Go", "Python", "JavaScript", "Java", "Другой"}},
		{ID: 3, Text: "Какой у вас опыт в разработке?", Type: "select",
			Options: []string{"Менее 1 года", "1–3 года", "3–5 лет", "Более 5 лет"}},
		{ID: 4, Text: "Что вы хотите изучить в этом курсе?", Type: "text"},
		{ID: 5, Text: "Откуда вы узнали о курсе?", Type: "radio",
			Options: []string{"Друзья / коллеги", "Соцсети", "Поисковик", "Реклама"}},
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO questions (id, text, type, options_json) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, q := range defaults {
		opts, err := json.Marshal(q.Options)
		if err != nil {
			return err
		}
		if opts == nil {
			opts = []byte("[]")
		}
		if _, err := stmt.Exec(q.ID, q.Text, q.Type, string(opts)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) ListQuestions() ([]Question, error) {
	rows, err := s.db.Query(`SELECT id, text, type, options_json FROM questions ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Question
	for rows.Next() {
		var q Question
		var optsJSON string
		if err := rows.Scan(&q.ID, &q.Text, &q.Type, &optsJSON); err != nil {
			return nil, err
		}
		if optsJSON != "" && optsJSON != "[]" {
			if err := json.Unmarshal([]byte(optsJSON), &q.Options); err != nil {
				return nil, err
			}
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

func (s *Store) CreateSubmission(answers []Answer) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	createdAt := time.Now().UTC().Format(time.RFC3339)
	res, err := tx.Exec(`INSERT INTO submissions (created_at) VALUES (?)`, createdAt)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(`INSERT INTO answers (submission_id, question_id, value) VALUES (?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, a := range answers {
		if _, err := stmt.Exec(id, a.QuestionID, a.Value); err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

type SubmissionSummary struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Answers   []Answer  `json:"answers"`
}

func (s *Store) ListSubmissions() ([]SubmissionSummary, error) {
	rows, err := s.db.Query(`SELECT id, created_at FROM submissions ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}

	type row struct {
		id        int64
		createdAt time.Time
	}
	var headers []row
	for rows.Next() {
		var id int64
		var createdAt string
		if err := rows.Scan(&id, &createdAt); err != nil {
			rows.Close()
			return nil, err
		}
		t, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			rows.Close()
			return nil, err
		}
		headers = append(headers, row{id: id, createdAt: t})
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	out := make([]SubmissionSummary, 0, len(headers))
	for _, h := range headers {
		answers, err := s.getAnswers(h.id)
		if err != nil {
			return nil, err
		}
		out = append(out, SubmissionSummary{
			ID:        h.id,
			CreatedAt: h.createdAt,
			Answers:   answers,
		})
	}
	return out, nil
}

func (s *Store) GetSubmission(id int64) (*SubmissionSummary, error) {
	var createdAt string
	err := s.db.QueryRow(`SELECT created_at FROM submissions WHERE id = ?`, id).Scan(&createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, err
	}
	answers, err := s.getAnswers(id)
	if err != nil {
		return nil, err
	}
	return &SubmissionSummary{ID: id, CreatedAt: t, Answers: answers}, nil
}

func (s *Store) getAnswers(submissionID int64) ([]Answer, error) {
	rows, err := s.db.Query(
		`SELECT question_id, value FROM answers WHERE submission_id = ? ORDER BY question_id`,
		submissionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Answer
	for rows.Next() {
		var a Answer
		if err := rows.Scan(&a.QuestionID, &a.Value); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	if out == nil {
		out = []Answer{}
	}
	return out, rows.Err()
}

type Stats struct {
	TotalSubmissions int              `json:"total_submissions"`
	LanguageCounts   map[string]int   `json:"language_counts"`
	ExperienceCounts map[string]int   `json:"experience_counts"`
}

func (s *Store) GetStats() (*Stats, error) {
	stats := &Stats{
		LanguageCounts:   map[string]int{},
		ExperienceCounts: map[string]int{},
	}

	if err := s.db.QueryRow(`SELECT COUNT(*) FROM submissions`).Scan(&stats.TotalSubmissions); err != nil {
		return nil, err
	}

	langRows, err := s.db.Query(`
SELECT value, COUNT(*) AS cnt
FROM answers
WHERE question_id = 2
GROUP BY value
ORDER BY cnt DESC`)
	if err != nil {
		return nil, err
	}
	defer langRows.Close()
	for langRows.Next() {
		var value string
		var cnt int
		if err := langRows.Scan(&value, &cnt); err != nil {
			return nil, err
		}
		stats.LanguageCounts[value] = cnt
	}
	if err := langRows.Err(); err != nil {
		return nil, err
	}

	expRows, err := s.db.Query(`
SELECT value, COUNT(*) AS cnt
FROM answers
WHERE question_id = 3
GROUP BY value
ORDER BY cnt DESC`)
	if err != nil {
		return nil, err
	}
	defer expRows.Close()
	for expRows.Next() {
		var value string
		var cnt int
		if err := expRows.Scan(&value, &cnt); err != nil {
			return nil, err
		}
		stats.ExperienceCounts[value] = cnt
	}
	return stats, expRows.Err()
}
