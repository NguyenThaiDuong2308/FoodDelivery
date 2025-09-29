package handlers

//
//import (
//	"net/http"
//	"shipper-assignment-service/internal/models"
//	"shipper-assignment-service/internal/service"
//
//	"github.com/gin-gonic/gin"
//)
//
//type assignRequest struct {
//	OrderEvent models.OrderEvent `json:"order_event" binding:"required"`
//}
//
//type assignResponse struct {
//	Assignment *models.Assignment `json:"assignment,omitempty"`
//}
//
//type ShipperHandler struct {
//	shipperAssignmentService service.ShipperAssignmentService
//}
//
//func NewShipperAssignmentHandler(s service.ShipperAssignmentService) *ShipperHandler {
//	return &ShipperHandler{
//		shipperAssignmentService: s,
//	}
//}
//
//func (h *ShipperHandler) AssignNearestShipper(c *gin.Context) {
//	var req assignRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	assignment, err := h.shipperAssignmentService.AssignNearestShipper(c.Request.Context(), &req.OrderEvent)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, assignResponse{
//		Assignment: assignment,
//	})
//}
