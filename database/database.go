package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/kilianmandscharo/activities/schemas"
)

var tables = []schemas.TableSchema {
	{Name: "users", Columns: "(id serial PRIMARY KEY, name text, email text, password text)"},
	{Name: "activities", Columns: "(id serial PRIMARY KEY, name text, user_id int references users(id) ON DELETE CASCADE)"},
	{Name: "blocks", Columns:  "(id serial PRIMARY KEY, start_time timestamp, end_time timestamp, activity_id int references activities(id) ON DELETE CASCADE)"},
	{Name: "pauses", Columns: "(id serial PRIMARY KEY, start_time timestamp, end_time timestamp, block_id int references blocks(id) ON DELETE CASCADE)"},
}

func databaseError(message string, err error) error {
	errorMessage := fmt.Sprintf("%s, Error: %s", message, err)
	return errors.New(errorMessage)
}

func InitTables(db *sql.DB) error {
	for _, table := range tables {
		err := createTable(db, table.Name, table.Columns)	
		if err != nil {
			return err	
		}
	}

	return nil
}

func ClearTables(db *sql.DB) error {
	for _, table := range tables {
		err := clearTable(db, table.Name)
		if err != nil {
			return err	
		}
	}

	return nil
} 

func DeleteTables(db *sql.DB) error {
	var reverseTableNames []string
	for _, table := range tables {
		reverseTableNames = append(reverseTableNames, table.Name)
	}

	for i, j := 0, len(reverseTableNames)-1; i < j; i, j = i+1, j-1 {
		reverseTableNames[i], reverseTableNames[j] = reverseTableNames[j], reverseTableNames[i]
	}

	for _, name := range reverseTableNames {
		err := deleteTable(db, name)
		if err != nil {
			return err	
		}
	}

	return nil
}

func createTable(db *sql.DB, name string, columns string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s %s", name, columns)
	_, err := db.Exec(query)
	if err != nil {
		return databaseError("Could not create table", err) 
	}

	return nil
}

func deleteTable(db *sql.DB, name string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
	_, err := db.Exec(query)
	if err != nil {
		return databaseError("Could not delete table", err)
	}

	return nil
}

func clearTable(db *sql.DB, name string) error {
	query := fmt.Sprintf("DELETE FROM %s", name)
	_, err := db.Exec(query)
	if err != nil {
		return databaseError("Could not clear table", err)
	}

	return nil
}

func AddUser(db *sql.DB, name string, email string, password string) error {
	_, err := db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", name, email, password)
	if err != nil {
		return databaseError("Could not add the user to the database", err) 
	}

	return nil
}

func AddActivity(db *sql.DB, name string, user_id int) error {
	_, err := db.Exec("INSERT INTO activities (name, user_id) VALUES ($1, $2)", name, user_id)
	if err != nil {
		return databaseError("Could not add the activity to the database", err) 
	}
	
	return nil
}

func AddBlock(db *sql.DB, start_time string, end_time string, activity_id int) error {
	_, err := db.Exec("INSERT INTO blocks (start_time, end_time, activity_id) VALUES ($1, $2, $3)", start_time, end_time, activity_id)
	if err != nil {
		return databaseError("Could not add the block to the database", err) 
	}

	return nil
}

func AddPause(db *sql.DB, start_time string, end_time string, block_id int) error {
	_, err := db.Exec("INSERT INTO pauses (start_time, end_time, block_id) VALUES ($1, $2, $3)", start_time, end_time, block_id)
	if err != nil {
		return databaseError("Could not add the pause to the database", err) 
	}

	return nil
}

func DeleteByTableAndId(db *sql.DB, table string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = %d", table, id)
	_, err := db.Exec(query)
	if err != nil {
		return databaseError("Could not delete the row from the database", err)
	}

	return nil
}

func GetAllActivities(db *sql.DB, userId int) []schemas.Activity {
	var activities []schemas.Activity

	rows, err := db.Query("SELECT * FROM activities WHERE user_id = $1", userId)	
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var (
			id int
			name string
			userId int
		)		

		if err := rows.Scan(&id, &name, &userId); err != nil {
			log.Fatal(err)
		}

		var activity schemas.Activity

		activity.Id = id
		activity.Name = name
		activity.UserId = userId
		activity.Blocks = GetAllBlocks(db, id)

		activities = append(activities, activity)
	}

	return activities
}

func GetAllBlocks(db *sql.DB, activityId int) []schemas.Block {
	var blocks []schemas.Block

	rows, err := db.Query("SELECT * FROM blocks WHERE activity_id = $1", activityId)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var (
			id int
			startTime string
			endTime string
			activityId int
		)

		if err := rows.Scan(&id, &startTime, &endTime, &activityId); err != nil {
			log.Fatal(err)
		}

		var block schemas.Block

		block.Id = id
		block.StartTime = startTime
		block.EndTime = endTime
		block.ActivityId = activityId
		block.Pauses = GetAllPauses(db, id)

		blocks = append(blocks, block)
	}

	return blocks
}

func GetAllPauses(db *sql.DB, blockId int) []schemas.Pause {
	var pauses []schemas.Pause

	rows, err := db.Query("SELECT * FROM pauses WHERE id = $1", blockId)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var (
			id int
			startTime string
			endTime string
			blockId int
		)

		if err := rows.Scan(&id, &startTime, &endTime, &blockId); err != nil {
			log.Fatal(err)	
		}

		var pause schemas.Pause

		pause.Id = id
		pause.StartTime = startTime
		pause.EndTime = endTime
		pause.BlockId = blockId

		pauses = append(pauses, pause)
	}

	return pauses
}
