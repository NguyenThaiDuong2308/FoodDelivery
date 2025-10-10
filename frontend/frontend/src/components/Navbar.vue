<template>
  <header class="navbar">
    <div class="container">
      <div class="navbar-brand">
        <Truck :size="32" color="#ea580c" />
        <h1>DeliveryApp</h1>
      </div>

      <nav class="navbar-menu" :class="{ active: mobileMenuOpen }">
        <router-link to="/" @click="closeMobile">
          <Store :size="20" /> Restaurants
        </router-link>
        <router-link to="/orders" @click="closeMobile">
          <Package :size="20" /> Orders
        </router-link>
        <router-link to="/tracking" @click="closeMobile">
          <MapPin :size="20" /> Track
        </router-link>
        <router-link
            v-if="authStore.user?.role === 'restaurant_admin'"
            to="/admin/restaurants"
            @click="closeMobile"
        >
          <Settings :size="20" /> Manage Restaurant
        </router-link>
        <router-link to="/cart" class="cart-link" @click="closeMobile">
          <ShoppingCart :size="24" />
          <span v-if="cartStore.itemCount > 0" class="badge">
            {{ cartStore.itemCount }}
          </span>
        </router-link>
        <router-link to="/profile" class="user-profile" @click="closeMobile">
          <User :size="20" />
          <span>{{ authStore.user?.name }}</span>
          <button @click="handleLogout" class="logout-btn">
            <LogOut :size="20" />
          </button>
        </router-link>
      </nav>

      <button class="mobile-toggle" @click="mobileMenuOpen = !mobileMenuOpen">
        <Menu v-if="!mobileMenuOpen" :size="24" />
        <X v-else :size="24" />
      </button>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Truck, Store, Package, MapPin, ShoppingCart, User, LogOut, Menu, Settings, X } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'

const router = useRouter()
const authStore = useAuthStore()
const cartStore = useCartStore()
const mobileMenuOpen = ref(false)

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}

const closeMobile = () => {
  mobileMenuOpen.value = false
}
</script>

<style scoped>
.navbar {
  background: white;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.navbar-brand h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
}

.navbar-menu {
  display: flex;
  align-items: center;
  gap: 2rem;
}

.navbar-menu a {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #4b5563;
  text-decoration: none;
  transition: color 0.2s;
}

.navbar-menu a:hover,
.navbar-menu a.router-link-active {
  color: #ea580c;
}

.cart-link {
  position: relative;
}

.badge {
  position: absolute;
  top: -8px;
  right: -8px;
  background: #ef4444;
  color: white;
  font-size: 0.75rem;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding-left: 1rem;
  border-left: 1px solid #e5e7eb;
}

.logout-btn {
  background: none;
  border: none;
  color: #ef4444;
  cursor: pointer;
  padding: 0.5rem;
  display: flex;
  align-items: center;
}

.mobile-toggle {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
}

@media (max-width: 768px) {
  .mobile-toggle {
    display: block;
  }

  .navbar-menu {
    display: none;
  }

  .navbar-menu.active {
    display: flex;
    flex-direction: column;
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: white;
    padding: 1rem;
    box-shadow: 0 4px 8px rgba(0,0,0,0.1);
    align-items: flex-start;
  }
}
</style>