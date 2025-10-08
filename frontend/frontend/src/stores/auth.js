import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '../services/auth.service'
import { jwtDecode }from 'jwt-decode';
import apiClient from "@/services/api.service.js";
export const useAuthStore = defineStore('auth', () => {
    let user = ref(null)
    const token = ref(null)
    const refreshToken = ref(null)

    const isAuthenticated = computed(() => !!token.value)

    function setAuth(userData, authToken, refToken = null) {
        user.value = userData
        token.value = authToken
        refreshToken.value = refToken
        sessionStorage.setItem('token', authToken)
        sessionStorage.setItem('user', JSON.stringify(userData))
        if(refToken){
            sessionStorage.setItem('refreshToken', refToken)
        }
    }

    function clearAuth() {
        user.value = null
        token.value = null
        refreshToken.value = null
        sessionStorage.removeItem('token')
        sessionStorage.removeItem('user')
        sessionStorage.removeItem('refreshToken')
    }

    async function login(credentials) {
        try {
            const data = await authService.login(credentials)
            const accessToken = data.accessToken
            const refreshToken = data.refreshToken
            console.log(accessToken)
            console.log(refreshToken)
            const decoded = jwtDecode(accessToken)
            console.log(decoded)
            const userId = decoded.user_id
            console.log(userId)

            apiClient.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`
            const currentUser = await authService.getCurrentUser(userId)
            if (currentUser) {
                setAuth(currentUser, accessToken, refreshToken)
            }
            console.log("login successful:", data)
            return data
        } catch (error){
            console.error("Login failed", error)
            throw error
        }
    }

    async function register(userData) {
        try{
           return await authService.register(userData)
        } catch (error){
            console.error("Registration failed:", error)
            throw error;
        }
    }

    async function logout() {
        try {
            await authService.logout()
        } finally {
            clearAuth()
        }
    }

    async function resetPassword(passwordData) {
        return await authService.resetPassword(passwordData)
    }

    async function forgotPassword(email) {
        return await authService.forgotPassword(email)
    }

    async function resetForgotPassword(resetData) {
        return await authService.resetForgotPassword(resetData)
    }

    async function refreshAccessToken() {
        if (!refreshToken.value) {
            throw new Error('No refresh token available')
        }
        try {
            const data = await authService.refreshToken(refreshToken.value)
            token.value = data.accessToken
            sessionStorage.setItem('token', data.accessToken)
            if (data.refreshToken) {
                refreshToken.value = data.refreshToken
                sessionStorage.setItem('refreshToken', data.refreshToken)
            }
            return data
        } catch (error) {
            clearAuth()
            throw error
        }
    }

    function checkAuth() {
        const savedToken = sessionStorage.getItem('token')
        const savedUser = sessionStorage.getItem('user')
        if (savedToken && savedUser) {
            token.value = savedToken
            user.value = JSON.parse(savedUser)
        }
    }

    return {
        user,
        token,
        isAuthenticated,
        login,
        register,
        logout,
        resetPassword,
        forgotPassword,
        resetForgotPassword,
        refreshAccessToken,
        checkAuth,
    }
})