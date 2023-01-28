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
	db.Init()
	db.Clear()

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/user", func(c *gin.Context) {
		var user schemas.UserCreate

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "unauthorized"})
			return
		}

		if _, err := db.AddUser(user.Name, user.Email, user.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "user created"})
		}
	})

	router.POST("/activity", func(c *gin.Context) {
		var activity schemas.ActivityCreate

		if err := c.BindJSON(&activity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "unauthorized"})
			return
		}

		if _, err := db.AddActivity(activity.Name, activity.UserId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "activity created"})
		}
	})

	// router.DELETE("/activity/:activityId", func(c *gin.Context) {
	// 	activityId, _ := strconv.Atoi(c.Param("activityId"))
	// 	err := db.DeleteByTableAndId("activities", activityId)
	// 	if err != nil {
	// 		c.JSON(http.StatusNotFound, gin.H{"status": "activity not found"})
	// 	} else {
	// 		c.JSON(http.StatusOK, gin.H{"status": "activity deleted"})
	// 	}
	// })

	router.GET("/activities/:userId", func(c *gin.Context) {
		userId, _ := strconv.Atoi(c.Param("userId"))
		activities, _ := db.GetActivities(userId)
		c.JSON(http.StatusOK, activities)
	})

	router.POST("/block", func(c *gin.Context) {
		var block schemas.BlockCreate

		if err := c.BindJSON(&block); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "unauthorized"})
			return
		}

		if _, err := db.AddBlock(block.StartTime, block.EndTime, block.ActivityId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "block created"})
		}
	})

	router.Run(":8080")
}
