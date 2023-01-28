package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kilianmandscharo/activities/schemas"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

var tables = []schemas.TableSchema{
	{Name: "users", Columns: "(id serial PRIMARY KEY, name text, email text, password text)"},
	{Name: "activities", Columns: "(id serial PRIMARY KEY, name text, user_id int references users(id) ON DELETE CASCADE)"},
	{Name: "blocks", Columns: "(id serial PRIMARY KEY, start_time timestamp, end_time timestamp, activity_id int references activities(id) ON DELETE CASCADE)"},
	{Name: "pauses", Columns: "(id serial PRIMARY KEY, start_time timestamp, end_time timestamp, block_id int references blocks(id) ON DELETE CASCADE)"},
}

func New(connStr string) (*Database, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (db *Database) Init() error {
	for _, table := range tables {
		err := createTable(db.db, table.Name, table.Columns)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) Close() {
	db.db.Close()
}

func (db *Database) Clear() error {
	_, err := db.db.Exec("TRUNCATE users RESTART IDENTITY CASCADE")
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) AddUser(name string, email string, password string) (int, error) {
	row := db.db.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		name,
		email,
		password)
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (db *Database) AddActivity(name string, user_id int) (int, error) {
	row := db.db.QueryRow(
		"INSERT INTO activities (name, user_id) VALUES ($1, $2) RETURNING id",
		name,
		user_id)
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (db *Database) AddBlock(start_time string, end_time string, activity_id int) (int, error) {
	row := db.db.QueryRow(
		"INSERT INTO blocks (start_time, end_time, activity_id) VALUES ($1, $2, $3) RETURNING id",
		start_time,
		end_time,
		activity_id)
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (db *Database) AddPause(start_time string, end_time string, block_id int) (int, error) {
	row := db.db.QueryRow(
		"INSERT INTO pauses (start_time, end_time, block_id) VALUES ($1, $2, $3) RETURNING id",
		start_time,
		end_time,
		block_id)
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (db *Database) DeleteByTableAndId(table string, id int) error {
	_, err := db.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = %d", table, id))
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetActivities(userId int) ([]schemas.Activity, error) {
	var activities []schemas.Activity

	rows, err := db.db.Query("SELECT * FROM activities WHERE user_id = $1", userId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id     int
			name   string
			userId int
		)
		if err := rows.Scan(&id, &name, &userId); err != nil {
			return nil, err
		}
		blocks, err := db.GetBlocks(id)
		if err != nil {
			return nil, err
		}
		activities = append(activities, schemas.Activity{
			Id:     id,
			Name:   name,
			UserId: userId,
			Blocks: blocks})
	}
	return activities, nil
}

func (db *Database) GetBlocks(activityId int) ([]schemas.Block, error) {
	var blocks []schemas.Block

	rows, err := db.db.Query("SELECT * FROM blocks WHERE activity_id = $1", activityId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id         int
			startTime  string
			endTime    string
			activityId int
		)
		if err := rows.Scan(&id, &startTime, &endTime, &activityId); err != nil {
			return nil, err
		}
		pauses, err := db.GetPauses(id)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, schemas.Block{
			Id:         id,
			StartTime:  startTime,
			EndTime:    endTime,
			ActivityId: activityId,
			Pauses:     pauses})
	}
	return blocks, nil
}

func (db *Database) GetPauses(blockId int) ([]schemas.Pause, error) {
	var pauses []schemas.Pause

	rows, err := db.db.Query("SELECT * FROM pauses WHERE id = $1", blockId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int
			startTime string
			endTime   string
			blockId   int
		)
		if err := rows.Scan(&id, &startTime, &endTime, &blockId); err != nil {
			return nil, err
		}
		pauses = append(pauses, schemas.Pause{
			Id:        id,
			StartTime: startTime,
			EndTime:   endTime,
			BlockId:   blockId})
	}
	return pauses, nil
}

func createTable(db *sql.DB, name string, columns string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s %s", name, columns))
	if err != nil {
		return err
	}
	return nil
}

func deleteTable(db *sql.DB, name string) error {
	_, err := db.Exec("DROP TABLE IF EXISTS ?", name)
	if err != nil {
		return err
	}
	return nil
}

// func updateActivityName(db *sql.DB, activityId int, newName string) error {
// 	query := fmt.Sprintf("UPDATE activities SET name = %s WHERE id = %d", newName, activityId)
// 	_, err := db.Exec(query)
// 	if err != nil {
// 		return databaseError("Could not delete table", err)
// 	}
//
// 	return nil
// }

// func DeleteTables(db *sql.DB) error {
// 	var reverseTableNames []string
// 	for _, table := range tables {
// 		reverseTableNames = append(reverseTableNames, table.Name)
// 	}
//
// 	for i, j := 0, len(reverseTableNames)-1; i < j; i, j = i+1, j-1 {
// 		reverseTableNames[i], reverseTableNames[j] = reverseTableNames[j], reverseTableNames[i]
// 	}
//
// 	for _, name := range reverseTableNames {
// 		err := deleteTable(db, name)
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }

// func databaseError(message string, err error) error {
// 	errorMessage := fmt.Sprintf("%s, Error: %s", message, err)
// 	return errors.New(errorMessage)
// }

// func clearTable(db *sql.DB, name string) error {
// 	query := fmt.Sprintf("DELETE FROM %s", name)
// 	_, err := db.Exec(query)
// 	if err != nil {
// 		return databaseError("Could not clear table", err)
// 	}
//
// 	return nil
// }
