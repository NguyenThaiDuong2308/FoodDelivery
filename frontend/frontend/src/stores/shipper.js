// frontend/src/stores/shipper.js
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { shipperService } from '../services/shipper.service'

export const useShipperStore = defineStore('shipper', () => {
    const shippers = ref([])
    const selectedShipper = ref(null)
    const shipperLocation = ref(null)
    const loading = ref(false)
    const error = ref(null)

    // Fetch all shippers
    async function fetchShippers() {
        loading.value = true
        error.value = null
        try {
            const res = await shipperService.getAll()
            shippers.value = res.shippers
            console.log(shippers.value)
        } catch (err) {
            error.value = err.message || 'Failed to fetch shippers'
            console.error('Error fetching shippers:', err)
        } finally {
            loading.value = false
        }
    }

    // Fetch shipper by ID
    async function fetchShipperById(id) {
        loading.value = true
        error.value = null
        try {
            selectedShipper.value = await shipperService.getById(id)
            return selectedShipper.value
        } catch (err) {
            error.value = err.message || 'Failed to fetch shipper'
            console.error('Error fetching shipper:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // Update shipper location
    async function updateLocation(id, location) {
        loading.value = true
        error.value = null
        try {
            await shipperService.updateLocation(id, location)

            // Update location in local state
            const shipper = shippers.value.find(s => s.id === id)
            if (shipper) {
                shipper.location = location
            }

            if (selectedShipper.value?.id === id) {
                selectedShipper.value.location = location
            }

            shipperLocation.value = location
        } catch (err) {
            error.value = err.message || 'Failed to update location'
            console.error('Error updating location:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // Update shipper status
    async function updateStatus(id, status) {
        loading.value = true
        error.value = null
        try {
            await shipperService.updateStatus(id, status)

            // Update status in local state
            const shipper = shippers.value.find(s => s.id === id)
            if (shipper) {
                shipper.status = status
            }

            if (selectedShipper.value?.id === id) {
                selectedShipper.value.status = status
            }
        } catch (err) {
            error.value = err.message || 'Failed to update status'
            console.error('Error updating status:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // Get shipper location
    async function getLocation(id) {
        loading.value = true
        error.value = null
        try {
            shipperLocation.value = await shipperService.getLocation(id)
            return shipperLocation.value
        } catch (err) {
            error.value = err.message || 'Failed to get location'
            console.error('Error getting location:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // Clear error
    function clearError() {
        error.value = null
    }

    // Clear selected shipper
    function clearSelectedShipper() {
        selectedShipper.value = null
        shipperLocation.value = null
    }

    return {
        // State
        shippers,
        selectedShipper,
        shipperLocation,
        loading,
        error,

        // Actions
        fetchShippers,
        fetchShipperById,
        updateLocation,
        updateStatus,
        getLocation,
        clearError,
        clearSelectedShipper
    }
})