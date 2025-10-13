import apiClient from './api.service'
import { API_ENDPOINTS } from '../config/api'

export const userService = {
    async getAll() {
        const response = await apiClient.get(API_ENDPOINTS.USERS)
        return response.data
    },

    async getById(id) {
        const response = await apiClient.get(API_ENDPOINTS.USER_BY_ID(id))
        return response.data
    },

    async update(id, userData) {
        const response = await apiClient.put(API_ENDPOINTS.USER_BY_ID(id), userData)
        return response.data
    },

    async delete(id) {
        const response = await apiClient.delete(API_ENDPOINTS.USER_BY_ID(id))
        return response.data
    },
    async getLocation(id){
        const response = await apiClient.get(API_ENDPOINTS.USER_LOCATION(id))
        return response.data
    }
}