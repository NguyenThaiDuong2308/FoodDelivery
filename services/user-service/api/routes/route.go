package routes

import (
	"user-service/api/middlewares"
	"user-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) {
	authRoute := r.Group("/auth")
	authRoute.POST("/register", authHandler.Register)
	authRoute.POST("/login", authHandler.Login)
	authRoute.GET("/logout", middlewares.AuthMiddleware(), authHandler.Logout)
	authRoute.POST("/reset-password", middlewares.AuthMiddleware(), authHandler.ResetPassword)
	authRoute.POST("/refresh-token", authHandler.RefreshToken)
	authRoute.POST("/forgot-password", authHandler.ForgotPassword)
	authRoute.POST("/reset-forgot-password", authHandler.ResetForgotPassword)
	userRoute := r.Group("/user")
	userRoute.GET("/:id/get-location", userHandler.GetUserAddressByID)
	userRoute.Use(middlewares.AuthMiddleware())
	{
		userRoute.GET("/:id", userHandler.GetUserByID)
		userRoute.PUT("/:id", userHandler.UpdateUser)
		userRoute.DELETE("/:id", userHandler.DeleteUser)
		userRoute.GET("", userHandler.GetAllUserInfo)
	}
}
