package database

import (
	"database/sql"
)

type Task struct {
	ID      string
	Date    string
	Title   string
	Comment string
	Repeat  string
}

func InsertTask(date string, title string, comment string, repeat string) (int, error) {
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
