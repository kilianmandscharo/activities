package main

import (
	"database/sql"
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

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PW"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.ClearTables(db)
	database.InitTables(db)

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/user", func(c *gin.Context) {
		var user schemas.UserCreate

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "unauthorized"})
			return
		}

		if err := database.AddUser(db, user.Name, user.Email, user.Password); err != nil {
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

		if err := database.AddActivity(db, activity.Name, activity.UserId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "activity created"})
		}
	})

	router.DELETE("/activity/:activityId", func(c *gin.Context) {
		activityId, _ := strconv.Atoi(c.Param("activityId"))
		err := database.DeleteByTableAndId(db, "activities", activityId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "activity not found"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "activity deleted"})
		}
	})

	router.GET("/activities/:userId", func(c *gin.Context) {
		userId, _ := strconv.Atoi(c.Param("userId"))
		activities := database.GetAllActivities(db, userId)
		c.JSON(http.StatusOK, activities)
	})

	router.POST("/block", func(c *gin.Context) {
		var block schemas.BlockCreate

		if err := c.BindJSON(&block); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "unauthorized"})
			return
		}

		if err := database.AddBlock(db, block.StartTime, block.EndTime, block.ActivityId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "block created"})
		}
	})

	router.Run(":8080")
}
