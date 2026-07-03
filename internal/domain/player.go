package domain

type Player struct{
	ID uint `gorm:"primaryKey"`

	Name string
	Phone string 

	TimeSlots []TimeSlot `gorm:"many2many:player_time_slots:"`
}

type AssignPlayerRequest struct{
	PlayerID uint `json:"player_id"`
}