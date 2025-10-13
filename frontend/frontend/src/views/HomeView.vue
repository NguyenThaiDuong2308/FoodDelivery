<template>
  <div>
    <Navbar />
    <div class="container">
      <h2 class="page-title">Restaurants</h2>

      <div v-if="restaurantStore.loading" class="loading">
        <div class="spinner"></div>
        Loading restaurants...
      </div>

      <div v-else-if="restaurantStore.error" class="error-state">
        <AlertCircle :size="64" />
        <p>{{ restaurantStore.error }}</p>
        <button @click="restaurantStore.fetchRestaurants()" class="btn-retry">
          Try Again
        </button>
      </div>

      <div v-else-if="restaurantStore.restaurants.length === 0" class="empty-state">
        <Store :size="64" />
        <p>No restaurants available</p>
      </div>

      <div v-else class="restaurant-grid">
        <div
            v-for="restaurant in restaurantStore.restaurants"
            :key="restaurant.id"
            class="restaurant-card"
            @click="goToMenu(restaurant.id)"
        >
          <div class="restaurant-image">
            <Store :size="64" color="white" />
          </div>
          <div class="restaurant-info">
            <div class="restaurant-header">
              <h3>{{ restaurant.name }}</h3>
              <span class="status-badge" :class="restaurant.status">
                {{ restaurant.status }}
              </span>
            </div>

            <p class="description">{{ restaurant.description || 'No description' }}</p>

            <div class="restaurant-details">
              <p class="address">
                <MapPin :size="16" />
                {{ restaurant.address || 'No address' }}
              </p>

              <p v-if="restaurant.phone_number" class="phone">
                <Phone :size="16" />
                {{ restaurant.phone_number }}
              </p>

              <p v-if="restaurant.email" class="email">
                <Mail :size="16" />
                {{ restaurant.email }}
              </p>
            </div>

            <!-- Show menu item count if available -->
            <div v-if="restaurant.menu_items && restaurant.menu_items.length > 0" class="menu-count">
              <UtensilsCrossed :size="16" />
              {{ restaurant.menu_items.length }} menu items
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Store, MapPin, Phone, Mail, UtensilsCrossed, AlertCircle } from 'lucide-vue-next'
import Navbar from '../components/Navbar.vue'
import { useRestaurantStore } from '../stores/restaurant'

const router = useRouter()
const restaurantStore = useRestaurantStore()

onMounted(() => {
  restaurantStore.fetchRestaurants()
})

const goToMenu = (id) => {
  router.push(`/menu/${id}`)
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
  font-size: 1.125rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #f3f4f6;
  border-top-color: #ea580c;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

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

.error-state p {
  color: #6b7280;
  margin: 0 0 1.5rem 0;
  font-size: 1.125rem;
}

.btn-retry {
  padding: 0.75rem 1.5rem;
  background: #ea580c;
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-retry:hover {
  background: #c2410c;
  transform: translateY(-1px);
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.empty-state svg {
  color: #d1d5db;
  margin-bottom: 1rem;
}

.empty-state p {
  color: #6b7280;
  margin: 0;
  font-size: 1.125rem;
}

.restaurant-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
}

.restaurant-card {
  background: white;
  border-radius: 1rem;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  cursor: pointer;
  transition: all 0.3s;
}

.restaurant-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.15);
}

.restaurant-image {
  height: 160px;
  background: linear-gradient(135deg, #fb923c 0%, #ef4444 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.restaurant-info {
  padding: 1.5rem;
}

.restaurant-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.restaurant-info h3 {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
  flex: 1;
}

.description {
  color: #6b7280;
  font-size: 0.875rem;
  margin-bottom: 1rem;
  line-height: 1.5;
}

.restaurant-details {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.address,
.phone,
.email {
  color: #6b7280;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0;
}

.address svg,
.phone svg,
.email svg {
  flex-shrink: 0;
  color: #9ca3af;
}

.menu-count {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: #fef3c7;
  color: #92400e;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 600;
  margin-top: 0.75rem;
}

.menu-count svg {
  flex-shrink: 0;
}

.status-badge {
  display: inline-block;
  padding: 0.375rem 0.875rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: capitalize;
  white-space: nowrap;
}

.status-badge.open {
  background: #dcfce7;
  color: #15803d;
}

.status-badge.closed {
  background: #fee2e2;
  color: #991b1b;
}

@media (max-width: 768px) {
  .container {
    padding: 1rem 0.5rem;
  }

  .restaurant-grid {
    grid-template-columns: 1fr;
  }

  .restaurant-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>