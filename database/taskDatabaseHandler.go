package database

import (
	"database/sql"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func InsertTask(date, title, comment, repeat string) (int, error) {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	insertSQL := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(insertSQL, date, title, comment, repeat)
	if err != nil {
		return 0, err
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(taskID), nil
}

func GetAllTasks() ([]Task, error) {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTaskByID(taskID string) (*Task, error) {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, taskID)

	var task Task
	err = row.Scan(&task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	task.ID = taskID
	return &task, nil
}

func UpdateTaskByID(id, date, title, comment, repeat string) error {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	updateSQL := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err = db.Exec(updateSQL, date, title, comment, repeat, id)
	if err != nil {
		return err
	}

	return nil
}

// Проверка наличия задачи с указанным ID в базе данных
func TaskExists(taskID string) bool {
	dbFile := determineDbFile()
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return false
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(id) FROM scheduler WHERE id = ?", taskID).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func UpdateDayForTask(id, date string) error {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	updateSQL := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err = db.Exec(updateSQL, date, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(id string) error {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	deleteSQL := `DELETE FROM scheduler WHERE id = ?`
	_, err = db.Exec(deleteSQL, id)
	if err != nil {
		return err
	}

	return nil
}
