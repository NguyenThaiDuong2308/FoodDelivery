package handler

import (
	"net/http"
	"shipper-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShipperHandler struct {
	shipperService service.ShipperService
}

func NewShipperHandler(shipperService service.ShipperService) *ShipperHandler {
	return &ShipperHandler{shipperService: shipperService}
}

func (h *ShipperHandler) GetShipperInfoByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	shipper, err := h.shipperService.GetShipperInfoByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"shipper": shipper})
}

type UpdateShipperStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *ShipperHandler) UpdateShipperStatus(c *gin.Context) {
	var req UpdateShipperStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	existingShipper, err := h.shipperService.GetShipperInfoByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	existingShipper.Status = req.Status
	err = h.shipperService.UpdateShipperStatus(c.Request.Context(), existingShipper)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"shipper": existingShipper})
}

func (h *ShipperHandler) GetShipperList(c *gin.Context) {
	shippers, err := h.shipperService.GetShipperList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"shippers": shippers})
}

type UpdateLocationRequest struct {
	Longitude float64 `json:"longitude" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
}

func (h *ShipperHandler) UpdateLocation(c *gin.Context) {
	shipperID, _ := strconv.Atoi(c.Param("id"))
	var req UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := h.shipperService.UpdateShipperLocation(c.Request.Context(), uint(shipperID), req.Longitude, req.Latitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"shipperID": shipperID,
		"longitude": req.Longitude,
		"latitude":  req.Latitude,
	})

}

func (h *ShipperHandler) GetShipperLocation(c *gin.Context) {
	shipperID, _ := strconv.Atoi(c.Param("id"))
	location, err := h.shipperService.GetShipperLocation(c.Request.Context(), uint(shipperID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"shipper":  shipperID,
		"location": location,
	})
}
