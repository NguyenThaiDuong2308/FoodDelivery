package route

import (
	"restaurant-service/api/middlewares"
	"restaurant-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	restaurantHandler *handler.RestaurantHandler,
	menuItemHandler *handler.MenuItemHandler,
) {
	restaurantRoute := r.Group("/restaurant")
	restaurantRoute.POST("", middlewares.AuthMiddleware(), restaurantHandler.CreateRestaurant)
	restaurantRoute.GET("/:restaurant_id", restaurantHandler.GetRestaurantInfoByID)
	restaurantRoute.PUT("/:restaurant_id", middlewares.AuthMiddleware(), restaurantHandler.UpdateRestaurant)
	restaurantRoute.GET("", restaurantHandler.GetRestaurantList)
	restaurantRoute.GET("/:restaurant_id/menu", menuItemHandler.GetMenuListByRestaurantInfo)
	restaurantRoute.GET("/:restaurant_id/menu/:id", menuItemHandler.GetMenuItemByItemID)
	restaurantRoute.POST("/:restaurant_id/menu", middlewares.AuthMiddleware(), menuItemHandler.CreateMenuItem)
	restaurantRoute.PUT("/:restaurant_id/menu/:id", middlewares.AuthMiddleware(), menuItemHandler.UpdateMenuItem)
	restaurantRoute.DELETE("/:restaurant_id/menu/:id", middlewares.AuthMiddleware(), menuItemHandler.DeleteMenuItem)
}
