package models

type OrderEvent struct {
	EventName    string `json:"event_name"`
	OrderID      uint   `json:"order_id"`
	RestaurantID int    `json:"restaurant_id"`
}

type Assignment struct {
	OrderID   uint    `json:"order_id"`
	ShipperID int     `json:"shipper_id"`
	Distance  float64 `json:"distance"`
}

type Coordinate struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
