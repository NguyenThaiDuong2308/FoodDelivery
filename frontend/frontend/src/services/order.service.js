import apiClient from './api.service'
import { API_ENDPOINTS } from '../config/api'

export const orderService = {
    async create(orderData) {
        const response = await apiClient.post(API_ENDPOINTS.ORDERS, orderData)
        return response.data
    },

    async getCustomerOrders(customerId) {
        const response = await apiClient.get(API_ENDPOINTS.CUSTOMER_ORDERS(customerId))
        return response.data
    },

    async getOrderDetails(orderID){
        const response = await apiClient.get(API_ENDPOINTS.ORDER_BY_ID(orderID))
        return response.data
    },

    async updateStatus(orderId, status) {
        const response = await apiClient.post(API_ENDPOINTS.ORDER_BY_ID(orderId), {order_id: orderId, status })
        return response.data
    },
    async getRestaurantOrders(restaurantId) {
        const response = await apiClient.get(API_ENDPOINTS.RESTAURANT_ORDERS(restaurantId))
        return response.data
    },

    async getShipperOrders(shipperId) {
        const response = await apiClient.get(API_ENDPOINTS.SHIPPER_ORDERS(shipperId))
        return response.data
    }
}