package task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jedyEvgeny/YPGoFinalJob/database"
	dCalc "github.com/jedyEvgeny/YPGoFinalJob/datecalculator"
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, `{"error":"метод не разрешен"}`)
		return
	}
}

// Реализация обработчика /api/task
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

func postTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"ошибка десериализации JSON"}`)
		return
	}

	// Проверка обязательного поля title
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"не указан заголовок задачи"}`)
		return
	}

	// Проверка и преобразование даты задачи
	date := task.Date
	_, err = time.Parse("20060102", date)
	if err != nil && date != "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"дата JSON в неверном формате. Используйте формат 20060102"}`)
		return
	}

	currentDate := time.Now()
	currentDateFormat := currentDate.Format("20060102")
	if date == "" || date < currentDateFormat && task.Repeat == "" {
		date = currentDateFormat
	}
	if date < currentDateFormat {
		date, err = dCalc.NextDate(currentDate, date, task.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, `{"error":"не удалось вычислить ближайшее событие"}`)
			return
		}
	}

	// Сохранение информации в базу данных SQLite
	taskID, err := database.InsertTask(date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error":"ошибка сохранения задачи в базу данных"}`)
		return
	}

	response := fmt.Sprintf(`{"id":"%d"}`, taskID)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, response)
}

func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")

	// Проверяем, был ли передан идентификатор задачи
	if taskID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Не указан идентификатор"}`))
		return
	}

	// Получаем задачу из базы данных по её идентификатору
	t, err := database.GetTaskByID(taskID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Задача не найдена"}`))
		return
	}

	// Кодируем задачу в формат JSON и отправляем её клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func getTask(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")

	tasks, err := database.GetAllTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error":"ошибка при получении списка задач из базы данных"}`)
		return
	}

	// Формирование JSON-ответа
	var response struct {
		TasksSlice []Task `json:"tasks"`
	}

	for _, task := range tasks {
		response.TasksSlice = append(response.TasksSlice, Task{
			ID:      task.ID,
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	if len(response.TasksSlice) == 0 {
		response.TasksSlice = []Task{}
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error":"ошибка при формировании JSON-ответа"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(jsonResponse))
}
