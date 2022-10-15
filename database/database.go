package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type tableSchema struct {
	name string
	columns string
} 

func databaseError(message string, err error) error {
	errorMessage := fmt.Sprintf("%s, Error: %s", message, err)
	return errors.New(errorMessage)
}

func defineTables() []tableSchema {
	var tables []tableSchema

	tables = append(tables, tableSchema{"users", "(id serial PRIMARY KEY, name text, email text, password text)"})
	tables = append(tables, tableSchema{"activities", "(id serial PRIMARY KEY, name text, user_id int references users(id))"})
	tables = append(tables, tableSchema{"blocks", "(id serial PRIMARY KEY, start_time timestamp, end_time timestamp, activity_id int references activities(id))"})
	tables = append(tables, tableSchema{"pauses", "(id serial PRIMARY KEY, start_time timestamp, end_time timestamp, block_id int references blocks(id))"})

	return tables
}

func InitDatabaseTables(db *sql.DB) error {
	tables := defineTables()

	for _, table := range tables {
		err := createTable(db, table.name, table.columns)	
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func ClearDatabaseTables(db *sql.DB) error {
	tables := defineTables()

	var reverseTableNames []string
	for _, table := range tables {
		reverseTableNames = append(reverseTableNames, table.name)
	}

	for i, j := 0, len(reverseTableNames)-1; i < j; i, j = i+1, j-1 {
		reverseTableNames[i], reverseTableNames[j] = reverseTableNames[j], reverseTableNames[i]
	}

	for _, name := range reverseTableNames {
		err := deleteTable(db, name)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func createTable(db *sql.DB, name string, columns string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s %s", name, columns)
	_, err := db.Exec(query)
	if err != nil {
		return databaseError("Failed to create the database table", err) 
	}

	return nil
}

func deleteTable(db *sql.DB, name string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
	_, err := db.Exec(query)
	if err != nil {
		return databaseError("Failed to delete the database table", err)
	}

	return nil
}

func AddUser(db *sql.DB, name string, email string, password string) error {
	_, err := db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", name, email, password)
	if err != nil {
		return databaseError("Failed to add the user to the database", err) 
	}

	return nil
}

func AddActivity(db *sql.DB, name string, user_id int) error {
	_, err := db.Exec("INSERT INTO activities (name, user_id) VALUES ($1, $2)", name, user_id)
	if err != nil {
		return databaseError("Failed to add the activity to the database", err) 
	}
	
	return nil
}

func AddBlock(db *sql.DB, start_time string, end_time string, activity_id int) error {
	_, err := db.Exec("INSERT INTO blocks (start_time, end_time, activity_id) VALUES ($1, $2, $3)", start_time, end_time, activity_id)
	if err != nil {
		return databaseError("Failed to add the block to the database", err) 
	}

	return nil
}

func AddPause(db *sql.DB, start_time string, end_time string, block_id int) error {
	_, err := db.Exec("INSERT INTO pauses (start_time, end_time, block_id) VALUES ($1, $2, $3)", start_time, end_time, block_id)
	if err != nil {
		return databaseError("Failed to add the pause to the database", err) 
	}

	return nil
}

func GetUsers(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users") 
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var (
			id int
			name string
			email string
			password string
		)
		if err := rows.Scan(&id, &name, &email, &password); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, email, password)
	}
}
