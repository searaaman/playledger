package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/searaaman/playledger/internal/config"
	"github.com/searaaman/playledger/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`}

func Register(ctx *gin.Context) {

	var request RegisterRequest

	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := domain.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
	}

	err = config.DB.Create(&user).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    gin.H{
			"id":user.ID,
			"name":user.Name,
			"email":user.Email,
		},
	})
}

func Login(ctx *gin.Context){
	var request LoginRequest

	err :=ctx.ShouldBindJSON(&request)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"Error":err.Error(),
		})
		return
	}

	var user domain.User

	err=config.DB.Where("email =?",request.Email).First(&user).Error

	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"Error":err.Error(),
		})
		return
	}

	err=bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(request.Password),

	)

	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"Error":"Invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"user_id": user.ID,
	"email":   user.Email,
	})

	tokenString, err := token.SignedString([]byte(config.JWTSecret))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}


	ctx.JSON(http.StatusOK,gin.H{
		"Message":"Login Succesful",
		"token":   tokenString,
		"user":gin.H{
			"id":user.ID,
			"name":user.Name,
			"email":user.Email,
		},
	})
}