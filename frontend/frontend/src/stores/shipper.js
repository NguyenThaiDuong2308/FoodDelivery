import { defineStore } from 'pinia'
import { ref } from 'vue'
import { shipperService } from '../services/shipper.service'

export const useShipperStore = defineStore('shipper', () => {
    const shippers = ref([])
    const loading = ref(false)

    async function fetchShippers() {
        loading.value = true
        try {
            shippers.value = await shipperService.getAll()
        } finally {
            loading.value = false
        }
    }

    async function updateLocation(id, location) {
        await shipperService.updateLocation(id, location)
        await fetchShippers()
    }

    return {
        shippers,
        loading,
        fetchShippers,
        updateLocation
    }
})