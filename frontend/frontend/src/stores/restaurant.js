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
            restaurants.value = await restaurantService.getAll()
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

    return {
        restaurants,
        selectedRestaurant,
        loading,
        fetchRestaurants,
        fetchRestaurantById
    }
})