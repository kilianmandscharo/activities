package database

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var db *Database

const (
	testUserId       = 1
	testUserName     = "Apollo"
	testUserEmail    = "test@gmail.com"
	testUserPassword = "12345"

	testActivityId          = 1
	testActivityName        = "Running"
	testActivityNameUpdated = "Swimming"

	testBlockId               = 1
	testBlockStartTime        = "2023-02-01T14:00:00Z"
	testBlockEndTime          = "2023-02-01T14:30:00Z"
	testBlockStartTimeUpdated = "2023-04-05T16:00:00Z"
	testBlockEndTimeUpdated   = "2023-04-05T16:30:00Z"

	testPauseId               = 1
	testPauseStartTime        = "2023-02-01T14:15:00Z"
	testPauseEndTime          = "2023-02-01T14:20:00Z"
	testPauseStartTimeUpdated = "2023-04-05T16:15:00Z"
	testPauseEndTimeUpdated   = "2023-04-05T16:20:00Z"

	testStartTimeCurrentBlock = "2023-02-01T14:15:00Z"
	testEndTimeCurrentBlock   = ""
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("could not load .env file", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PW"),
		os.Getenv("DB_NAME_TEST"))

	database, err := New(connStr)
	if err != nil {
		log.Fatal("could not open database:", err)
	}
	db = database
	if err := db.Init(); err != nil {
		log.Fatal("could not initialize database:", err)
	}
	if err := db.Clear(); err != nil {
		log.Fatal("could not clear database:", err)
	}
}

func TestAddUser(t *testing.T) {
	testId, err := db.AddUser(testUserName, testUserEmail, testUserPassword)
	if err != nil {
		t.Fatalf("could not add user, %v", err)
	}
	user, err := db.GetUser(testUserId)
	if err != nil {
		t.Fatalf("could not retrieve user, %v", err)
	}
	assert.Equal(t, testId, user.Id)
	assert.Equal(t, testUserName, user.Name)
	assert.Equal(t, testUserEmail, user.Email)
	assert.Equal(t, testUserPassword, user.Password)
}

func TestGetUser(t *testing.T) {
	user, err := db.GetUser(testUserId)
	if err != nil {
		t.Fatalf("could not retrieve user, %v", err)
	}
	assert.Equal(t, testUserId, user.Id)
	assert.Equal(t, testUserName, user.Name)
	assert.Equal(t, testUserEmail, user.Email)
	assert.Equal(t, testUserPassword, user.Password)
}

func TestAddActivity(t *testing.T) {
	testId, err := db.AddActivity(testActivityName, testUserId)
	if err != nil {
		t.Fatalf("could not add activity, %v", err)
	}
	activity, err := db.GetActivity(testActivityId)
	if err != nil {
		t.Fatalf("could not retrieve activity, %v", err)
	}
	assert.Equal(t, testId, activity.Id)
	assert.Equal(t, testActivityName, activity.Name)
	assert.Equal(t, testUserId, activity.UserId)
}

func TestUpdateActivity(t *testing.T) {
	if err := db.UpdateActivity(testActivityId, testActivityNameUpdated); err != nil {
		t.Fatalf("could not update activity, %v", err)
	}
	activity, err := db.GetActivity(testActivityId)
	if err != nil {
		t.Fatalf("could not retrieve activity, %v", err)
	}
	assert.Equal(t, testActivityId, activity.Id)
	assert.Equal(t, testActivityNameUpdated, activity.Name)
	assert.Equal(t, testUserId, activity.UserId)
}

func TestAddBlock(t *testing.T) {
	testId, err := db.AddBlock(testBlockStartTime, testBlockEndTime, testActivityId)
	if err != nil {
		t.Fatalf("could not add block, %v", err)
	}
	block, err := db.GetBlock(testId)
	if err != nil {
		t.Fatalf("could not retrieve block, %v", err)
	}
	assert.Equal(t, testId, block.Id)
	assert.Equal(t, testBlockStartTime, block.StartTime)
	assert.Equal(t, testBlockEndTime, block.EndTime)
	assert.Equal(t, testActivityId, block.ActivityId)
}

