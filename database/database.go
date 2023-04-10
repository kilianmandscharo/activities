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

func (db *Database) GetActivity(activityId int) (schemas.Activity, error) {
	var activity schemas.Activity
	row := db.db.QueryRow("SELECT * FROM activities WHERE id = $1", activityId)
	var id int
	var name string
	var userId int
	if err := row.Scan(&id, &name, &userId); err != nil {
		return activity, err
	}
	blocks, err := db.GetBlocks(activityId)
	if err != nil {
		return activity, err
	}
	activity.Id = id
	activity.Name = name
	activity.UserId = userId
	activity.Blocks = blocks
	return activity, nil
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

func (db *Database) UpdateActivity(id int, name string) error {
	_, err := db.db.Exec("UPDATE activities SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetBlocks(activityId int) ([]schemas.Block, error) {
	var blocks []schemas.Block
	rows, err := db.db.Query("SELECT * FROM blocks WHERE activity_id = $1 AND end_time IS NOT NULL", activityId)
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

func (db *Database) GetBlock(blockId int) (schemas.Block, error) {
	var block schemas.Block
	row := db.db.QueryRow("SELECT * FROM blocks WHERE id = $1", blockId)
	var id int
	var startTime string
	var endTime string
	var activityId int
	if err := row.Scan(&id, &startTime, &endTime, &activityId); err != nil {
		return block, err
	}
	pauses, err := db.GetPauses(blockId)
	if err != nil {
		return block, err
	}

	block.Id = id
	block.StartTime = startTime
	block.EndTime = endTime
	block.ActivityId = activityId
	block.Pauses = pauses
	return block, nil
}

func (db *Database) GetCurrentBlock() (schemas.Block, error) {
	var block schemas.Block
	row := db.db.QueryRow("SELECT * FROM blocks WHERE end_time IS NULL")
	var id int
	var startTime string
	var endTime sql.NullString
	var activityId int
	err := row.Scan(&id, &startTime, &endTime, &activityId)
	if err == sql.ErrNoRows {
		return block, nil
	}
	if err != nil {
		return block, err
	}
	pauses, err := db.GetPauses(id)
	if err != nil {
		return block, err
	}

	block.Id = id
	block.StartTime = startTime
	block.EndTime = endTime.String
	block.ActivityId = activityId
	block.Pauses = pauses
	return block, nil
}

func (db *Database) AddBlock(startTime string, endTime string, activityId int) (int, error) {
	row := db.db.QueryRow(
		"INSERT INTO blocks (start_time, end_time, activity_id) VALUES ($1, $2, $3) RETURNING id",
		startTime,
		newNullString(endTime),
		activityId)
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (db *Database) UpdateBlock(id int, startTime string, endTime string) error {
	_, err := db.db.Exec("UPDATE blocks SET start_time = $1, end_time = $2 WHERE id = $3", startTime, endTime, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetPauses(blockId int) ([]schemas.Pause, error) {
	var pauses []schemas.Pause

	rows, err := db.db.Query("SELECT * FROM pauses WHERE block_id = $1", blockId)
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

func (db *Database) AddPause(startTime string, endTime string, blockId int) (int, error) {
	row := db.db.QueryRow(
		"INSERT INTO pauses (start_time, end_time, block_id) VALUES ($1, $2, $3) RETURNING id",
		startTime,
		endTime,
		blockId)
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (db *Database) UpdatePause(id int, startTime string, endTime string) error {
	_, err := db.db.Exec("UPDATE pauses SET start_time = $1, end_time = $2 WHERE id = $3", startTime, endTime, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DeletePauses(blockId int) error {
	_, err := db.db.Exec("DELETE FROM pauses WHERE block_id = $1", blockId)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DeleteByTableAndId(table string, id int) error {
	_, err := db.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = %d", table, id))
	if err != nil {
		return err
	}
	return nil
}

func createTable(db *sql.DB, name string, columns string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s %s", name, columns))
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DeleteTable(name string) error {
	_, err := db.db.Exec("DROP TABLE IF EXISTS ?", name)
	if err != nil {
		return err
	}
	return nil
}

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

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
