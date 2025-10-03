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

    async getMenu(id) {
        const response = await apiClient.get(API_ENDPOINTS.RESTAURANT_MENU(id))
        return response.data
    }
}