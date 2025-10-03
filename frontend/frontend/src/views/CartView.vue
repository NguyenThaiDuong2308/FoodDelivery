<template>
  <div>
    <Navbar />
    <div class="container">
      <h2 class="page-title">Your Cart</h2>

      <div v-if="cartStore.items.length === 0" class="empty-cart">
        <ShoppingCart :size="64" color="#d1d5db" />
        <p>Your cart is empty</p>
        <router-link to="/" class="btn btn-primary">Browse Restaurants</router-link>
      </div>

      <div v-else class="cart-content">
        <div class="cart-items">
          <div
              v-for="(item, index) in cartStore.items"
              :key="index"
              class="cart-item"
          >
            <div class="item-info">
              <h4>{{ item.name }}</h4>
              <p class="price">${{ item.price }} each</p>
            </div>
            <div class="item-actions">
              <div class="quantity-control">
                <button @click="decreaseQty(index)">-</button>
                <span>{{ item.quantity }}</span>
                <button @click="increaseQty(index)">+</button>
              </div>
              <button @click="removeItem(index)" class="remove-btn">
                <Trash2 :size="20" />
              </button>
            </div>
          </div>
        </div>

        <div class="cart-summary">
          <h3>Order Summary</h3>
          <div class="summary-row">
            <span>Subtotal</span>
            <span>${{ cartStore.totalAmount.toFixed(2) }}</span>
          </div>
          <div class="summary-row">
            <span>Delivery Fee</span>
            <span>$5.00</span>
          </div>
          <div class="summary-row total">
            <span>Total</span>
            <span>${{ (cartStore.totalAmount + 5).toFixed(2) }}</span>
          </div>
          <button @click="placeOrder" class="btn btn-primary" :disabled="loading">
            {{ loading ? 'Placing Order...' : 'Place Order' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ShoppingCart, Trash2 } from 'lucide-vue-next'
import Navbar from '../components/Navbar.vue'
import { useCartStore } from '../stores/cart'
import { useOrderStore } from '../stores/order'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const cartStore = useCartStore()
const orderStore = useOrderStore()
const authStore = useAuthStore()
const loading = ref(false)

const increaseQty = (index) => {
  cartStore.updateQuantity(index, cartStore.items[index].quantity + 1)
}

const decreaseQty = (index) => {
  cartStore.updateQuantity(index, cartStore.items[index].quantity - 1)
}

const removeItem = (index) => {
  cartStore.removeItem(index)
}

const placeOrder = async () => {
  loading.value = true
  try {
    const orderData = {
      customer_id: authStore.user.id,
      restaurant_id: cartStore.restaurantId,
      order_items: cartStore.items.map(item => ({
        menu_item_id: item.menu_item_id,
        quantity: item.quantity
      }))
    }

    await orderStore.createOrder(orderData)
    cartStore.clearCart()
    alert('Order placed successfully!')
    router.push('/orders')
  } catch (error) {
    alert('Failed to place order: ' + (error.response?.data?.message || error.message))
  } finally {
    loading.value = false
  }
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

.empty-cart {
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.empty-cart p {
  color: #6b7280;
  margin: 1rem 0 2rem;
  font-size: 1.125rem;
}

.cart-content {
  display: grid;
  grid-template-columns: 1fr 380px;
  gap: 2rem;
}

.cart-items {
  background: white;
  border-radius: 1rem;
  padding: 1.5rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.cart-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 0;
  border-bottom: 1px solid #e5e7eb;
}

.cart-item:last-child {
  border-bottom: none;
}

.item-info h4 {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 0.25rem;
}

.price {
  color: #6b7280;
  font-size: 0.875rem;
}

.item-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.quantity-control {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  background: #f3f4f6;
  padding: 0.5rem;
  border-radius: 0.5rem;
}

.quantity-control button {
  width: 32px;
  height: 32px;
  background: white;
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
  font-size: 1.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.quantity-control span {
  min-width: 32px;
  text-align: center;
  font-weight: 600;
}

.remove-btn {
  background: none;
  border: none;
  color: #ef4444;
  cursor: pointer;
  padding: 0.5rem;
}

.cart-summary {
  background: white;
  border-radius: 1rem;
  padding: 1.5rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  height: fit-content;
}

.cart-summary h3 {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 1rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e5e7eb;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem 0;
  color: #6b7280;
}

.summary-row.total {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1f2937;
  padding-top: 1rem;
  margin-top: 1rem;
  border-top: 2px solid #e5e7eb;
}

.btn {
  width: 100%;
  padding: 0.75rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  margin-top: 1.5rem;
}

.btn-primary {
  background: #ea580c;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #c2410c;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .cart-content {
    grid-template-columns: 1fr;
  }
}
</style>