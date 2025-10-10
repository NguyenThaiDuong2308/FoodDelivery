<template>
  <div>
    <Navbar />

    <div class="container">
      <div class="header">
        <h2 class="page-title">Restaurant Management</h2>
        <button
            class="btn btn-primary"
            @click="openForm()"
            :disabled="!canCreateRestaurant"
        >
          <Plus :size="20" />
          Add Restaurant
        </button>
      </div>

      <!-- Warning message when user already has a restaurant -->
      <div v-if="!canCreateRestaurant" class="warning-banner">
        <AlertCircle :size="20" />
        <span>You already have a restaurant. Each user can only create one restaurant.</span>
      </div>

      <div v-if="restaurantStore.loading" class="loading">
        Loading restaurants...
      </div>

      <div v-else class="restaurant-grid">
        <div v-for="r in restaurantStore.restaurants" :key="r.id" class="restaurant-card">
          <div class="restaurant-header">
            <div class="avatar">
              <Store :size="32" />
            </div>
            <div class="restaurant-title">
              <h3>{{ r.name }}</h3>
              <span class="status-badge" :class="r.status">{{ r.status }}</span>
            </div>
          </div>

          <div class="restaurant-info">
            <div class="info-row">
              <MapPin :size="18" />
              <span>{{ r.address || 'No address' }}</span>
            </div>
            <div class="info-row">
              <Phone :size="18" />
              <span>{{ r.phone_number || 'No phone' }}</span>
            </div>
            <div class="info-row">
              <Mail :size="18" />
              <span>{{ r.email || 'No email' }}</span>
            </div>
            <div v-if="r.description" class="info-row description">
              <FileText :size="18" />
              <span>{{ r.description }}</span>
            </div>
          </div>

          <div class="card-actions">
            <button @click="openMenu(r)" class="btn btn-primary">
              <Menu :size="18" />
              Manage Menu
            </button>
            <button @click="openForm(r)" class="btn btn-secondary">
              <Edit :size="18" />
              Edit
            </button>
          </div>
        </div>
      </div>

      <div v-if="!restaurantStore.loading && restaurantStore.restaurants.length === 0" class="empty-state">
        <Store :size="64" />
        <h3>No restaurants yet</h3>
        <p>Get started by adding your first restaurant</p>
        <button class="btn btn-primary" @click="openForm()">
          <Plus :size="20" />
          Add Restaurant
        </button>
      </div>
    </div>

    <!-- Restaurant Form Modal -->
    <div v-if="showForm" class="modal" @click.self="closeForm">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ formMode === 'edit' ? 'Edit Restaurant' : 'Add New Restaurant' }}</h3>
          <button @click="closeForm" class="close-btn">
            <X :size="24" />
          </button>
        </div>

        <form @submit.prevent="submitForm">
          <div class="form-group">
            <label>Restaurant Name *</label>
            <input v-model="form.name" placeholder="Enter restaurant name" class="input" required />
          </div>

          <div class="form-group">
            <label>Description</label>
            <textarea v-model="form.description" placeholder="Enter description" class="input textarea" rows="3"></textarea>
          </div>

          <div class="form-group">
            <label>Address</label>
            <input v-model="form.address" placeholder="Enter address" class="input" />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Phone Number</label>
              <input v-model="form.phone_number" placeholder="Enter phone" class="input" />
            </div>

            <div class="form-group">
              <label>Email</label>
              <input v-model="form.email" type="email" placeholder="Enter email" class="input" />
            </div>
          </div>

          <div class="form-group">
            <label>Status</label>
            <select v-model="form.status" class="input">
              <option value="open">Open</option>
              <option value="closed">Closed</option>
            </select>
          </div>

          <div class="modal-actions">
            <button type="submit" class="btn btn-primary" :disabled="restaurantStore.loading">
              {{ restaurantStore.loading ? 'Saving...' : 'Save Restaurant' }}
            </button>
            <button type="button" @click="closeForm" class="btn btn-secondary">
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Menu Modal -->
    <div v-if="showMenu" class="modal" @click.self="closeMenu">
      <div class="modal-content modal-large">
        <div class="modal-header">
          <div>
            <h3>Menu Items</h3>
            <p class="modal-subtitle">{{ selectedRestaurant.name }}</p>
          </div>
          <button @click="closeMenu" class="close-btn">
            <X :size="24" />
          </button>
        </div>

        <div class="menu-section">
          <button @click="openMenuItemForm()" class="btn btn-primary">
            <Plus :size="20" />
            Add Menu Item
          </button>

          <div v-if="selectedRestaurant.menu_items && selectedRestaurant.menu_items.length > 0" class="menu-list">
            <div v-for="item in selectedRestaurant.menu_items" :key="item.id" class="menu-item">
              <div class="menu-item-icon">
                <UtensilsCrossed :size="24" />
              </div>
              <div class="menu-item-info">
                <h4>{{ item.name }}</h4>
                <p v-if="item.description" class="description">{{ item.description }}</p>
                <div class="menu-item-details">
                  <span class="price">${{ item.price }}</span>
                  <span class="availability-badge" :class="{ available: item.available, unavailable: !item.available }">
                    {{ item.available ? 'Available' : 'Unavailable' }}
                  </span>
                </div>
              </div>
              <div class="menu-item-actions">
                <button @click="openMenuItemForm(item)" class="btn-icon" title="Edit">
                  <Edit :size="18" />
                </button>
                <button @click="deleteMenuItem(item.id)" class="btn-icon btn-icon-danger" title="Delete">
                  <Trash2 :size="18" />
                </button>
              </div>
            </div>
          </div>

          <div v-else class="empty-menu">
            <UtensilsCrossed :size="48" />
            <p>No menu items yet</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Menu Item Form Modal -->
    <div v-if="showMenuItemForm" class="modal" @click.self="closeMenuItemForm">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ menuItemFormMode === 'edit' ? 'Edit Menu Item' : 'Add Menu Item' }}</h3>
          <button @click="closeMenuItemForm" class="close-btn">
            <X :size="24" />
          </button>
        </div>

        <form @submit.prevent="submitMenuItemForm">
          <div class="form-group">
            <label>Item Name *</label>
            <input
                v-model="menuItemForm.name"
                placeholder="e.g., Margherita Pizza"
                class="input"
                required
            />
          </div>

          <div class="form-group">
            <label>Description</label>
            <textarea
                v-model="menuItemForm.description"
                placeholder="Describe your dish..."
                class="input textarea"
                rows="3"
            ></textarea>
          </div>

          <div class="form-group">
            <label>Price *</label>
            <input
                v-model.number="menuItemForm.price"
                type="number"
                step="0.01"
                min="0"
                placeholder="0.00"
                class="input"
                required
            />
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input
                  v-model="menuItemForm.available"
                  type="checkbox"
                  class="checkbox"
              />
              <span>Available for order</span>
            </label>
          </div>

          <div class="modal-actions">
            <button type="submit" class="btn btn-primary" :disabled="restaurantStore.loading">
              {{ restaurantStore.loading ? 'Saving...' : 'Save Menu Item' }}
            </button>
            <button type="button" @click="closeMenuItemForm" class="btn btn-secondary">
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, computed } from 'vue'
import { Store, MapPin, Phone, Mail, FileText, Menu, Edit, Trash2, Plus, X, UtensilsCrossed, AlertCircle } from 'lucide-vue-next'
import { useRestaurantStore } from '../stores/restaurant'
import Navbar from '../components/Navbar.vue'

