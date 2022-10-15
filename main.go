package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(("Could not load .env file"))
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PW"), os.Getenv("DB_NAME")) 
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable(db, "users", "(id serial, name text, password text, email text)")
	createTable(db, "activities", "(id serial, name text, user_id int)")
	createTable(db, "blocks", "(id serial, start_time timestamp, end_time timestamp, activity_id int)")
	createTable(db, "pauses", "(id serial, start_time timestamp, end_time timestamp, block_id int)")

	addUser(db, "Kilian", "kilian187@gmail.com", "6fd7gf9dfgh90h900")
	addUser(db, "Peter", "peter291@gmail.com", "34ffvfdghash9vcbv")

	addActivity(db, "Work", 1)
	addActivity(db, "Guitar Practice", 1)
	addActivity(db, "Gymnastics", 2)

	addBlock(db, "2022-10-15 15:30:00", "2022-10-15 17:30", 1)

	addPause(db, "2022-10-15 16:00", "2022-10-15 16:15", 1)
}

func createTable(db *sql.DB, name string, columns string) {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s %s", name, columns)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func addUser(db *sql.DB, name string, email string, password string) {
	query := fmt.Sprintf("INSERT INTO users (name, email, password) VALUES ('%s', '%s', '%s')", name, email, password)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func addActivity(db *sql.DB, name string, user_id int) {
	query := fmt.Sprintf("INSERT INTO activities (name, user_id) VALUES ('%s', %d)", name, user_id)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func addBlock(db *sql.DB, start_time string, end_time string, activity_id int) {
	query := fmt.Sprintf("INSERT INTO blocks (start_time, end_time, activity_id) VALUES ('%s', '%s', %d)", start_time, end_time, activity_id)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func addPause(db *sql.DB, start_time string, end_time string, block_id int) {
	query := fmt.Sprintf("INSERT INTO pauses (start_time, end_time, block_id) VALUES ('%s', '%s', %d)", start_time, end_time, block_id)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
