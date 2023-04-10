package main

import (
	"fmt"
	"log"

	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kilianmandscharo/activities/database"
	"github.com/kilianmandscharo/activities/schemas"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(("Could not load .env file"))
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PW"),
		os.Getenv("DB_NAME"))

	db, err := database.New(connStr)
	if err != nil {
		log.Fatal("could not open database")
	}

	defer db.Close()
	err = db.Init()
	if err != nil {
		log.Fatal("could not init database", err)
	}
	err = db.Clear()
	if err != nil {
		log.Fatal("could not clear database", err)
	}
	_, err = db.AddUser("Apollo", "test@gmail.com", "12345")
	if err != nil {
		log.Fatal("could not add user", err)
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/user", func(c *gin.Context) {
		var user schemas.UserCreate
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "could not read body"})
			return
		}
		if id, err := db.AddUser(user.Name, user.Email, user.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not add user"})
		} else {
			c.JSON(http.StatusOK, gin.H{"id": id})
		}
	})

	router.GET("/activities/:userId", func(c *gin.Context) {
		userId, _ := strconv.Atoi(c.Param("userId"))
		activities, err := db.GetActivities(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not get activities"})
		} else {
			c.JSON(http.StatusOK, activities)
		}
	})

	router.GET("/activity/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		activity, err := db.GetActivity(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not get activity"})
		} else {
			c.JSON(http.StatusOK, activity)
		}
	})

	router.POST("/activity", func(c *gin.Context) {
		var activity schemas.ActivityCreate
		if err := c.BindJSON(&activity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "could not read body"})
			return
		}
		if id, err := db.AddActivity(activity.Name, activity.UserId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not add activity"})
		} else {
			c.JSON(http.StatusOK, gin.H{"id": id})
		}
	})

	router.PUT("/activity", func(c *gin.Context) {
		var activity schemas.Activity
		if err := c.BindJSON(&activity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "could not read body"})
			return
		}
		err := db.UpdateActivity(activity.Id, activity.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not update activity"})
		} else {
			c.Status(http.StatusOK)
		}
	})

	router.DELETE("/activity/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := db.DeleteByTableAndId("activities", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not delete activity"})
		} else {
			c.Status(http.StatusOK)
		}
	})

	router.GET("/current", func(c *gin.Context) {
		block, err := db.GetCurrentBlock()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not get current block"})
		} else {
			c.JSON(http.StatusOK, block)
		}
	})

	router.GET("/blocks/:activityId", func(c *gin.Context) {
		activityId, _ := strconv.Atoi(c.Param("activityId"))
		blocks, err := db.GetBlocks(activityId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not get blocks"})
		} else {
			c.JSON(http.StatusOK, blocks)
		}
	})

	router.GET("/block/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		block, err := db.GetBlock(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "coult not get block"})
		} else {
			c.JSON(http.StatusOK, block)
		}
	})

	router.POST("/block", func(c *gin.Context) {
		var block schemas.BlockCreate
		if err := c.BindJSON(&block); err != nil {
      fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "could not read block"})
			return
		}
		id, err := db.AddBlock(block.StartTime, block.EndTime, block.ActivityId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not add block"})
			return
		}
		for _, pause := range block.Pauses {
			_, err := db.AddPause(pause.StartTime, pause.EndTime, id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "could not add pause"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	router.PUT("/block", func(c *gin.Context) {
		var block schemas.Block
		if err := c.BindJSON(&block); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "could not read block"})
			return
		}
		if err := db.UpdateBlock(block.Id, block.StartTime, block.EndTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not update block"})
			return
		}
		if err := db.DeletePauses(block.Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not update pauses"})
			return
		}
		for _, pause := range block.Pauses {
			_, err := db.AddPause(pause.StartTime, pause.EndTime, block.Id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "could not update pause"})
				return
			}
		}
		c.Status(http.StatusOK)
	})

	router.DELETE("/block/:id", func(c *gin.Context) {
		blockId, _ := strconv.Atoi(c.Param("id"))
		err := db.DeleteByTableAndId("blocks", blockId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not delete block"})
		} else {
			c.Status(http.StatusOK)
		}
	})

	router.GET("/pause/:blockId", func(c *gin.Context) {
		blockId, _ := strconv.Atoi(c.Param("blockId"))
		pauses, err := db.GetBlocks(blockId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not get pauses"})
		} else {
			c.JSON(http.StatusOK, pauses)
		}
	})

	router.POST("/pause", func(c *gin.Context) {
		var pause schemas.PauseCreate
		if err := c.BindJSON(&pause); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "could not read pause"})
			return
		}
		if id, err := db.AddPause(pause.StartTime, pause.EndTime, pause.BlockId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not add pause"})
		} else {
			c.JSON(http.StatusOK, gin.H{"id": id})
		}
	})

	router.PUT("/pause", func(c *gin.Context) {
		var pause schemas.Pause
		if err := c.BindJSON(&pause); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "could not read pause"})
			return
		}
		if err := db.UpdatePause(pause.Id, pause.StartTime, pause.EndTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not update pause"})
		} else {
			c.Status(http.StatusOK)
		}
	})

	router.DELETE("/pause/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := db.DeleteByTableAndId("blocks", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not delete pause"})
		} else {
			c.Status(http.StatusOK)
		}
	})

	router.Run(":8080")
}
