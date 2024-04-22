// /
// 	database
// 		-database.go
// 	datecalculator
// 		-datecalculator.go
// 	main.go
// 	go.mod

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jedyEvgeny/YPGoFinalJob/database"
	dCalc "github.com/jedyEvgeny/YPGoFinalJob/datecalculator"
	task "github.com/jedyEvgeny/YPGoFinalJob/tasks"
)

const Port = ":7540"
const WebDir = "./web"

func main() {
	database.InitDatabase()
	port := determinePort()
	http.Handle("/", http.FileServer(http.Dir(WebDir)))
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", task.NewTaskMaker)
	http.HandleFunc("/api/tasks", task.NewTaskMakerGet)
	fmt.Println("Сервер запущен на порту", port)
	http.ListenAndServe(port, nil)
}

func determinePort() string {
	envPort := os.Getenv("TODO_PORT")
	if envPort != "" {
		return envPort
	}
	return Port
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Ошибка преобразования параметра now", http.StatusBadRequest)
		return
	}

	nextDate, err := dCalc.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, nextDate)
}
