import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '../services/auth.service'

export const useAuthStore = defineStore('auth', () => {
    const user = ref(null)
    const token = ref(null)

    const isAuthenticated = computed(() => !!token.value)

    function setAuth(userData, authToken) {
        user.value = userData
        token.value = authToken
        sessionStorage.setItem('token', authToken)
        sessionStorage.setItem('user', JSON.stringify(userData))
    }

    function clearAuth() {
        user.value = null
        token.value = null
        sessionStorage.removeItem('token')
        sessionStorage.removeItem('user')
    }

    async function login(credentials) {
        const data = await authService.login(credentials)
      //  const users = await authService.getCurrentUser()
        if (users.length > 0) {
            setAuth(users[0], data.accessToken)
        }
        console.log("login successful:", data)
        return data
    }

    async function register(userData) {
        return await authService.register(userData)
    }

    async function logout() {
        try {
            await authService.logout()
        } finally {
            clearAuth()
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
        checkAuth
    }
})