package handlers
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/searaaman/playledger/internal/domain"
	"github.com/searaaman/playledger/internal/config"
)

func CreateSession(ctx *gin.Context){
	var session domain.Session
	err:=ctx.ShouldBindBodyWithJSON(&session)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":"err.Error()",
		})
		return
	}
	err=config.DB.Create(&session).Error
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"error":err.Error(),
		})
		return


	}
	ctx.JSON(http.StatusCreated,session)
	
}
