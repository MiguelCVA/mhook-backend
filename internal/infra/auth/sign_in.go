package auth_controller

import (
	"net/http"
	"time"

	"github.com/MiguelCVA/mhook-backend/internal/database"
	"github.com/MiguelCVA/mhook-backend/internal/models"
	"github.com/MiguelCVA/mhook-backend/internal/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Body struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func SignIn(ctx *gin.Context) {
	var requestBody Body
	db := database.GetDatabase()

	err := ctx.ShouldBindJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Can't bind JSON: " + err.Error(),
		})
		return
	}

	var user models.User
	// if err := db.Where("email = ? AND password = ?", requestBody.Email, requestBody.Password).First(&user).Error; err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		ctx.JSON(http.StatusUnauthorized, gin.H{
	// 			"error": "Incorrect email or password",
	// 		})
	// 	} else {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{
	// 			"error": "Error finding user: " + err.Error(),
	// 		})
	// 	}
	// 	return
	// }

	if err := db.Where("email = ?", requestBody.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Incorrect email or password",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error finding user: " + err.Error(),
			})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Incorrect email or password",
		})
		return
	}

	session := models.Session{
		UserID: user.ID,
	}
	err = db.Create(&session).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Can't create session: " + err.Error(),
		})
		return
	}

	jwt, err := services.GenerateJWT(map[string]interface{}{
		"email":   requestBody.Email,
		"session": session.ID,
	}, 168*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   jwt,
	})
}
