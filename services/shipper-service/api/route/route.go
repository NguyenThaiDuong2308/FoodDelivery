package route

import (
	"shipper-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	shipperHandler *handler.ShipperHandler,
) {
	shipperRoute := r.Group("/shipper")
	shipperRoute.GET("/:id", shipperHandler.GetShipperInfoByID)
	shipperRoute.PUT("/:id", shipperHandler.UpdateShipperStatus)
	shipperRoute.GET("", shipperHandler.GetShipperList)
	shipperRoute.PUT("/:id/location", shipperHandler.UpdateLocation)
	shipperRoute.GET("/:id/location", shipperHandler.GetShipperLocation)
}
