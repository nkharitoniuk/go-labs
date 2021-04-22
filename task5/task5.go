package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strings"
	"time"
)

type Status string

const (
	New         Status = "new"
	In_progress Status = "in_progress"
	Canceled    Status = "canceled"
	Completed   Status = "completed"
	Expired     Status = "expired"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "task5"
	dbname   = "task5"
)

type Task struct {
	ID          string    `db:"task5" json:"id"`
	Name        string    `db:"task5" json:"name"`
	Description string    `db:"task5" json:"description"`
	DueDate     time.Time `db:"task5" json:"due_date"`
	Status      Status    `db:"task5" json:"status"`
}

func task(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		getTask(w, r)
	case "PUT":
		updateTask(w, r)
	case "DELETE":
		deleteTask(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Some error"}`))
	}
}

func task_(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		getTasks(w, r)
	case "POST":
		createTask(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Some error"}`))
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := getURLIdParam(w, r)
	err := store.DeleteTask(id)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(`{"message": "Can not delete task"}`))
		return
	} else {
		w.Write([]byte(`{"message": "Task deleted"}`))
	}
}

func getURLIdParam(w http.ResponseWriter, r *http.Request) string {
	id := strings.Replace(r.URL.Path, "/task/", "", 1)
	if id == "" {
		w.Write([]byte(`{"message": "Invalid URL"}`))
	}
	return id
}

func createTask(w http.ResponseWriter, r *http.Request) {
	task := Task{}
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(`{"message": "Bad request body"}`))
		return
	}
	err = store.CreateTask(&task)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(`{"message": "Can not create task"}`))
		return
	} else {
		w.Write([]byte(`{"message": "Task created"}`))
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	ts, err := store.GetTasks()
	taskListBytes, err := json.Marshal(ts)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(taskListBytes)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	ts, err := store.GetTask(getURLIdParam(w, r))
	if ts == nil {
		w.Write([]byte(`{"message": "Task not found"}`))
		return
	}
	taskBytes, err := json.Marshal(ts)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(taskBytes)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(`{"message": "Bad request body"}`))
		return
	}
	id := getURLIdParam(w, r)
	ts, err := store.GetTask(id)
	if len(ts) == 0 {
		w.Write([]byte(`{"message": "Task not found"}`))
	} else {
		task.ID = id
		err = store.UpdateTask(&task)
		if err != nil {
			fmt.Println(err)
		} else {
			w.Write([]byte(`{"message": "Task updated"}`))
		}
	}
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()

	if err != nil {
		panic(err)
	}
	InitStore(&dbStore{db: db})

	http.HandleFunc("/tasks", task_)
	http.HandleFunc("/task/", task)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Store interface {
	CreateTask(task *Task) error
	UpdateTask(task *Task) error
	GetTasks() ([]*Task, error)
	DeleteTask(id string) error
	GetTask(id string) ([]*Task, error)
}

type dbStore struct {
	db *sql.DB
}

func (store *dbStore) DeleteTask(id string) error {
	_, err := store.db.Exec("delete from task where id = $1", id)
	if err != nil {
		panic(err)
	}
	return err
}

func (store *dbStore) CreateTask(task *Task) error {
	_, err := store.db.Query("INSERT INTO task(name, description, due_date, status) VALUES ($1,$2,$3,$4)",
		task.Name, task.Description, task.DueDate, task.Status)
	return err
}

func (store *dbStore) UpdateTask(task *Task) error {
	_, err := store.db.Query("update task set name = $1, description=$2, due_date=$3, status=$4 where id = $5",
		task.Name, task.Description, task.DueDate, task.Status, task.ID)
	return err
}

func (store *dbStore) GetTasks() ([]*Task, error) {
	rows, err := store.db.Query("SELECT id, name, description, due_date, status from task")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next() {
		task := &Task{}
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.DueDate, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (store *dbStore) GetTask(id string) ([]*Task, error) {
	rows, err := store.db.Query("SELECT id, name, description, due_date, status from task where id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []*Task{}
	for rows.Next() {
		task := &Task{}
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.DueDate, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

var store Store

func InitStore(s Store) {
	store = s
}
