import axios from 'axios'
import { API_BASE_URL } from '../config/api'

const apiClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json'
    }
})

apiClient.interceptors.request.use(
    (config) => {
        const token = sessionStorage.getItem('token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => Promise.reject(error)
)

apiClient.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true

            try {
                const refreshToken = sessionStorage.getItem('refreshToken')
                if (refreshToken) {
                    const { data } = await axios.post(`${API_BASE_URL}/auth/refresh-token`, {
                        refresh_token: refreshToken
                    })

                    sessionStorage.setItem('token', data.access_token)

                    originalRequest.headers.Authorization = `Bearer ${data.accessToken}`
                    return apiClient(originalRequest)
                }
            } catch (refreshError) {
                sessionStorage.removeItem('token')
                sessionStorage.removeItem('user')
                sessionStorage.removeItem('refreshToken')
                window.location.href = '/login'
                return Promise.reject(refreshError)
            }
        }

        if (error.response?.status === 401) {
            sessionStorage.removeItem('token')
            sessionStorage.removeItem('user')
            sessionStorage.removeItem('refreshToken')
            window.location.href = '/login'
        }

        return Promise.reject(error)
    }
)

export default apiClient