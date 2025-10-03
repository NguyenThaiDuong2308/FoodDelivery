<template>
  <div>
    <Navbar />
    <div class="container">
      <button @click="goBack" class="back-btn">
        ← Back to Restaurants
      </button>

      <div v-if="restaurantStore.loading" class="loading">
        Loading menu...
      </div>

      <div v-else-if="restaurantStore.selectedRestaurant">
        <div class="restaurant-header">
          <h2>{{ restaurantStore.selectedRestaurant.name }}</h2>
          <p>{{ restaurantStore.selectedRestaurant.address }}</p>
        </div>

        <h3 class="menu-title">Menu</h3>

        <div class="menu-grid">
          <div
              v-for="item in restaurantStore.selectedRestaurant.menu_items"
              :key="item.id"
              class="menu-item"
          >
            <h4>{{ item.name }}</h4>
            <p class="item-description">{{ item.description }}</p>
            <div class="item-footer">
              <span class="price">${{ item.price }}</span>
              <button
                  v-if="item.available"
                  @click="addToCart(item)"
                  class="btn btn-add"
              >
                Add to Cart
              </button>
              <span v-else class="unavailable">Unavailable</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Navbar from '../components/Navbar.vue'
import { useRestaurantStore } from '../stores/restaurant'
import { useCartStore } from '../stores/cart'

const route = useRoute()
const router = useRouter()
const restaurantStore = useRestaurantStore()
const cartStore = useCartStore()

onMounted(() => {
  const id = route.params.id
  restaurantStore.fetchRestaurantById(id)
  cartStore.setRestaurant(id)
})

const goBack = () => {
  router.push('/')
}

const addToCart = (item) => {
  cartStore.addItem(item)
  alert(`${item.name} added to cart!`)
}
</script>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.back-btn {
  background: none;
  border: none;
  color: #ea580c;
  font-size: 1rem;
  cursor: pointer;
  margin-bottom: 1.5rem;
}

.back-btn:hover {
  text-decoration: underline;
}

.loading {
  text-align: center;
  padding: 3rem;
  color: #6b7280;
}

.restaurant-header {
  background: white;
  padding: 2rem;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  margin-bottom: 2rem;
}

.restaurant-header h2 {
  font-size: 2rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 0.5rem;
}

.restaurant-header p {
  color: #6b7280;
}

.menu-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 1.5rem;
}

.menu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.5rem;
}

.menu-item {
  background: white;
  padding: 1.5rem;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.menu-item h4 {
  font-size: 1.125rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 0.5rem;
}

.item-description {
  color: #6b7280;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.item-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price {
  font-size: 1.5rem;
  font-weight: 700;
  color: #ea580c;
}

.btn-add {
  background: #ea580c;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  cursor: pointer;
  font-weight: 600;
}

.btn-add:hover {
  background: #c2410c;
}

.unavailable {
  color: #ef4444;
  font-size: 0.875rem;
}
</style>