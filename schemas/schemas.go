package schemas

import "database/sql"

type TableSchema struct {
	Name    string
	Columns string
}

type Pause struct {
	Id        int    `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	BlockId   int    `json:"blockId"`
}

type Block struct {
	Id         int     `json:"id"`
	StartTime  string  `json:"startTime"`
	EndTime    string  `json:"endTime"`
	ActivityId int     `json:"activityId"`
	Pauses     []Pause `json:"pauses"`
}

type Activity struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	UserId int     `json:"userId"`
	Blocks []Block `json:"blocks"`
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type UserCreate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ActivityCreate struct {
	Name   string `json:"name" binding:"required"`
	UserId int    `json:"userId" binding:"required"`
}

type BlockCreate struct {
	StartTime  string        `json:"startTime" binding:"required"`
	EndTime    string        `json:"endTime" binding:"required"`
	ActivityId int           `json:"activityId" binding:"required"`
	Pauses     []PauseCreate `json:"pauses"`
}

type PauseCreate struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
	BlockId   int    `json:"blockId"`
}

type CurrentBlock struct {
	Id         int            `json:"id"`
	StartTime  string         `json:"startTime"`
	EndTime    sql.NullString `json:"endTime"`
	ActivityId int            `json:"activityId"`
	Pauses     []Pause        `json:"pauses"`
}
