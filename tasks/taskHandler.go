package task

import (
	"fmt"
	"net/http"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Реализация обработчика /api/task
func NewTaskMaker(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		postTask(w, r)
		return
	case "GET":
		getTaskByIDHandler(w, r)
		return
	case "PUT":
		updateTask(w, r)
		return
	case "DELETE":
		deleteTaskByID(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, `{"error":"метод не разрешен"}`)
		return
	}
}

// Реализация обработчика /api/tasks
func NewTaskMakerGet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTask(w)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, `{"error":"метод не разрешен"}`)
		return
	}
}

func MarkTaskAsDone(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		postTaskDone(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, `{"error":"метод не разрешен"}`)
		return
	}
}
