package schemas

type TableSchema struct {
	Name string
	Columns string
} 

type Pause struct {
	Id int
	StartTime string	
	EndTime string
	BlockId int
}

type Block struct {
	Id int
	StartTime string
	EndTime string
	ActivityId int
	Pauses []Pause
}

type Activity struct {
  Id int `json:"id"`
	Name string `json:"name"`
  UserId int `json:"userId"`
  Blocks []Block `json:"blocks"`
}

type User struct {
	Id int
	Name string
	Email string
	Password string
}

type UserCreate struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ActivityCreate struct {
	Name string `json:"name" binding:"required"`
	UserId int `json:"userId" binding:"required"`
}

type BlockCreate struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime string `json:"endTime" binding:"required"`
	ActivityId int `json:"activityId" binding:"required"`
}

type PauseCreate struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime string `json:"endTime" binding:"required"`
	BlockId int `json:"blockId" binding:"required"`
}
