<template>
  <div>
    <Navbar />
    <div class="container">
      <h2 class="page-title">Your Orders</h2>

      <div v-if="orderStore.loading" class="loading">
        Loading orders...
      </div>

      <div v-else-if="orderStore.orders.length === 0" class="empty-state">
        <Package :size="64" color="#d1d5db" />
        <p>No orders yet</p>
        <router-link to="/" class="btn btn-primary">Start Ordering</router-link>
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
            </div>
            <span class="status-badge" :class="order.status">
              {{ order.status }}
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
            <Clock :size="16" />
            <span>Placed on {{ formatDate(order.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { Package, Clock } from 'lucide-vue-next'
import Navbar from '../components/Navbar.vue'
import { useOrderStore } from '../stores/order'
import { useAuthStore } from '../stores/auth'

const orderStore = useOrderStore()
const authStore = useAuthStore()

onMounted(() => {
  orderStore.fetchCustomerOrders(authStore.user.id)
})

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
  margin-bottom: 0.25rem;
}

.restaurant-info {
  color: #6b7280;
  font-size: 0.875rem;
}

.status-badge {
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: capitalize;
}

.status-badge.completed {
  background: #dcfce7;
  color: #15803d;
}

.status-badge.pending {
  background: #fef3c7;
  color: #92400e;
}

.status-badge.processing {
  background: #dbeafe;
  color: #1e40af;
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
  align-items: center;
  gap: 0.5rem;
  color: #9ca3af;
  font-size: 0.875rem;
  padding-top: 1rem;
  border-top: 1px solid #e5e7eb;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  text-decoration: none;
  display: inline-block;
}

.btn-primary {
  background: #ea580c;
  color: white;
}

.btn-primary:hover {
  background: #c2410c;
}
</style>