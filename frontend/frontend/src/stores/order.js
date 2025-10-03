import { defineStore } from 'pinia'
import { ref } from 'vue'
import { orderService } from '../services/order.service'

export const useOrderStore = defineStore('order', () => {
    const orders = ref([])
    const loading = ref(false)

    async function fetchCustomerOrders(customerId) {
        loading.value = true
        try {
            orders.value = await orderService.getCustomerOrders(customerId)
        } finally {
            loading.value = false
        }
    }

    async function createOrder(orderData) {
        loading.value = true
        try {
            const newOrder = await orderService.create(orderData)
            orders.value.unshift(newOrder)
            return newOrder
        } finally {
            loading.value = false
        }
    }

    return {
        orders,
        loading,
        fetchCustomerOrders,
        createOrder
    }
})