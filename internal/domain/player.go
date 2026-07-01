package domain

type Player struct{
	ID uint `gorm:"primaryKey"`
	Name string
	PhoneNumber string 
}