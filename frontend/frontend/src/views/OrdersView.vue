<template>
  <div>
    <Navbar />
    <div class="container">
      <h2 class="page-title">{{ getPageTitle() }}</h2>

      <div v-if="orderStore.loading" class="loading">
        Loading orders...
      </div>

      <div v-else-if="!orderStore.orders || orderStore.orders.length === 0" class="empty-state">
        <Package :size="64" color="#d1d5db" />
        <p>No orders yet</p>
        <router-link v-if="authStore.user?.role === 'customer'" to="/" class="btn btn-primary">
          Start Ordering
        </router-link>
      </div>

      <div v-else class="orders-list">
        <div
            v-for="order in orderStore.orders"
            :key="order.id"
            class="order-card"
        >
          <div class="order-header">
            <div>
              <h3>Order #{{ order.id }}</h3>
              <p class="restaurant-info">Restaurant ID: {{ order.restaurant_id }}</p>
              <p v-if="authStore.user?.role === 'restaurant_admin'" class="customer-info">
                Customer ID: {{ order.customer_id }}
              </p>
              <p v-if="authStore.user?.role === 'shipper' && order.shipper_id" class="shipper-info">
                Assigned to you
              </p>
              <p class="price-info">Total: ${{ order.total_price }}</p>
            </div>
            <span class="status-badge" :class="order.status">
              {{ getStatusLabel(order.status) }}
            </span>
          </div>

          <div class="order-items">
            <div
                v-for="(item, index) in order.order_items"
                :key="index"
                class="order-item"
            >
              <span>Menu Item #{{ item.menu_item_id }}</span>
              <span class="quantity">Qty: {{ item.quantity }}</span>
            </div>
          </div>

          <div class="order-footer">
            <div class="footer-info">
              <Clock :size="16" />
              <span>Placed on {{ formatDate(order.created_at) }}</span>
            </div>
            <div class="footer-actions">
              <button @click="viewOrderDetails(order.id)" class="btn btn-secondary">
                <Eye :size="16" />
                View Details
              </button>

              <!-- Customer: Mark as Completed -->
              <button
                  v-if="authStore.user?.role === 'customer' && order.status === 'delivering'"
                  @click="markAsCompleted(order.id)"
                  class="btn btn-success"
                  :disabled="updatingOrder === order.id"
              >
                <CheckCircle :size="16" />
                {{ updatingOrder === order.id ? 'Updating...' : 'Mark Completed' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Order Details Modal -->
    <div v-if="showDetailsModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Order #{{ selectedOrder?.id }} Details</h3>
          <button @click="closeModal" class="btn-close">
            <X :size="24" />
          </button>
        </div>

        <div v-if="loadingDetails" class="modal-loading">
          Loading details...
        </div>

        <div v-else-if="selectedOrder" class="modal-body">
          <!-- Status -->
          <div class="detail-section">
            <h4>Status</h4>
            <span class="status-badge large" :class="selectedOrder.status">
              {{ getStatusLabel(selectedOrder.status) }}
            </span>
            <p class="status-description">{{ getStatusDescription(selectedOrder.status) }}</p>
          </div>

          <!-- Customer Info (for restaurant_admin) -->
          <div v-if="authStore.user?.role === 'restaurant_admin'" class="detail-section">
            <h4>Customer</h4>
            <p>Customer ID: {{ selectedOrder.customer_id }}</p>
          </div>

          <!-- Restaurant Info -->
          <div class="detail-section">
            <h4>Restaurant</h4>
            <p>Restaurant ID: {{ selectedOrder.restaurant_id }}</p>
          </div>

          <!-- Shipper Info -->
          <div v-if="selectedOrder.shipper_id" class="detail-section">
            <h4>Shipper</h4>
            <p>Shipper ID: {{ selectedOrder.shipper_id }}</p>
          </div>

          <!-- Order Items -->
          <div class="detail-section">
            <h4>Order Items</h4>
            <div class="detail-items">
              <div
                  v-for="(item, index) in selectedOrder.order_items"
                  :key="index"
                  class="detail-item"
              >
                <span>Menu Item #{{ item.menu_item_id }}</span>
                <span>Qty: {{ item.quantity }}</span>
              </div>
            </div>
          </div>

          <!-- Total -->
          <div class="detail-section">
            <h4>Total Price</h4>
            <p class="total-price">${{ selectedOrder.total_price }}</p>
          </div>

          <!-- Date -->
          <div class="detail-section">
            <h4>Order Date</h4>
            <p>{{ formatDate(selectedOrder.created_at) }}</p>
          </div>

          <!-- Actions -->
          <div class="modal-actions">
            <!-- Customer: Mark as Completed -->
            <button
                v-if="authStore.user?.role === 'customer' && selectedOrder.status === 'delivering'"
                @click="markAsCompleted(selectedOrder.id)"
                class="btn btn-success btn-large"
                :disabled="updatingOrder === selectedOrder.id"
            >
              <CheckCircle :size="18" />
              {{ updatingOrder === selectedOrder.id ? 'Updating...' : 'Mark as Completed' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Package, Clock, Eye, CheckCircle, X } from 'lucide-vue-next'
import Navbar from '../components/Navbar.vue'
import { useOrderStore } from '../stores/order'
import { useAuthStore } from '../stores/auth'
import { useRestaurantStore } from '../stores/restaurant'
import { useShipperStore } from '../stores/shipper'

const orderStore = useOrderStore()
const authStore = useAuthStore()
const restaurantStore = useRestaurantStore()
const shipperStore = useShipperStore()

const showDetailsModal = ref(false)
const selectedOrder = ref(null)
const loadingDetails = ref(false)
const updatingOrder = ref(null)

onMounted(async () => {
  const user = authStore.user
  if (!user) return

  try {
    switch (user.role) {
      case 'customer':
        await orderStore.fetchCustomerOrders(user.id)
        break

      case 'restaurant_admin':
        await restaurantStore.fetchRestaurants()
        const userRestaurant = restaurantStore.restaurants.find(r => r.manager_id === user.id)
        if (userRestaurant) {
          await orderStore.fetchRestaurantOrders(userRestaurant.id)
        }
        break

      case 'shipper':
        // Fetch shipper's orders
        // First get shipper record for this user
        await shipperStore.fetchShippers()
        const shipperRecord = shipperStore.shippers.find(s => s.user_id === user.id)
        if (shipperRecord) {
          await orderStore.fetchShipperOrders(shipperRecord.id)
        }
        break

      default:
        console.warn('Unknown user role:', user.role)
    }
  } catch (error) {
    console.error('Failed to fetch orders:', error)
    showNotification('Failed to load orders', 'error')
  }
})

const getPageTitle = () => {
  const role = authStore.user?.role
  const titles = {
    customer: 'Your Orders',
    restaurant_admin: 'Restaurant Orders',
    shipper: 'Your Deliveries'
  }
  return titles[role] || 'Orders'
}

const formatDate = (date) => {
  if (!date) return 'Recently'
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusLabel = (status) => {
  const labels = {
    created: 'Created',
    cancelled: 'Cancelled',
    delivering: 'Delivering',
    completed: 'Completed'
  }
  return labels[status] || status
}

const getStatusDescription = (status) => {
  const descriptions = {
    created: 'Your order has been created and is waiting for shipper assignment.',
    cancelled: 'No shipper was found. Order has been cancelled.',
    delivering: 'Your order is on the way! Shipper has been assigned.',
    completed: 'Order has been completed successfully.'
  }
  return descriptions[status] || ''
}

const viewOrderDetails = async (orderId) => {
  showDetailsModal.value = true
  loadingDetails.value = true

  try {
    const orderDetails = await orderStore.fetchOrderDetails(orderId)
    selectedOrder.value = orderDetails
  } catch (error) {
    console.error('Failed to fetch order details:', error)
    showNotification('Failed to load order details', 'error')
  } finally {
    loadingDetails.value = false
  }
}

const closeModal = () => {
  showDetailsModal.value = false
  selectedOrder.value = null
}

const markAsCompleted = async (orderId) => {
  if (updatingOrder.value) return

  if (!confirm('Mark this order as completed?')) return

  updatingOrder.value = orderId

  try {
    await orderStore.updateOrderStatus(orderId, 'completed')
    showNotification('✅ Order marked as completed!', 'success')

    // Refresh orders list based on role
    const user = authStore.user
    if (user.role === 'customer') {
      await orderStore.fetchCustomerOrders(user.id)
    }

    // Update modal if open
    if (selectedOrder.value?.id === orderId) {
      selectedOrder.value.status = 'completed'
    }
  } catch (error) {
    console.error('Failed to update order status:', error)
    showNotification('❌ Failed to update order status', 'error')
  } finally {
    updatingOrder.value = null
  }
}

const showNotification = (message, type = 'success') => {
  const notification = document.createElement('div')
  notification.className = `notification ${type}`
  notification.textContent = message
  document.body.appendChild(notification)

  setTimeout(() => {
    notification.classList.add('show')
  }, 10)

  setTimeout(() => {
    notification.classList.remove('show')
    setTimeout(() => {
      document.body.removeChild(notification)
    }, 300)
  }, 3000)
}
</script>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 2rem;
}

.loading {
  text-align: center;
  padding: 3rem;
  color: #6b7280;
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.empty-state p {
  color: #6b7280;
  margin: 1rem 0 2rem;
  font-size: 1.125rem;
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.order-card {
  background: white;
  border-radius: 1rem;
  padding: 1.5rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  transition: all 0.3s;
}

.order-card:hover {
  box-shadow: 0 6px 16px rgba(0,0,0,0.15);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e5e7eb;
}

.order-header h3 {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.25rem 0;
}

.restaurant-info,
.customer-info,
.shipper-info {
  color: #6b7280;
  font-size: 0.875rem;
  margin: 0.25rem 0;
}

.customer-info {
  color: #4f46e5;
}

.shipper-info {
  color: #15803d;
  font-weight: 600;
}

.price-info {
  color: #ea580c;
  font-size: 1.125rem;
  font-weight: 700;
  margin: 0.5rem 0 0 0;
}

.status-badge {
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: capitalize;
}

.status-badge.large {
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
}

.status-badge.completed {
  background: #dcfce7;
  color: #15803d;
}

.status-badge.created {
  background: #dbeafe;
  color: #1e40af;
}

.status-badge.delivering {
  background: #fef3c7;
  color: #92400e;
}

.status-badge.cancelled {
  background: #fee2e2;
  color: #991b1b;
}

.order-items {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.order-item {
  display: flex;
  justify-content: space-between;
  color: #4b5563;
  padding: 0.5rem 0;
}

.quantity {
  font-weight: 600;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 1rem;
  border-top: 1px solid #e5e7eb;
  gap: 1rem;
  flex-wrap: wrap;
}

.footer-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #9ca3af;
  font-size: 0.875rem;
}

.footer-actions {
  display: flex;
  gap: 0.75rem;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: #ea580c;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #c2410c;
}

.btn-secondary {
  background: #f3f4f6;
  color: #1f2937;
  border: 1px solid #d1d5db;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.btn-success {
  background: #22c55e;
  color: white;
}

.btn-success:hover:not(:disabled) {
  background: #16a34a;
}

.btn-large {
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  width: 100%;
  justify-content: center;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-content {
  background: white;
  border-radius: 1rem;
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}

.btn-close {
  padding: 0.5rem;
  background: transparent;
  border: none;
  cursor: pointer;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 0.5rem;
  transition: all 0.2s;
}

.btn-close:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.modal-loading {
  padding: 3rem;
  text-align: center;
  color: #6b7280;
}

.modal-body {
  padding: 1.5rem;
}

.detail-section {
  margin-bottom: 1.5rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.detail-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.detail-section h4 {
  font-size: 0.875rem;
  font-weight: 700;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0 0 0.75rem 0;
}

.detail-section p {
  color: #1f2937;
  font-size: 1rem;
  margin: 0.5rem 0;
}

.status-description {
  color: #6b7280 !important;
  font-size: 0.875rem !important;
  margin-top: 0.75rem !important;
  font-style: italic;
}

.detail-items {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem;
  background: #f9fafb;
  border-radius: 0.5rem;
  color: #4b5563;
}

.total-price {
  font-size: 1.5rem !important;
  font-weight: 700 !important;
  color: #ea580c !important;
}

.modal-actions {
  margin-top: 1.5rem;
}

/* Notification */
:global(.notification) {
  position: fixed;
  top: 2rem;
  right: 2rem;
  padding: 1rem 1.5rem;
  border-radius: 0.75rem;
  box-shadow: 0 8px 20px rgba(0,0,0,0.2);
  font-weight: 600;
  z-index: 9999;
  opacity: 0;
  transform: translateY(-20px);
  transition: all 0.3s;
}

:global(.notification.show) {
  opacity: 1;
  transform: translateY(0);
}

:global(.notification.success) {
  background: #15803d;
  color: white;
}

:global(.notification.error) {
  background: #dc2626;
  color: white;
}

/* Responsive */
@media (max-width: 768px) {
  .container {
    padding: 1rem 0.5rem;
  }

  .order-header {
    flex-direction: column;
    gap: 1rem;
  }

  .status-badge {
    align-self: flex-start;
  }

  .order-footer {
    flex-direction: column;
    align-items: flex-start;
  }

  .footer-actions {
    width: 100%;
    flex-direction: column;
  }

  .footer-actions .btn {
    width: 100%;
    justify-content: center;
  }

  .modal-content {
    max-height: 95vh;
  }
}
</style>