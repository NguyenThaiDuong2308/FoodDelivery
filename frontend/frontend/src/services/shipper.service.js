import apiClient from './api.service'
import { API_ENDPOINTS } from '../config/api'

export const shipperService = {
    async getAll() {
        const response = await apiClient.get(API_ENDPOINTS.SHIPPERS)
        return response.data
    },

    async getById(id) {
        const response = await apiClient.get(API_ENDPOINTS.SHIPPER_BY_ID(id))
        return response.data
    },

    async updateLocation(id, location) {
        const response = await apiClient.put(API_ENDPOINTS.SHIPPER_LOCATION(id), location)
        return response.data
    },

    async updateStatus(id, status) {
        const response = await apiClient.put(API_ENDPOINTS.SHIPPER_BY_ID(id), { status })
        return response.data
    }
}