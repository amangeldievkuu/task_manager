package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var tasks []Task

// Handler for GET /tasks
func getTasks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getTasks(w, r)
	} else if r.Method == http.MethodPost {
		createTask(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func updateTask(w http.ResponseWriter, r *http.Request, id int) {
	var updated Task
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for i, t := range tasks {
		if t.ID == id {
			updated.ID = id
			tasks[i] = updated
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updated)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request, id int) {
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func getOneTask(w http.ResponseWriter, _ *http.Request, id int) {
	var task Task

	for i, t := range tasks {
		if t.ID == id {
			task = tasks[i]
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func taskByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid task ID ", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPut {
		updateTask(w, r, id)
		return
	} else if r.Method == http.MethodGet {
		getOneTask(w, r, id)
		return
	} else if r.Method == http.MethodDelete {
		deleteTask(w, r, id)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	http.HandleFunc("/tasks", taskHandler)
	http.HandleFunc("/tasks/", taskByIdHandler)

	fmt.Println("listens on port: 8000")
	http.ListenAndServe(":8000", nil)
}