func TestUpdateBlock(t *testing.T) {
	if err := db.UpdateBlock(testBlockId, testBlockStartTimeUpdated, testBlockEndTimeUpdated); err != nil {
		t.Fatalf("could not update block, %v", err)
	}
	block, err := db.GetBlock(testBlockId)
	if err != nil {
		t.Fatalf("could not retrieve block, %v", err)
	}
	assert.Equal(t, testBlockId, block.Id)
	assert.Equal(t, testBlockStartTimeUpdated, block.StartTime)
	assert.Equal(t, testBlockEndTimeUpdated, block.EndTime)
}

func TestAddPause(t *testing.T) {
	testId, err := db.AddPause(testPauseStartTime, testPauseEndTime, testBlockId)
	if err != nil {
		t.Fatalf("could not add pause, %v", err)
	}
	pause, err := db.GetPause(testPauseId)
	if err != nil {
		t.Fatalf("could not retrieve pause, %v", err)
	}
	assert.Equal(t, testId, pause.Id)
	assert.Equal(t, testPauseStartTime, pause.StartTime)
	assert.Equal(t, testPauseEndTime, pause.EndTime)
	assert.Equal(t, testBlockId, pause.BlockId)
}

func TestUpdatePause(t *testing.T) {
	if err := db.UpdatePause(testPauseId, testPauseStartTimeUpdated, testPauseEndTimeUpdated); err != nil {
		t.Fatalf("could not update pause, %v", err)
	}
	pause, err := db.GetPause(testPauseId)
	if err != nil {
		t.Fatalf("could not retrieve pause, %v", err)
	}
	assert.Equal(t, testPauseId, pause.Id)
	assert.Equal(t, testPauseStartTimeUpdated, pause.StartTime)
	assert.Equal(t, testPauseEndTimeUpdated, pause.EndTime)
}

func TestGetActivities(t *testing.T) {
	activities, err := db.GetActivities(testUserId)
	if err != nil {
		t.Fatalf("could not retrieve activities, %v", err)
	}
	assert.Equal(t, len(activities), 1)
	activity := activities[0]
	assert.Equal(t, testActivityId, activity.Id)
	assert.Equal(t, testActivityNameUpdated, activity.Name)
	assert.Equal(t, testUserId, activity.UserId)

	blocks := activity.Blocks
	assert.Equal(t, len(blocks), 1)
	block := blocks[0]
	assert.Equal(t, testBlockId, block.Id)
	assert.Equal(t, testBlockStartTimeUpdated, block.StartTime)
	assert.Equal(t, testBlockEndTimeUpdated, block.EndTime)

	pauses := block.Pauses
	assert.Equal(t, len(pauses), 1)
	pause := pauses[0]
	assert.Equal(t, testPauseId, pause.Id)
	assert.Equal(t, testPauseStartTimeUpdated, pause.StartTime)
	assert.Equal(t, testPauseEndTimeUpdated, pause.EndTime)
}

func TestGetActivity(t *testing.T) {
	activity, err := db.GetActivity(testActivityId)
	if err != nil {
		t.Fatalf("could not retrieve activity, %v", err)
	}
	assert.Equal(t, activity.Id, testActivityId)
	assert.Equal(t, activity.Name, testActivityNameUpdated)
	assert.Equal(t, activity.UserId, testUserId)

	blocks := activity.Blocks
	assert.Equal(t, len(blocks), 1)
	block := blocks[0]
	assert.Equal(t, testBlockId, block.Id)
	assert.Equal(t, testBlockStartTimeUpdated, block.StartTime)
	assert.Equal(t, testBlockEndTimeUpdated, block.EndTime)

	pauses := block.Pauses
	assert.Equal(t, len(pauses), 1)
	pause := pauses[0]
	assert.Equal(t, testPauseId, pause.Id)
	assert.Equal(t, testPauseStartTimeUpdated, pause.StartTime)
	assert.Equal(t, testPauseEndTimeUpdated, pause.EndTime)
}

