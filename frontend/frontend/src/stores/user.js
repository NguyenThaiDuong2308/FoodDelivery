import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userService } from '../services/user.service'
import { authService } from '../services/auth.service'

export const useUserStore = defineStore('user', () => {
    const users = ref([])
    const currentUser = ref(null)
    const loading = ref(false)
    const error = ref(null)

    const userById = computed(() => (id) => {
        return users.value.find(user => user.id === id)
    })

    async function fetchUsers() {
        loading.value = true
        error.value = null
        try {
            const result = await userService.getAll()
            users.value = result.users
            console.log(users)
        } catch (err) {
            error.value = err.message
            throw err
        } finally {
            console.log(users)
            loading.value = false
        }
    }

    async function fetchUserById(id) {
        loading.value = true
        error.value = null
        try {
            currentUser.value = await userService.getById(id)
            return currentUser.value
        } catch (err) {
            error.value = err.message
            throw err
        } finally {
            loading.value = false
        }
    }

    async function updateUser(id, userData) {
        loading.value = true
        error.value = null
        try {
            const updatedUser = await userService.update(id, userData)
            // Update in users array
            const index = users.value.findIndex(u => u.id === id)
            if (index !== -1) {
                users.value[index] = updatedUser
            }

            // Update currentUser if it's the same
            if (currentUser.value?.id === id) {
                currentUser.value = updatedUser
            }

            return updatedUser
        } catch (err) {
            error.value = err.message
            throw err
        } finally {
            loading.value = false
        }
    }

    async function deleteUser(id) {
        loading.value = true
        error.value = null
        try {
            await userService.delete(id)

            // Remove from users array
            users.value = users.value.filter(u => u.id !== id)

            // Clear currentUser if it's the deleted user
            if (currentUser.value?.id === id) {
                currentUser.value = null
            }
        } catch (err) {
            error.value = err.message
            throw err
        } finally {
            loading.value = false
        }
    }

    async function resetPassword(passwordData) {
        loading.value = true
        error.value = null
        try {
            await authService.resetPassword(passwordData)
        } catch (err) {
            error.value = err.message
            throw err
        } finally {
            loading.value = false
        }
    }

    async function fetchUserLocation(userId){
        loading.value = true
        error.value = null
        try {
            const locationData = await userService.getLocation(userId)
            // locationData = { id: number, address: string }
            console.log(`User ${userId} location:`, locationData)
            return locationData
        } catch (err) {
            error.value = err.message
            console.error(`Error fetching location for user ${userId}:`, err)
            throw err
        } finally {
            loading.value = false
        }
    }

    function setCurrentUser(user) {
        currentUser.value = user
    }

    function clearCurrentUser() {
        currentUser.value = null
    }

    function clearError() {
        error.value = null
    }

    return {
        users,
        currentUser,
        loading,
        error,
        userById,
        fetchUsers,
        fetchUserById,
        updateUser,
        deleteUser,
        resetPassword,
        fetchUserLocation,
        setCurrentUser,
        clearCurrentUser,
        clearError
    }
})