const restaurantStore = useRestaurantStore()

const showForm = ref(false)
const showMenu = ref(false)
const showMenuItemForm = ref(false)
const formMode = ref('add') // 'add' | 'edit'
const menuItemFormMode = ref('add') // 'add' | 'edit'

const form = reactive({
  id: null,
  name: '',
  description: '',
  address: '',
  phone_number: '',
  email: '',
  status: 'open'
})

const selectedRestaurant = reactive({
  id: null,
  name: '',
  menu_items: []
})

const menuItemForm = reactive({
  id: null,
  name: '',
  description: '',
  price: 0,
  available: true
})

// Check if user can create a restaurant (limit: 1 restaurant per user)
const canCreateRestaurant = computed(() => {
  return restaurantStore.restaurants.length === 0
})

// ---------- Life cycle ----------
onMounted(() => {
  restaurantStore.fetchRestaurants()
})

// ---------- Restaurant Actions ----------
const openForm = (restaurant = null) => {
  if (restaurant) {
    formMode.value = 'edit'
    Object.assign(form, restaurant)
  } else {
    // Check if user already has a restaurant
    if (!canCreateRestaurant.value) {
      alert('You already have a restaurant. Each user can only create one restaurant.')
      return
    }
    formMode.value = 'add'
    Object.assign(form, {
      id: null,
      name: '',
      description: '',
      address: '',
      phone_number: '',
      email: '',
      status: 'open'
    })
  }
  showForm.value = true
}

