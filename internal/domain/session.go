package domain
import "time"

type Session struct{
	ID uint `gorm:"primaryKey"`
	StartTime time.Time
	EndTime time.Time
	CourtPrice float64
	/*TimeSlots []*TimeSlot*/

}

