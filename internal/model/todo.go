package model

type Todo struct {
	ID        uint `gorm:"primarykey"`
	Title     string
	Completed bool
}
