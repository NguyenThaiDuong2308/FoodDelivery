<template>
  <div>
    <Navbar />
    <div class="container">
      <h2 class="page-title">Restaurants</h2>

      <div v-if="restaurantStore.loading" class="loading">
        Loading restaurants...
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
            <h3>{{ restaurant.name }}</h3>
            <p class="description">{{ restaurant.description || 'No description' }}</p>
            <p class="address">
              <MapPin :size="16" />
              {{ restaurant.address || 'No address' }}
            </p>
            <div class="status-badge" :class="restaurant.status">
              {{ restaurant.status }}
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
import { Store, MapPin } from 'lucide-vue-next'
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
}

.restaurant-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
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

.restaurant-info h3 {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 0.5rem;
}

.description {
  color: #6b7280;
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
}

.address {
  color: #9ca3af;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  margin-bottom: 1rem;
}

.status-badge {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: capitalize;
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
  .restaurant-grid {
    grid-template-columns: 1fr;
  }
}
</style>