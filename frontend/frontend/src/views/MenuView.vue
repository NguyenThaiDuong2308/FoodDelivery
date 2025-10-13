<template>
  <div>
    <Navbar />
    <div class="container">
      <!-- Back Button -->
      <button @click="goBack" class="back-btn">
        <ArrowLeft :size="20" />
        Back to Restaurants
      </button>

      <!-- Loading State -->
      <div v-if="restaurantStore.loading" class="loading">
        <Loader2 :size="48" class="spinner" />
        <p>Loading menu...</p>
      </div>

      <!-- Menu Content -->
      <div v-else-if="restaurantStore.selectedRestaurant">
        <!-- Restaurant Header -->
        <div class="restaurant-header">
          <div class="header-icon">
            <Store :size="40" />
          </div>
          <div class="header-info">
            <h2>{{ restaurantStore.selectedRestaurant.name }}</h2>
            <div class="header-details">
              <span v-if="restaurantStore.selectedRestaurant.address" class="detail-item">
                <MapPin :size="16" />
                {{ restaurantStore.selectedRestaurant.address }}
              </span>
              <span v-if="restaurantStore.selectedRestaurant.phone_number" class="detail-item">
                <Phone :size="16" />
                {{ restaurantStore.selectedRestaurant.phone_number }}
              </span>
            </div>
            <div class="status-badge" :class="restaurantStore.selectedRestaurant.status">
              {{ restaurantStore.selectedRestaurant.status === 'open' ? '🟢 Open' : '🔴 Closed' }}
            </div>
          </div>
        </div>

        <!-- Menu Section -->
        <div class="menu-section">
          <div class="menu-header">
            <h3 class="menu-title">
              <UtensilsCrossed :size="24" />
              Menu Items
            </h3>
            <span class="item-count">{{ menuItems.length }} items</span>
          </div>

          <!-- Empty Menu -->
          <div v-if="menuItems.length === 0" class="empty-menu">
            <UtensilsCrossed :size="64" />
            <h4>No menu items available</h4>
            <p>This restaurant hasn't added any dishes yet.</p>
          </div>

          <!-- Menu Grid -->
          <div v-else class="menu-grid">
            <div
                v-for="item in menuItems"
                :key="item.id"
                class="menu-item"
                :class="{ unavailable: !item.available }"
            >
              <div class="item-badge" v-if="!item.available">
                <AlertCircle :size="16" />
                Unavailable
              </div>

              <div class="item-content">
                <h4>{{ item.name }}</h4>
                <p class="item-description">{{ item.description || 'No description available' }}</p>

                <div class="item-footer">
                  <span class="price">${{ parseFloat(item.price).toFixed(2) }}</span>

                  <!-- Add to Cart Button: Only for customers -->
                  <button
                      v-if="authStore.user?.role === 'customer' && item.available"
                      @click="addToCart(item)"
                      class="btn-add"
                  >
                    <ShoppingCart :size="18" />
                    Add to Cart
                  </button>

                  <!-- Unavailable Text -->
                  <span v-else-if="!item.available" class="unavailable-text">
                    Out of Stock
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Cart Summary (Fixed Bottom) - Only for customers -->
        <div v-if="authStore.user?.role === 'customer' && cartStore.items.length > 0" class="cart-summary">
          <div class="cart-info">
            <ShoppingCart :size="20" />
            <span>{{ cartStore.items.length }} items</span>
            <span class="cart-total">${{ cartTotal }}</span>
          </div>
          <button @click="goToCart" class="btn-checkout">
            View Cart
            <ArrowRight :size="18" />
          </button>
        </div>
      </div>

      <!-- Error State -->
      <div v-else class="error-state">
        <AlertCircle :size="64" />
        <h3>Restaurant not found</h3>
        <button @click="goBack" class="btn-back">Go Back</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Store,
  MapPin,
  Phone,
  UtensilsCrossed,
  ShoppingCart,
  ArrowLeft,
  ArrowRight,
  Loader2,
  AlertCircle
} from 'lucide-vue-next'
import Navbar from '../components/Navbar.vue'
import { useRestaurantStore } from '../stores/restaurant'
import { useCartStore } from '../stores/cart'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const restaurantStore = useRestaurantStore()
const cartStore = useCartStore()
const authStore = useAuthStore()

// Biến riêng để lưu menu items
const menuItems = ref([])

onMounted(async () => {
  const id = route.params.id

  // Fetch restaurant info
  await restaurantStore.fetchRestaurantById(id)

  // Fetch menu items và lưu vào biến local
  const items = await restaurantStore.fetchMenu(id)
  menuItems.value = items || []

  // Set restaurant for cart (only for customers)
  if (authStore.user?.role === 'customer') {
    cartStore.setRestaurant(id)
  }
})

const cartTotal = computed(() => {
  const total = cartStore.items.reduce((sum, item) => {
    return sum + (parseFloat(item.price) * item.quantity)
  }, 0)
  return total.toFixed(2)
})

const goBack = () => {
  router.push('/')
}

const goToCart = () => {
  router.push('/cart')
}

