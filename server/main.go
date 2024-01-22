package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
)

var db *sql.DB

// Task структура представляет собой задачу в To-Do List
type Task struct {
	ID int `json:"id"`
	Title string `json:"title"`
}

type SuccesResponse struct {
	Status string `json:"status"`
}

// Инициализация базы данных SQLite
func initDB() *sql.DB {
	database, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	return database
}

// Получение списка всех задач
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Создание новой задачи
func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO tasks (title) VALUES (?)", task.Title)
	if err != nil {
		log.Fatal(err)
	}

	status := "success"
	w.Header().Set("Content-Type", "application/json")
	// Создаем объект SuccesResponse и кодируем его в JSON
	response := SuccesResponse{Status: status}
	json.NewEncoder(w).Encode(response)
}

// Обновление существующей задачи
func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// taskID := vars["id"]
	taskIDStr := vars["id"]

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Преобразование строки taskIDStr в целое число
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE tasks SET title = ? WHERE id = ?", task.Title, taskID)
	if err != nil {
		log.Fatal(err)
	}

	status := "success"
	w.Header().Set("Content-Type", "application/json")
	// Создаем объект SuccesResponse и кодируем его в JSON
	response := SuccesResponse{Status: status}
	json.NewEncoder(w).Encode(response)
}

// Удаление задачи
func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIDStr := vars["id"]

	// Преобразование строки taskIDStr в целое число
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	if err != nil {
		log.Fatal(err)
	}

	status := "success"
	w.Header().Set("Content-Type", "application/json")
	// Создаем объект SuccesResponse и кодируем его в JSON
	response := SuccesResponse{Status: status}
	json.NewEncoder(w).Encode(response)
}

func main() {
	db = initDB()
	defer db.Close()

	router := mux.NewRouter()

	// Обслуживание API
	router.HandleFunc("/tasks", getAllTasks).Methods("GET")
	router.HandleFunc("/task", createTask).Methods("POST")
	router.HandleFunc("/task/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")

	// Обслуживание статических файлов из папки "static"
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

