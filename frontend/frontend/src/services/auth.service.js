import apiClient from './api.service'
import { API_ENDPOINTS } from '../config/api'

export const authService = {
    async login(credentials) {
        const response = await apiClient.post(API_ENDPOINTS.LOGIN, credentials)
        console.log("data after login:", response.data)
        return response.data
    },

    async register(userData) {
        const response = await apiClient.post(API_ENDPOINTS.REGISTER, userData)
        return response.data
    },

    async logout() {
        const response = await apiClient.get(API_ENDPOINTS.LOGOUT)
        return response.data
    },

    async getCurrentUser(id) {
        const response = await apiClient.get(API_ENDPOINTS.USER_BY_ID(id))
        return response.data
    },
    async resetPassword(passwordData) {
        const response = await apiClient.post(API_ENDPOINTS.RESET_PASSWORD, passwordData)
        return response.data
    },

    async forgotPassword(email) {
        const response = await apiClient.post(API_ENDPOINTS.FORGOT_PASSWORD, { email })
        return response.data
    },

    async resetForgotPassword(resetData) {
        const response = await apiClient.post(API_ENDPOINTS.RESET_FORGOT_PASSWORD, resetData)
        return response.data
    },

    async refreshAccessToken(refreshToken) {
        const response = await apiClient.post(API_ENDPOINTS.REFRESH_TOKEN, {
            refresh_token: refreshToken
        })
        return response.data
    }
}