const addToCart = (item) => {
  if (!item.available) return
  if (authStore.user?.role !== 'customer') return

  cartStore.addItem(item)

  // Show success message
  showNotification(`✅ ${item.name} added to cart!`)
}

// Simple notification (you can replace with a toast library)
const showNotification = (message) => {
  const notification = document.createElement('div')
  notification.className = 'notification'
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
  }, 2000)
}
</script>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem 8rem;
}

/* Back Button */
.back-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: white;
  border: 2px solid #e5e7eb;
  color: #ea580c;
  font-size: 1rem;
  font-weight: 600;
  padding: 0.75rem 1.25rem;
  border-radius: 0.75rem;
  cursor: pointer;
  margin-bottom: 1.5rem;
  transition: all 0.2s;
}

.back-btn:hover {
  background: #fff7ed;
  border-color: #fed7aa;
}

/* Loading State */
.loading {
  text-align: center;
  padding: 4rem 2rem;
  color: #6b7280;
}

.spinner {
  animation: spin 1s linear infinite;
  color: #ea580c;
  margin-bottom: 1rem;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Restaurant Header */
.restaurant-header {
  display: flex;
  gap: 1.5rem;
  background: white;
  padding: 2rem;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  margin-bottom: 2rem;
}

.header-icon {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #fb923c 0%, #ea580c 100%);
  border-radius: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.header-info {
  flex: 1;
}

.header-info h2 {
  font-size: 2rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.75rem 0;
}

.header-details {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #6b7280;
  font-size: 0.9rem;
}

.status-badge {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 700;
}

.status-badge.open {
  background: #dcfce7;
  color: #15803d;
}

.status-badge.closed {
  background: #fee2e2;
  color: #991b1b;
}

/* Menu Section */
.menu-section {
  margin-bottom: 2rem;
}

.menu-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.menu-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1.75rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}

.item-count {
  background: #fff7ed;
  color: #ea580c;
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 600;
}

/* Empty Menu */
.empty-menu {
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
}

.empty-menu svg {
  color: #d1d5db;
  margin-bottom: 1rem;
}

.empty-menu h4 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.5rem 0;
}

.empty-menu p {
  color: #6b7280;
  margin: 0;
}

/* Menu Grid */
.menu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.menu-item {
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  overflow: hidden;
  transition: all 0.3s;
  position: relative;
}

.menu-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 20px rgba(0,0,0,0.15);
}

.menu-item.unavailable {
  opacity: 0.6;
}

.item-badge {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: #fee2e2;
  color: #991b1b;
  padding: 0.375rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  z-index: 1;
}

.item-content {
  padding: 1.5rem;
}

.menu-item h4 {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.5rem 0;
}

.item-description {
  color: #6b7280;
  font-size: 0.9rem;
  margin: 0 0 1.25rem 0;
  line-height: 1.5;
  min-height: 3rem;
}

.item-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.price {
  font-size: 1.75rem;
  font-weight: 700;
  color: #ea580c;
}

.btn-add {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #ea580c;
  color: white;
  border: none;
  padding: 0.75rem 1.25rem;
  border-radius: 0.75rem;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
  transition: all 0.2s;
}

.btn-add:hover {
  background: #c2410c;
  transform: scale(1.05);
}

.unavailable-text {
  color: #ef4444;
  font-size: 0.875rem;
  font-weight: 600;
}

/* Cart Summary (Fixed Bottom) */
.cart-summary {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  border-top: 2px solid #e5e7eb;
  padding: 1.25rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 -4px 12px rgba(0,0,0,0.1);
  z-index: 100;
}

.cart-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  color: #1f2937;
  font-weight: 600;
}

.cart-total {
  font-size: 1.5rem;
  color: #ea580c;
  font-weight: 700;
}

.btn-checkout {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #ea580c;
  color: white;
  border: none;
  padding: 1rem 2rem;
  border-radius: 0.75rem;
  cursor: pointer;
  font-weight: 700;
  font-size: 1rem;
  transition: all 0.2s;
}

.btn-checkout:hover {
  background: #c2410c;
  transform: scale(1.05);
}

/* Error State */
.error-state {
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.error-state svg {
  color: #ef4444;
  margin-bottom: 1rem;
}

.error-state h3 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 1.5rem 0;
}

.btn-back {
  background: #ea580c;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 0.75rem;
  cursor: pointer;
  font-weight: 600;
}

.btn-back:hover {
  background: #c2410c;
}

/* Notification */
:global(.notification) {
  position: fixed;
  top: 2rem;
  right: 2rem;
  background: #15803d;
  color: white;
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

/* Responsive */
@media (max-width: 768px) {
  .container {
    padding: 1rem 0.5rem 8rem;
  }

  .restaurant-header {
    flex-direction: column;
    text-align: center;
  }

  .header-icon {
    margin: 0 auto;
  }

  .header-details {
    align-items: center;
  }

  .menu-grid {
    grid-template-columns: 1fr;
  }

  .menu-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .cart-summary {
    flex-direction: column;
    gap: 0.75rem;
  }

  .btn-checkout {
    width: 100%;
  }
}
</style>