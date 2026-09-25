package main

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
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
		slog.Error("Failed to connect to database", "error", err)
	}
	if _, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS todos 
			(
				id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
				item TEXT NOT NULL
			)
	`); err != nil {
		slog.Error("Failed to create tables", "error", err)
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		slog.Error("Failed to parse templates", "error", err)
	}

	ctx := context.Background()

	s := Service{
		db:  db,
		t:   tmpl,
		ctx: ctx,
	}

	// Set up routes
	http.HandleFunc("GET /style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "style.css")
	})
	http.HandleFunc("GET /", s.index)
	http.HandleFunc("POST /todo", s.add)
	http.HandleFunc("POST /todo/{id}", s.remove)

	port := 8080
	address := fmt.Sprintf(":%d", port)
	slog.Info(fmt.Sprintf("Starting server on: http://localhost:%d", port))
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}

// Function called on request to `/`
func (s *Service) index(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// Function called on request to `/todo`
func (s *Service) add(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// Function called on request to `/todo/{id}`
func (s *Service) remove(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// A todo item (the values stored in the database)
type TodoItem struct {
	ID   int
	Item string
}

// Get the list of todos
func getTodos(db *pgxpool.Pool, ctx context.Context) ([]TodoItem, error) {
	todoRows, err := db.Query(ctx, "") // TODO
	if err != nil {
		return nil, err
	}

	// Close the rows when the function returns
	defer todoRows.Close()

	var todos []TodoItem

	// Convert the returned rows to a list of TodoItems
	for todoRows.Next() {
		var todo TodoItem
		if err := todoRows.Scan(&todo.ID, &todo.Item); err != nil {
			slog.Error("Failed to scan todo item", "error", err)
			continue
		}
		todos = append(todos, todo)
	}


	return todos, nil
}

// Add an item to the list
func addTodo(db *pgxpool.Pool, ctx context.Context, item string) error {
	_, err := db.Exec(ctx, "", item) // TODO
	return err
}

// Remove an item from the list
func removeTodo(db *pgxpool.Pool, ctx context.Context, id int64) error {
	_, err := db.Exec(ctx, "", id) // TODO
	return err
}
