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
	Id int
	Name string
	UserId int
	Blocks []Block
}

type User struct {
	Id int
	Name string
	Email string
	Password string
}
