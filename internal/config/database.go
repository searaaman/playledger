package config 

import (
	"fmt"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/searaaman/playledger/internal/domain"
)

var DB *gorm.DB
func ConnectDatabase(){
	dsn:="host=localhost user=postgres password=12345678 dbname=playledger port =5432 sslmode=disable"
	db,err:=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil{
		log.Fatal("Failed to connect to db")
	}
	DB=db
	err=DB.AutoMigrate(
		&domain.Session{},
		/*&domain.Player{},
		&domain.TimeSlot{},*/
	)
	if err!=nil{
		log.Fatal("Automigrate was not succesfull",err)
	}
	fmt.Println("Succesfully connected to PostgreSQL")
}