const closeForm = () => (showForm.value = false)

const submitForm = async () => {
  try {
    if (formMode.value === 'add') {
      // Double check before creating
      if (!canCreateRestaurant.value) {
        alert('You already have a restaurant. Each user can only create one restaurant.')
        return
      }
      await restaurantStore.createRestaurant(form)
    } else {
      await restaurantStore.updateRestaurant(form.id, form)
    }
    closeForm()
  } catch (error) {
    alert('Failed to save restaurant: ' + error.message)
  }
}


// ---------- Menu Actions ----------
const openMenu = async (restaurant) => {
  try {
    await restaurantStore.fetchMenu(restaurant.id)
    Object.assign(selectedRestaurant, restaurant)
    showMenu.value = true
  } catch (error) {
    alert('Failed to load menu: ' + error.message)
  }
}

const closeMenu = () => (showMenu.value = false)

// ---------- Menu Item Form Actions ----------
const openMenuItemForm = (item = null) => {
  if (item) {
    menuItemFormMode.value = 'edit'
    Object.assign(menuItemForm, {
      id: item.id,
      name: item.name,
      description: item.description || '',
      price: item.price,
      available: item.available !== undefined ? item.available : true
    })
  } else {
    menuItemFormMode.value = 'add'
    Object.assign(menuItemForm, {
      id: null,
      name: '',
      description: '',
      price: 0,
      available: true
    })
  }
  showMenuItemForm.value = true
}

const closeMenuItemForm = () => (showMenuItemForm.value = false)

const submitMenuItemForm = async () => {
  try {
    const payload = {
      name: menuItemForm.name,
      description: menuItemForm.description,
      price: menuItemForm.price,
      available: menuItemForm.available
    }

    if (menuItemFormMode.value === 'add') {
      await restaurantStore.createMenuItem(selectedRestaurant.id, payload)
    } else {
      await restaurantStore.updateMenuItem(selectedRestaurant.id, menuItemForm.id, payload)
    }

    await restaurantStore.fetchMenu(selectedRestaurant.id)
    closeMenuItemForm()
  } catch (error) {
    alert('Failed to save menu item: ' + error.message)
  }
}

const deleteMenuItem = async (menuItemId) => {
  if (confirm('Delete this menu item?')) {
    try {
      await restaurantStore.deleteMenuItem(selectedRestaurant.id, menuItemId)
      await restaurantStore.fetchMenu(selectedRestaurant.id)
    } catch (error) {
      alert('Failed to delete menu item: ' + error.message)
    }
  }
}
</script>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}

.warning-banner {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  background: #fef3c7;
  border: 1px solid #fcd34d;
  border-radius: 0.75rem;
  padding: 1rem 1.25rem;
  margin-bottom: 1.5rem;
  color: #92400e;
  font-size: 0.9rem;
  font-weight: 500;
}

.warning-banner svg {
  flex-shrink: 0;
}

.loading {
  text-align: center;
  padding: 3rem;
  color: #6b7280;
  font-size: 1.125rem;
}

.restaurant-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.restaurant-card {
  background: white;
  border-radius: 1rem;
  padding: 1.5rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  transition: transform 0.2s, box-shadow 0.2s;
}

.restaurant-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0,0,0,0.15);
}