func TestGetBlocks(t *testing.T) {
	blocks, err := db.GetBlocks(testActivityId)
	if err != nil {
		t.Fatalf("could not retrieve blocks, %v", err)
	}
	assert.Equal(t, len(blocks), 1)
	block := blocks[0]
	assert.Equal(t, testBlockId, block.Id)
	assert.Equal(t, testBlockStartTimeUpdated, block.StartTime)
	assert.Equal(t, testBlockEndTimeUpdated, block.EndTime)

	pauses := block.Pauses
	assert.Equal(t, len(pauses), 1)
	pause := pauses[0]
	assert.Equal(t, testPauseId, pause.Id)
	assert.Equal(t, testPauseStartTimeUpdated, pause.StartTime)
	assert.Equal(t, testPauseEndTimeUpdated, pause.EndTime)
}

func TestGetBlock(t *testing.T) {
	block, err := db.GetBlock(testBlockId)
	if err != nil {
		t.Fatalf("could not retrieve block, %v", err)
	}
	assert.Equal(t, testBlockId, block.Id)
	assert.Equal(t, testBlockStartTimeUpdated, block.StartTime)
	assert.Equal(t, testBlockEndTimeUpdated, block.EndTime)

	pauses := block.Pauses
	assert.Equal(t, len(pauses), 1)
	pause := pauses[0]
	assert.Equal(t, testPauseId, pause.Id)
	assert.Equal(t, testPauseStartTimeUpdated, pause.StartTime)
	assert.Equal(t, testPauseEndTimeUpdated, pause.EndTime)
}

func TestGetPauses(t *testing.T) {
	pauses, err := db.GetPauses(testBlockId)
	if err != nil {
		t.Fatalf("could not retrieve pauses, %v", err)
	}
	assert.Equal(t, len(pauses), 1)
	pause := pauses[0]
	assert.Equal(t, testPauseId, pause.Id)
	assert.Equal(t, testPauseStartTimeUpdated, pause.StartTime)
	assert.Equal(t, testPauseEndTimeUpdated, pause.EndTime)
}

func TestGetPause(t *testing.T) {
	pause, err := db.GetPause(testPauseId)
	if err != nil {
		t.Fatalf("could not retrieve pause, %v", err)
	}
	assert.Equal(t, testPauseId, pause.Id)
	assert.Equal(t, testPauseStartTimeUpdated, pause.StartTime)
	assert.Equal(t, testPauseEndTimeUpdated, pause.EndTime)
}

func TestGetCurrentBlock(t *testing.T) {
	id, err := db.AddBlock(testStartTimeCurrentBlock, "", testActivityId)
	if err != nil {
		t.Fatalf("could not add block, %v", err)
	}
	block, err := db.GetCurrentBlock()
	if err != nil {
		t.Fatalf("could not get get current block, %v", err)
	}
	assert.Equal(t, id, block.Id)
	assert.Equal(t, testStartTimeCurrentBlock, block.StartTime)
	assert.Equal(t, testEndTimeCurrentBlock, block.EndTime)
}

func TestDeleteByTableAndId(t *testing.T) {
	if err := db.DeleteByTableAndId("pauses", testPauseId); err != nil {
		t.Fatalf("could not delete pause, %v", err)
	}
	_, err := db.GetPause(testPauseId)
	assert.NotEqual(t, nil, err)
	if err := db.DeleteByTableAndId("blocks", testBlockId); err != nil {
		t.Fatalf("could not delete block, %v", err)
	}
	_, err = db.GetBlock(testBlockId)
	assert.NotEqual(t, nil, err)
	if err := db.DeleteByTableAndId("activities", testActivityId); err != nil {
		t.Fatalf("could not delete activity, %v", err)
	}
	_, err = db.GetActivity(testActivityId)
	assert.NotEqual(t, nil, err)
	if err := db.DeleteByTableAndId("users", testUserId); err != nil {
		t.Fatalf("could not delete user, %v", err)
	}
	_, err = db.GetUser(testUserId)
	assert.NotEqual(t, nil, err)
}

func TestClose(t *testing.T) {
	err := db.Close()
	if err != nil {
		t.Fatalf("could not close database, %v", err)
	}
}
