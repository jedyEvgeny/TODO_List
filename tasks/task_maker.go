package task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jedyEvgeny/YPGoFinalJob/database"
	dCalc "github.com/jedyEvgeny/YPGoFinalJob/datecalculator"
)

func postTask(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Add("Content-Type", "application/json")
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

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"ошибка десериализации JSON"}`)
		return
	}

	// Проверка обязательного поля id
	if task.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"не указан идентификатор задачи"}`)
		return
	}

	_, err = database.GetTaskByID(task.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"указанный идентификатор задачи не найден в базе данных"}`)
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

	currentDateFormat := time.Now().Format("20060102")
	if date == "" || date < currentDateFormat && task.Repeat == "" {
		date = currentDateFormat
	}
	if date < currentDateFormat {
		date, err = dCalc.NextDate(time.Now(), date, task.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, `{"error":"не удалось вычислить ближайшее событие"}`)
			return
		}
	}

	// Проверка обязательного поля title
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"не указан заголовок задачи"}`)
		return
	}

	// Обновление информации в базе данных SQLite
	err = database.UpdateTaskByID(task.ID, date, task.Title, task.Comment, task.Repeat)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"error":"Задача не найдена"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{}`)
}

func postTaskDone(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	task, err := database.GetTaskByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Задача не найдена"}`))
		return
	}

	if task.Repeat == "" {
		err = database.DeleteTask(task.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error":"ошибка при удалении задачи: %s"}`, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{}`)
		return
	}
	currentDate := time.Now()
	date, err := dCalc.NextDate(currentDate, task.Date, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"не удалось вычислить ближайшее событие"}`)
		return
	}
	err = database.UpdateDayForTask(id, date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"не удалось обновить дату задачи"}`)
		return
	}
	fmt.Fprintln(w, `{}`)
}

func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "отсутствует ID"}`))
		return
	}
	if _, err := strconv.Atoi(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"идентефикатор задачи должен быть числом"}`)
		return
	}
	err := database.DeleteTask(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "ошибка при удалении задачи"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{}`)
}
