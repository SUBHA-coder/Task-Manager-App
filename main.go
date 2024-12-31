package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Priority string `json:"priority"`
	DueDate  string `json:"dueDate"`
	Complete bool   `json:"complete"`
}

var db *sql.DB

func main() {
	// Initialize the database
	initDB()
	defer db.Close()

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	r.HandleFunc("/tasks/filter/{priority}", filterTasksByPriority).Methods("GET")

	// Serve static files and templates
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", serveHome)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./db/tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tasks table if not exists
	createTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		priority TEXT NOT NULL,
		due_date TEXT NOT NULL,
		complete BOOLEAN NOT NULL
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, priority, due_date, complete FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Priority, &task.DueDate, &task.Complete)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println("Error decoding task:", err)
		return
	}

	log.Println("Received task:", task)

	// Set default completion status to false
	task.Complete = false

	// Insert into database
	result, err := db.Exec("INSERT INTO tasks (title, priority, due_date, complete) VALUES (?, ?, ?, ?)", task.Title, task.Priority, task.DueDate, task.Complete)
	if err != nil {
		http.Error(w, "Failed to insert task", http.StatusInternalServerError)
		log.Println("Error inserting task:", err)
		return
	}

	// Return the created task with ID
	id, _ := result.LastInsertId()
	task.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}

	// Update task in database
	_, err = db.Exec("UPDATE tasks SET title = ?, priority = ?, due_date = ?, complete = ? WHERE id = ?", task.Title, task.Priority, task.DueDate, task.Complete, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Delete task from database
	_, err = db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func filterTasksByPriority(w http.ResponseWriter, r *http.Request) {
	priority := mux.Vars(r)["priority"]

	rows, err := db.Query("SELECT id, title, priority, due_date, complete FROM tasks WHERE priority = ?", priority)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Priority, &task.DueDate, &task.Complete)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/index.html")
}
