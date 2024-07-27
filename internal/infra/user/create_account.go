package user_controller

import (
	"net/http"

	"github.com/MiguelCVA/mhook-backend/internal/database"
	"github.com/MiguelCVA/mhook-backend/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ResponseUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateUserController(ctx *gin.Context) {
	db := database.GetDatabase()

	var user models.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "Can't bind JSON: " + err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Can't hash password: " + err.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	err = db.Create(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Can't create user: " + err.Error(),
		})

		return
	}

	responseUser := ResponseUser{
		Name:  user.Name,
		Email: user.Email,
	}

	ctx.JSON(http.StatusCreated, responseUser)
}
