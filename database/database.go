package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const DbFileDefault = "scheduler.db"

func InitDatabase() error {
	dbFile := determineDbFile()
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println("Ошибка открытия базы данных:", err)
		return fmt.Errorf("ошибка открытия базы данных: %v", err)
	}
	defer db.Close()

	err = createSchedulerTable(db)

	if err != nil {
		return err
	}

	return nil
}

func determineDbFile() string {
	envDbFile := os.Getenv("TODO_DBFILE")
	if envDbFile != "" {
		return envDbFile
	}
	return DbFileDefault
}

func createSchedulerTable(db *sql.DB) error {
	createTableSQL := `
	  CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		title TEXT,
		comment TEXT,
		repeat TEXT(128)
	  );

	  CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		fmt.Println("Ошибка создания таблицы:", err)
		return fmt.Errorf("ошибка создания таблицы: %v", err)
	}

	fmt.Println("База данных успешно создана")
	return nil
}
