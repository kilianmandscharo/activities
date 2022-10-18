package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/kilianmandscharo/activities/database"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(("Could not load .env file"))
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PW"), os.Getenv("DB_NAME")) 
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.DeleteTables(db)
	database.InitTables(db)

	database.AddUser(db, "Kilian", "kilian187@gmail.com", "6fd7gf9dfgh90h900")
	database.AddUser(db, "Peter", "peter291@gmail.com", "34ffvfdghash9vcbv")

	database.AddActivity(db, "Work", 1)
	database.AddActivity(db, "Guitar Practice", 1)
	database.AddActivity(db, "Gymnastics", 2)

	database.AddBlock(db, "2022-10-15 15:30:00", "2022-10-15 17:30", 1)
	database.AddBlock(db, "2022-10-15 15:30:00", "2022-10-15 17:30", 2)
	database.AddBlock(db, "2022-10-15 15:30:00", "2022-10-15 17:30", 3)

	database.AddPause(db, "2022-10-15 16:00", "2022-10-15 16:15", 1)
	database.AddPause(db, "2022-10-15 16:00", "2022-10-15 16:15", 2)
	database.AddPause(db, "2022-10-15 16:00", "2022-10-15 16:15", 3)

	database.DeleteByTableAndId(db, "users", 1)
}

