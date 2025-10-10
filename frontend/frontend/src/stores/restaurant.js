import { defineStore } from 'pinia'
import { ref } from 'vue'
import { restaurantService } from '../services/restaurant.service'

export const useRestaurantStore = defineStore('restaurant', () => {
    const restaurants = ref([])
    const selectedRestaurant = ref(null)
    const loading = ref(false)

    async function fetchRestaurants() {
        loading.value = true
        try {
            const res = await restaurantService.getAll()
            restaurants.value = res.restaurants
            console.log(JSON.stringify(restaurants.value))
        } finally {
            loading.value = false
        }
    }

    async function fetchRestaurantById(id) {
        loading.value = true
        try {
            selectedRestaurant.value = await restaurantService.getById(id)
        } finally {
            loading.value = false
        }
    }

    async function createRestaurant(data) {
        loading.value = true
        try {
            const res = await restaurantService.create(data)
            restaurants.value.push(res)
            return res
        } finally {
            loading.value = false
        }
    }

    async function updateRestaurant(id, data) {
        loading.value = true
        try {
            const res = await restaurantService.update(id, data)
            // Cập nhật trong mảng local
            const index = restaurants.value.findIndex(r => r.id === id)
            if (index !== -1) restaurants.value[index] = res
            return res
        } finally {
            loading.value = false
        }
    }

    async function fetchMenu(restaurantId) {
        loading.value = true
        try {
            const res = await restaurantService.getMenu(restaurantId)
            return res
        } finally {
            loading.value = false
        }
    }

    async function createMenuItem(restaurantId, menuItemData) {
        loading.value = true
        try {
            const res = await restaurantService.createMenuItem(restaurantId, menuItemData)
            console.log(res)
            const restaurant = restaurants.value.find(r => r.id === restaurantId)
            if (restaurant) restaurant.menu_items.push(res)
            if (selectedRestaurant.value?.id === restaurantId) selectedRestaurant.value.menu_items.push(res)
            return res
        } finally {
            loading.value = false
        }
    }

    async function updateMenuItem(restaurantId, menuItemId, menuItemData) {
        loading.value = true
        try {
            const res = await restaurantService.updateMenuItem(restaurantId, menuItemId, menuItemData)
            const restaurant = restaurants.value.find(r => r.id === restaurantId)
            if (restaurant) {
                const index = restaurant.menu_items.findIndex(i => i.id === menuItemId)
                if (index !== -1) restaurant.menu_items[index] = res
            }
            if (selectedRestaurant.value?.id === restaurantId) {
                const index = selectedRestaurant.value.menu_items.findIndex(i => i.id === menuItemId)
                if (index !== -1) selectedRestaurant.value.menu_items[index] = res
            }
            return res
        } finally {
            loading.value = false
        }
    }

    async function deleteMenuItem(restaurantId, menuItemId) {
        loading.value = true
        try {
            await restaurantService.deleteMenuItem(restaurantId, menuItemId)
            const restaurant = restaurants.value.find(r => r.id === restaurantId)
            if (restaurant) restaurant.menu_items = restaurant.menu_items.filter(i => i.id !== menuItemId)
            if (selectedRestaurant.value?.id === restaurantId) {
                selectedRestaurant.value.menu_items = selectedRestaurant.value.menu_items.filter(i => i.id !== menuItemId)
            }
        } finally {
            loading.value = false
        }
    }

    return {
        restaurants,
        selectedRestaurant,
        loading,
        fetchRestaurants,
        fetchRestaurantById,
        createRestaurant,
        updateRestaurant,
        fetchMenu,
        createMenuItem,
        updateMenuItem,
        deleteMenuItem
    }
})