package main

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db  *pgxpool.Pool
	ctx context.Context
	t   *template.Template
}

func main() {
	db, err := pgxpool.New(context.Background(), "postgresql://web:web@db:5432/web")
	if err != nil {
		slog.Error("Failed to connect to database:", err)
	}
	if _, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS todos (id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, item TEXT NOT NULL)
	`); err != nil {
		slog.Error("Failed to create tables:", err)
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		slog.Error("Failed to parse templates:", err)
	}

	ctx := context.Background()

	s := Service{
		db:  db,
		t:   tmpl,
		ctx: ctx,
	}

	http.HandleFunc("GET /style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "style.css")
	})
	http.HandleFunc("GET /", s.index)
	http.HandleFunc("POST /todo", s.add)
	http.HandleFunc("POST /todo/{id}", s.remove)

	port := 8080
	address := fmt.Sprintf(":%d", port)
	slog.Info("Starting server", "address", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}

type IndexData struct {
	Items []TodoItem
}

func (s *Service) index(w http.ResponseWriter, r *http.Request) {
	items, err := getTodos(s.db, s.ctx)
	if err != nil {
		slog.Error("Failed to get todos", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	data := IndexData{
		Items: items,
	}

	if err := s.t.ExecuteTemplate(w, "index.html", data); err != nil {
		slog.Error("Failed to execute template", "error", err)
	}
}

func (s *Service) add(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("Failed to parse form:", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	item := r.FormValue("item")

	slog.Info("adding item: %s", item)

	if err := addTodo(s.db, s.ctx, item); err != nil {
		slog.Error("Failed to add todo item:", err)
		http.Error(w, "Failed to add todo item", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Service) remove(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error("Failed to porse price:", err)
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := removeTodo(s.db, s.ctx, id); err != nil {
		slog.Error("Failed to remove item:", err)
		http.Error(w, "Failed to remove item", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type TodoItem struct {
	ID   int
	Item string
}

func getTodos(db *pgxpool.Pool, ctx context.Context) ([]TodoItem, error) {
	todoRows, err := db.Query(ctx, "SELECT id, item FROM todos ORDER BY id")
	if err != nil {
		return nil, err
	}

	defer todoRows.Close()

	var todos []TodoItem

	for todoRows.Next() {
		var todo TodoItem
		if err := todoRows.Scan(&todo.ID, &todo.Item); err != nil {
			slog.Error("Failed to scan todo item:", err)
			continue
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func addTodo(db *pgxpool.Pool, ctx context.Context, item string) error {
	_, err := db.Exec(ctx, "INSERT INTO todos (item) VALUES ($1)", item)
	return err
}

func removeTodo(db *pgxpool.Pool, ctx context.Context, id int64) error {
	_, err := db.Exec(ctx, "DELETE FROM todos WHERE id = $1", id)
	return err
}
