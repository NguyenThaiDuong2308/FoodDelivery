package route

import (
	"shipper-assignment-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	shipperAssignmentHanlder *handlers.ShipperHandler,
) {
	assignRoute := r.Group("/assignment")
	assignRoute.POST("", shipperAssignmentHanlder.AssignNearestShipper)
}
