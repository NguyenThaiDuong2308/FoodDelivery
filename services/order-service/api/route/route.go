package route

import (
	"order-service/api/middlewares"
	"order-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	orderHandler *handler.OrderHandler,
) {
	orderRoute := r.Group("/order")
	orderRoute.POST("", middlewares.AuthMiddleware(), orderHandler.CreateOrder)
	orderRoute.GET("/:id", middlewares.AuthMiddleware(), orderHandler.GetOrderByID)
	orderRoute.POST("/:id", middlewares.AuthMiddleware(), orderHandler.UpdateOrderStatus)
	orderRoute.GET("/customer/:customer_id", middlewares.AuthMiddleware(), orderHandler.GetOrderByCustomerID)
}
