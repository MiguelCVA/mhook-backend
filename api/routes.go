package api

import (
	"log"
	"net/http"
	"time"

	auth_controller "github.com/MiguelCVA/mhook-backend/internal/infra/auth"
	user_controller "github.com/MiguelCVA/mhook-backend/internal/infra/user"
	"github.com/MiguelCVA/mhook-backend/pkg"
	"github.com/gin-gonic/gin"
)

func ConfigRoutes(router *gin.Engine) *gin.Engine {
	api := router.Group("api")
	{
		v1 := api.Group("/v1")
		{
			user := v1.Group("user")
			{
				user.GET("/", func(ctx *gin.Context) {
					ctx.JSON(200, gin.H{
						"message": "Hello World",
					})
				})
				user.POST("/", user_controller.CreateUserController)
			}
			auth := v1.Group("auth")
			{
				auth.POST("/sign-in", auth_controller.SignIn)
			}
		}
		api.GET("/token", func(ctx *gin.Context) {
			genJWT, err := pkg.GenerateJWT(map[string]interface{}{
				"a": "b",
			}, 1*time.Hour)

			if err != nil {
				log.Fatalf("Error generating JWT: %v", err)
			}

			ctx.JSON(http.StatusCreated, gin.H{
				"token": genJWT,
			})
		})
	}

	return router
}
