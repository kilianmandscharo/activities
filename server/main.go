package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	router.POST("/user", func(c *gin.Context) {
		var user schemas.UserCreate

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "unauthorized"})
		}

		if err := database.AddUser(db, user.Name, user.Email, user.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "user created"})
		}

	})

	router.Run(":8080")
}

