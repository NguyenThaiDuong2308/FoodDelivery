import { defineStore } from 'pinia'
import { ref } from 'vue'
import { orderService } from '../services/order.service'

export const useOrderStore = defineStore('order', () => {
    const orders = ref([])
    const currentOrder = ref(null)
    const loading = ref(false)
    const error = ref(null)

    // Fetch all orders for a customer
    async function fetchCustomerOrders(customerId) {
        loading.value = true
        error.value = null
        try {
            orders.value = await orderService.getCustomerOrders(customerId)
            return orders.value
        } catch (err) {
            error.value = err.message || 'Failed to fetch orders'
            console.error('Error fetching customer orders:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    async function fetchRestaurantOrders(restaurantId) {
        loading.value = true
        error.value = null
        try {
            orders.value = await orderService.getRestaurantOrders(restaurantId)
            return orders.value
        } catch (err) {
            error.value = err.message || 'Failed to fetch restaurant orders'
            console.error('Error fetching restaurant orders:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // ✅ Fetch shipper orders (NEW)
    async function fetchShipperOrders(shipperId) {
        loading.value = true
        error.value = null
        try {
            orders.value = await orderService.getShipperOrders(shipperId)
            return orders.value
        } catch (err) {
            error.value = err.message || 'Failed to fetch shipper orders'
            console.error('Error fetching shipper orders:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // Fetch order details by ID
    async function fetchOrderDetails(orderId) {
        loading.value = true
        error.value = null
        try {
            const orderDetails = await orderService.getOrderDetails(orderId)
            currentOrder.value = orderDetails

            // Update order in the orders list if exists
            const index = orders.value.findIndex(o => o.id === orderId)
            if (index !== -1) {
                orders.value[index] = orderDetails
            }

            return orderDetails
        } catch (err) {
            error.value = err.message || 'Failed to fetch order details'
            console.error('Error fetching order details:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // Create a new order
    async function createOrder(orderData) {
        loading.value = true
        error.value = null
        try {
            const newOrder = await orderService.create(orderData)

            // Add to beginning of orders list
            orders.value.unshift(newOrder)
            currentOrder.value = newOrder

            return newOrder
        } catch (err) {
            error.value = err.message || 'Failed to create order'
            console.error('Error creating order:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // Update order status
    async function updateOrderStatus(orderId, status) {
        loading.value = true
        error.value = null
        try {
            const response = await orderService.updateStatus(orderId, status)

            // Update order in the orders list
            const index = orders.value.findIndex(o => o.id === orderId)
            if (index !== -1) {
                orders.value[index].status = status
            }

            // Update currentOrder if it's the same order
            if (currentOrder.value?.id === orderId) {
                currentOrder.value.status = status
            }

            return response
        } catch (err) {
            error.value = err.message || 'Failed to update order status'
            console.error('Error updating order status:', err)
            throw err
        } finally {
            loading.value = false
        }
    }

    // Clear current order
    function clearCurrentOrder() {
        currentOrder.value = null
    }

    // Clear all orders
    function clearOrders() {
        orders.value = []
        currentOrder.value = null
    }

    // Get order by ID from local state
    function getOrderById(orderId) {
        return orders.value.find(o => o.id === orderId)
    }

    // Get orders by status
    function getOrdersByStatus(status) {
        return orders.value.filter(o => o.status === status)
    }

    return {
        // State
        orders,
        currentOrder,
        loading,
        error,

        // Actions
        fetchCustomerOrders,
        fetchOrderDetails,
        createOrder,
        updateOrderStatus,
        clearCurrentOrder,
        clearOrders,
        fetchRestaurantOrders,
        fetchShipperOrders,

        // Getters
        getOrderById,
        getOrdersByStatus
    }
})