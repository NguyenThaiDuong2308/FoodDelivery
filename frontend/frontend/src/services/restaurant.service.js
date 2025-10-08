import apiClient from './api.service'
import { API_ENDPOINTS } from '../config/api'
export const restaurantService = {
    async getAll() {
        const response = await apiClient.get(API_ENDPOINTS.RESTAURANTS)
        return response.data
    },

    async getById(id) {
        const response = await apiClient.get(API_ENDPOINTS.RESTAURANT_BY_ID(id))
        return response.data
    },

    async create(restaurantData) {
        const response = await apiClient.post(API_ENDPOINTS.RESTAURANTS, restaurantData)
        return response.data
    },

    async update(id, restaurantData) {
        const response = await apiClient.put(API_ENDPOINTS.RESTAURANT_BY_ID(id), restaurantData)
        return response.data
    },

    async getMenu(id) {
        const response = await apiClient.get(API_ENDPOINTS.RESTAURANT_MENU(id))
        return response.data
    },

    async createMenuItem(restaurantId, menuItemData) {
        const response = await apiClient.post(
            API_ENDPOINTS.RESTAURANT_MENU(restaurantId),
            menuItemData
        )
        return response.data
    },

    async getMenuItem(restaurantId, menuItemId) {
        const response = await apiClient.get(
            API_ENDPOINTS.MENU_ITEM(restaurantId, menuItemId)
        )
        return response.data
    },

    async updateMenuItem(restaurantId, menuItemId, menuItemData) {
        const response = await apiClient.put(
            API_ENDPOINTS.MENU_ITEM(restaurantId, menuItemId),
            menuItemData
        )
        return response.data
    },

    async deleteMenuItem(restaurantId, menuItemId) {
        const response = await apiClient.delete(
            API_ENDPOINTS.MENU_ITEM(restaurantId, menuItemId)
        )
        return response.data
    }
}