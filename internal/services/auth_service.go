package services

import(
	"errors"
	"github.com/searaaman/playledger/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB,name string,email string,password string)(*domain.User,error){
	var existingUser domain.User

	result :=db.Where("email=?",email).First(&existingUser)

	if result.Error==nil{
		return nil,errors.New("email already exists")
	}

	passwordHash,err:=bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err!=nil{
		return nil,err
	}

	user:=domain.User{
		Name:  name,
		Email: email,
		PasswordHash :string(passwordHash),


	}
	result =db.Create(&user)
	if result.Error !=nil{
		return nil,result.Error
	}
	return  &user,nil

}