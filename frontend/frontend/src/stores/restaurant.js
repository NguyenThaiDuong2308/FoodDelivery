import { defineStore } from 'pinia'
import { ref } from 'vue'
import { restaurantService } from '../services/restaurant.service'

export const useRestaurantStore = defineStore('restaurant', () => {
    const restaurants = ref([])
    const selectedRestaurant = ref(null)
    const loading = ref(false)
    const error = ref(null)

    // ✅ GET /restaurant
    async function fetchRestaurants() {
        loading.value = true
        error.value = null
        try {
            const res = await restaurantService.getAll()
            // API trả về array trực tiếp, KHÔNG có wrapper .restaurants
            restaurants.value = Array.isArray(res) ? res : (res.restaurants || [])
            console.log('Fetched restaurants:', restaurants.value)
            return restaurants.value
        } catch (err) {
            error.value = err.message
            console.error('Error fetching restaurants:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // ✅ GET /restaurant/{restaurantId}
    async function fetchRestaurantById(id) {
        loading.value = true
        error.value = null
        try {
            const res = await restaurantService.getById(id)
            selectedRestaurant.value = res

            // Update in restaurants array if exists
            const index = restaurants.value.findIndex(r => r.id === id)
            if (index !== -1) {
                restaurants.value[index] = res
            }

            return res
        } catch (err) {
            error.value = err.message
            console.error('Error fetching restaurant:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // ✅ POST /restaurant
    async function createRestaurant(data) {
        loading.value = true
        error.value = null
        try {
            // API expects: { name, description, address, phone_number, email, status }
            const res = await restaurantService.create(data)

            // Add to local array
            restaurants.value.push(res)

            return res
        } catch (err) {
            error.value = err.message
            console.error('Error creating restaurant:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // ✅ PUT /restaurant/{restaurantId}
    async function updateRestaurant(id, data) {
        loading.value = true
        error.value = null
        try {
            // API expects: { name, description, address, phone_number, email, status }
            const res = await restaurantService.update(id, data)

            // Update in local array
            const index = restaurants.value.findIndex(r => r.id === id)
            if (index !== -1) {
                restaurants.value[index] = res
            }

            // Update selectedRestaurant if it's the same
            if (selectedRestaurant.value?.id === id) {
                selectedRestaurant.value = res
            }

            return res
        } catch (err) {
            error.value = err.message
            console.error('Error updating restaurant:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // ✅ GET /restaurant/{restaurantId}/menu
    async function fetchMenu(restaurantId) {
        loading.value = true
        error.value = null
        try {
            // API returns array of menu items directly
            const menuItems = await restaurantService.getMenu(restaurantId)

            // Update restaurant's menu_items in local state
            const restaurant = restaurants.value.find(r => r.id === restaurantId)
            if (restaurant) {
                restaurant.menu_items = menuItems
            }

            if (selectedRestaurant.value?.id === restaurantId) {
                selectedRestaurant.value.menu_items = menuItems
            }

            return menuItems
        } catch (err) {
            error.value = err.message
            console.error('Error fetching menu:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // ✅ POST /restaurant/{restaurantId}/menu
    async function createMenuItem(restaurantId, menuItemData) {
        loading.value = true
        error.value = null
        try {
            // API expects: { name, description?, price, available }
            const res = await restaurantService.createMenuItem(restaurantId, menuItemData)

            // Add to local state
            const restaurant = restaurants.value.find(r => r.id === restaurantId)
            if (restaurant) {
                if (!restaurant.menu_items) restaurant.menu_items = []
                restaurant.menu_items.push(res)
            }

            if (selectedRestaurant.value?.id === restaurantId) {
                if (!selectedRestaurant.value.menu_items) selectedRestaurant.value.menu_items = []
                selectedRestaurant.value.menu_items.push(res)
            }

            return res
        } catch (err) {
            error.value = err.message
            console.error('Error creating menu item:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // ✅ PUT /restaurant/{restaurantId}/menu/{menuItemId}
    async function updateMenuItem(restaurantId, menuItemId, menuItemData) {
        loading.value = true
        error.value = null
        try {
            // API expects: { name, description?, price, available }
            const res = await restaurantService.updateMenuItem(restaurantId, menuItemId, menuItemData)

            // Update in local state
            const restaurant = restaurants.value.find(r => r.id === restaurantId)
            if (restaurant && restaurant.menu_items) {
                const index = restaurant.menu_items.findIndex(i => i.id === menuItemId)
                if (index !== -1) {
                    restaurant.menu_items[index] = res
                }
            }

            if (selectedRestaurant.value?.id === restaurantId && selectedRestaurant.value.menu_items) {
                const index = selectedRestaurant.value.menu_items.findIndex(i => i.id === menuItemId)
                if (index !== -1) {
                    selectedRestaurant.value.menu_items[index] = res
                }
            }

            return res
        } catch (err) {
            error.value = err.message
            console.error('Error updating menu item:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // ✅ DELETE /restaurant/{restaurantId}/menu/{menuItemId}
    async function deleteMenuItem(restaurantId, menuItemId) {
        loading.value = true
        error.value = null
        try {
            await restaurantService.deleteMenuItem(restaurantId, menuItemId)

            // Remove from local state
            const restaurant = restaurants.value.find(r => r.id === restaurantId)
            if (restaurant && restaurant.menu_items) {
                restaurant.menu_items = restaurant.menu_items.filter(i => i.id !== menuItemId)
            }

            if (selectedRestaurant.value?.id === restaurantId && selectedRestaurant.value.menu_items) {
                selectedRestaurant.value.menu_items = selectedRestaurant.value.menu_items.filter(i => i.id !== menuItemId)
            }
        } catch (err) {
            error.value = err.message
            console.error('Error deleting menu item:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // Clear store
    function clearStore() {
        restaurants.value = []
        selectedRestaurant.value = null
        error.value = null
    }

    return {
        // State
        restaurants,
        selectedRestaurant,
        loading,
        error,

        // Actions
        fetchRestaurants,
        fetchRestaurantById,
        createRestaurant,
        updateRestaurant,
        fetchMenu,
        createMenuItem,
        updateMenuItem,
        deleteMenuItem,
        clearStore
    }
})