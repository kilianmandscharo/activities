package database

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var db *Database

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("could not load .env file")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PW"),
		os.Getenv("DB_NAME"))

	database, err := New(connStr)
	if err != nil {
		log.Fatal("could not open database")
	}
	db = database
	if err := db.Init(); err != nil {
		log.Fatal("could not initialize database")
	}
	if err := db.Clear(); err != nil {
		log.Fatal("could not clear database")
	}
}

func TestAddUser(t *testing.T) {
	testName := "Apollo"
	testEmail := "test@gmail.com"
	testPassword := "12345"
	testId, err := db.AddUser(testName, testEmail, testPassword)
	if err != nil {
		t.Fatalf("could not add user, %v", err)
	}
	row := db.db.QueryRow("SELECT * FROM users WHERE id = $1", testId)
	var (
		id       int
		name     string
		email    string
		password string
	)
	if err := row.Scan(&id, &name, &email, &password); err != nil {
		t.Fatalf("could not retrieve user, %v", err)
	}
	if id != testId {
		t.Fatalf("id = %q, want %q", id, testId)
	}
	if name != testName {
		t.Fatalf("name = %q, want %q", name, testName)
	}
	if email != testEmail {
		t.Fatalf("email = %q, want %q", email, testEmail)
	}
	if password != testPassword {
		t.Fatalf("password = %q, want %q", password, testPassword)
	}
}

func TestAddActivity(t *testing.T) {
	testName := "Running"
	testUserId := 1
	testId, err := db.AddActivity(testName, testUserId)
	if err != nil {
		t.Fatalf("could not add activity, %v", err)
	}
	row := db.db.QueryRow("SELECT * FROM activities WHERE id = $1", testId)
	var (
		id     int
		name   string
		userId int
	)
	if err := row.Scan(&id, &name, &userId); err != nil {
		t.Fatalf("could not retrieve activity, %v", err)
	}
	if id != testId {
		t.Fatalf("id = %q, want %q", id, testId)
	}
	if name != testName {
		t.Fatalf("name = %q, want %q", name, testName)
	}
	if userId != testUserId {
		t.Fatalf("userId = %q, want %q", userId, testUserId)
	}
}

func TestAddBlock(t *testing.T) {
	testStartTime := "2023-02-01T14:00:00Z"
	testEndTime := "2023-02-01T14:30:00Z"
	testActivityId := 1
	testId, err := db.AddBlock(testStartTime, testEndTime, testActivityId)
	if err != nil {
		t.Fatalf("could not add block, %v", err)
	}
	row := db.db.QueryRow("SELECT * FROM blocks WHERE id = $1", testId)
	var (
		id         int
		startTime  string
		endTime    string
		activityId int
	)
	if err := row.Scan(&id, &startTime, &endTime, &activityId); err != nil {
		t.Fatalf("could not retrieve block, %v", err)
	}
	if id != testId {
		t.Fatalf("id = %q, want %q", id, testId)
	}
	if startTime != testStartTime {
		t.Fatalf("startTime = %q, want %q", startTime, testStartTime)
	}
	if endTime != testEndTime {
		t.Fatalf("endTime = %q, want %q", endTime, testEndTime)
	}
	if activityId != testActivityId {
		t.Fatalf("activityId = %q, want %q", activityId, testActivityId)
	}
}

func TestAddPause(t *testing.T) {
	testStartTime := "2023-02-01T14:15:00Z"
	testEndTime := "2023-02-01T14:20:00Z"
	testBlockId := 1
	testId, err := db.AddPause(testStartTime, testEndTime, testBlockId)
	if err != nil {
		t.Fatalf("could not add pause, %v", err)
	}
	row := db.db.QueryRow("SELECT * FROM pauses WHERE id = $1", testId)
	var (
		id        int
		startTime string
		endTime   string
		blockId   int
	)
	if err := row.Scan(&id, &startTime, &endTime, &blockId); err != nil {
		t.Fatalf("could not retrieve pause, %v", err)
	}
	if id != testId {
		t.Fatalf("id = %q, want %q", id, testId)
	}
	if startTime != testStartTime {
		t.Fatalf("startTime = %q, want %q", startTime, testStartTime)
	}
	if endTime != testEndTime {
		t.Fatalf("endTime = %q, want %q", endTime, testEndTime)
	}
	if blockId != testBlockId {
		t.Fatalf("blockId = %q, want %q", blockId, testBlockId)
	}
}

func TestClose(t *testing.T) {
  db.Close()
}