.restaurant-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e5e7eb;
}

.avatar {
  width: 60px;
  height: 60px;
  background: #fed7aa;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ea580c;
  flex-shrink: 0;
}

.restaurant-title {
  flex: 1;
}

.restaurant-title h3 {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.5rem 0;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
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

.restaurant-info {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.info-row {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  color: #6b7280;
  font-size: 0.9rem;
}

.info-row svg {
  flex-shrink: 0;
  margin-top: 2px;
}

.info-row.description {
  align-items: flex-start;
}

.card-actions {
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 0.5rem;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn-primary {
  background: #ea580c;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #c2410c;
}

.btn-primary:disabled {
  background: #d1d5db;
  cursor: not-allowed;
}

.btn-secondary {
  background: #f3f4f6;
  color: #1f2937;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.btn-danger {
  background: #fee2e2;
  color: #991b1b;
}

.btn-danger:hover {
  background: #fecaca;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

.empty-state h3 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.5rem 0;
}

.empty-state p {
  color: #6b7280;
  margin: 0 0 1.5rem 0;
}

/* Modal Styles */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-content {
  background: white;
  border-radius: 1rem;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.modal-large {
  max-width: 700px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}

.modal-subtitle {
  color: #6b7280;
  font-size: 0.9rem;
  margin: 0.25rem 0 0 0;
}

.close-btn {
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 0.375rem;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.modal-content form {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 600;
  color: #374151;
  margin-bottom: 0.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: all 0.2s;
}

.input:focus {
  outline: none;
  border-color: #ea580c;
  box-shadow: 0 0 0 3px rgba(234, 88, 12, 0.1);
}

.textarea {
  resize: vertical;
  min-height: 80px;
  font-family: inherit;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  font-weight: 500;
  color: #1f2937;
}

.checkbox {
  width: 20px;
  height: 20px;
  cursor: pointer;
  accent-color: #ea580c;
}

.modal-actions {
  display: flex;
  gap: 0.75rem;
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e5e7eb;
}

.modal-actions .btn {
  flex: 1;
}

/* Menu Section */
.menu-section {
  padding: 1.5rem;
}

.menu-list {
  margin-top: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.menu-item {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1rem;
  background: #f9fafb;
  border-radius: 0.75rem;
  transition: all 0.2s;
}

.menu-item:hover {
  background: #f3f4f6;
}

.menu-item-icon {
  width: 48px;
  height: 48px;
  background: #fed7aa;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ea580c;
  flex-shrink: 0;
}

.menu-item-info {
  flex: 1;
}

.menu-item-info h4 {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 0.25rem 0;
}

.menu-item-info .description {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0.25rem 0 0.5rem 0;
  line-height: 1.4;
}

.menu-item-details {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-top: 0.5rem;
}

.menu-item-info .price {
  font-size: 1.125rem;
  font-weight: 700;
  color: #ea580c;
}

.availability-badge {
  display: inline-block;
  padding: 0.25rem 0.625rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
}

.availability-badge.available {
  background: #dcfce7;
  color: #15803d;
}

.availability-badge.unavailable {
  background: #fee2e2;
  color: #991b1b;
}

.menu-item-actions {
  display: flex;
  gap: 0.5rem;
  flex-shrink: 0;
}

.btn-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 0.5rem;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-icon:hover {
  background: #f3f4f6;
  border-color: #9ca3af;
  color: #1f2937;
}

.btn-icon-danger:hover {
  background: #fee2e2;
  border-color: #fecaca;
  color: #991b1b;
}

.empty-menu {
  text-align: center;
  padding: 3rem 1rem;
  color: #9ca3af;
}

.empty-menu svg {
  margin-bottom: 1rem;
}

.empty-menu p {
  margin: 0;
  font-size: 0.9rem;
}

@media (max-width: 768px) {
  .restaurant-grid {
    grid-template-columns: 1fr;
  }

  .header {
    flex-direction: column;
    align-items: stretch;
    gap: 1rem;
  }

  .card-actions {
    grid-template-columns: 1fr;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .menu-item {
    flex-direction: column;
  }

  .menu-item-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>