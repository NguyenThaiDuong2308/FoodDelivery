<template>
  <div>
    <Navbar />
    <div class="container">
      <h2 class="page-title">My Profile</h2>

      <div v-if="userStore.loading" class="loading">
        Loading profile...
      </div>

      <div v-else-if="authStore.user" class="profile-card">
        <div class="profile-header">
          <div class="avatar">
            <User :size="48" />
          </div>
          <div class="header-info">
            <h3>{{ authStore.user.name }}</h3>
            <p class="role-badge" :class="authStore.user.role">
              {{ authStore.user.role }}
            </p>
          </div>
        </div>

        <!-- Shipper Status Section (Only for Shippers) -->
        <div v-if="authStore.user.role === 'shipper'" class="shipper-status-section">
          <div class="status-header">
            <h4>
              <Truck :size="20" />
              Shipper Status
            </h4>
            <span class="current-status" :class="currentStatus">
              {{ getStatusLabel(currentStatus) }}
            </span>
          </div>

          <div class="status-buttons">
            <button
                @click="handleStatusChange('available')"
                class="status-btn available"
                :class="{ active: currentStatus === 'available' }"
                :disabled="shipperStore.loading || !shipperId"
            >
              <CheckCircle :size="18" />
              Available
            </button>
            <button
                @click="handleStatusChange('busy')"
                class="status-btn busy"
                :class="{ active: currentStatus === 'busy' }"
                :disabled="shipperStore.loading || !shipperId"
            >
              <Clock :size="18" />
              Busy
            </button>
            <button
                @click="handleStatusChange('offline')"
                class="status-btn offline"
                :class="{ active: currentStatus === 'offline' }"
                :disabled="shipperStore.loading || !shipperId"
            >
              <XCircle :size="18" />
              Offline
            </button>
          </div>
        </div>

        <div class="profile-info">
          <div class="info-row">
            <Mail :size="20" />
            <span>{{ authStore.user.email }}</span>
          </div>
          <div class="info-row">
            <Phone :size="20" />
            <span>{{ authStore.user.phone_number || 'No phone number' }}</span>
          </div>
          <div class="info-row">
            <MapPin :size="20" />
            <span>{{ authStore.user.address || 'No address' }}</span>
          </div>
        </div>

        <button @click="toggleEditForm" class="btn btn-primary">
          <Edit :size="18" />
          {{ showEditForm ? 'Cancel Edit' : 'Edit Profile' }}
        </button>

        <!-- Edit Profile Form -->
        <div v-if="showEditForm" class="edit-form">
          <h4>Edit Profile</h4>
          <div class="form-group">
            <label>Name</label>
            <input v-model="editForm.name" placeholder="Name" class="input" />
          </div>
          <div class="form-group">
            <label>Email</label>
            <input v-model="editForm.email" type="email" placeholder="Email" class="input" />
          </div>
          <div class="form-group">
            <label>Phone Number</label>
            <input v-model="editForm.phone_number" placeholder="Phone" class="input" />
          </div>
          <div class="form-group">
            <label>Address</label>
            <input v-model="editForm.address" placeholder="Address" class="input" />
          </div>
          <div class="form-actions">
            <button @click="handleUpdate" class="btn btn-primary" :disabled="userStore.loading">
              {{ userStore.loading ? 'Saving...' : 'Save Changes' }}
            </button>
            <button @click="showEditForm = false" class="btn btn-secondary">
              Cancel
            </button>
          </div>
        </div>

        <!-- Change Password Section -->
        <div class="password-section">
          <h4>
            <Lock :size="20" />
            Change Password
          </h4>
          <div class="form-group">
            <label>Current Password</label>
            <input
                v-model="passwordForm.old_pass"
                type="password"
                placeholder="Enter current password"
                class="input"
            />
          </div>
          <div class="form-group">
            <label>New Password</label>
            <input
                v-model="passwordForm.new_pass"
                type="password"
                placeholder="Enter new password"
                class="input"
            />
          </div>
          <button @click="handlePasswordChange" class="btn btn-primary" :disabled="userStore.loading">
            {{ userStore.loading ? 'Updating...' : 'Update Password' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  User,
  Mail,
  Phone,
  MapPin,
  Edit,
  Lock,
  Truck,
  CheckCircle,
  Clock,
  XCircle
} from 'lucide-vue-next'
import Navbar from '../components/Navbar.vue'
import { useUserStore } from '../stores/user'
import { useAuthStore } from '../stores/auth'
import { useShipperStore } from '../stores/shipper'

const userStore = useUserStore()
const authStore = useAuthStore()
const shipperStore = useShipperStore()

const showEditForm = ref(false)
const currentStatus = ref('offline')
const shipperId = ref(null) // ← Lưu Shipper ID riêng

const editForm = ref({
  name: '',
  email: '',
  phone_number: '',
  address: ''
})

const passwordForm = ref({
  old_pass: '',
  new_pass: ''
})

onMounted(async () => {
  if (authStore.user) {
    editForm.value = {
      name: authStore.user.name,
      email: authStore.user.email,
      phone_number: authStore.user.phone_number || '',
      address: authStore.user.address || ''
    }

    // Fetch shipper data if user is a shipper
    if (authStore.user.role === 'shipper') {
      try {
        // Fetch all shippers and find the one with matching user_id
        await shipperStore.fetchShippers()
        const currentShipper = shipperStore.shippers.find(
            s => s.user_id === authStore.user.id || s.user?._id === authStore.user.id || s.user?.id === authStore.user.id
        )

        if (currentShipper) {
          shipperId.value = currentShipper.id || currentShipper._id
          currentStatus.value = currentShipper.status || 'offline'
          console.log('Found shipper:', currentShipper)
          console.log('Shipper ID:', shipperId.value)
        } else {
          console.warn('Shipper not found for user:', authStore.user.id)
          showNotification('⚠️ Shipper profile not found', 'warning')
        }
      } catch (error) {
        console.error('Failed to fetch shipper data:', error)
        showNotification('❌ Failed to load shipper data', 'error')
      }
    }
  }
})

const toggleEditForm = () => {
  showEditForm.value = !showEditForm.value
}

const handleUpdate = async () => {
  try {
    const updated = await userStore.updateUser(authStore.user.id, editForm.value)
    console.log("Updated user:", updated)

    // Update auth store with new user data
    const updatedUser = updated.message
    authStore.user.name = updatedUser.name
    authStore.user.email = updatedUser.email
    authStore.user.phone_number = updatedUser.phone_number
    authStore.user.address = updatedUser.address

    showNotification('✅ Profile updated successfully!')
    showEditForm.value = false
  } catch (error) {
    showNotification('❌ Failed to update profile: ' + error.message, 'error')
  }
}

const handlePasswordChange = async () => {
  if (!passwordForm.value.old_pass || !passwordForm.value.new_pass) {
    showNotification('⚠️ Please fill in all password fields', 'warning')
    return
  }

  if (passwordForm.value.new_pass.length < 6) {
    showNotification('⚠️ New password must be at least 6 characters', 'warning')
    return
  }

  try {
    await userStore.resetPassword(passwordForm.value)
    showNotification('✅ Password updated successfully!')
    passwordForm.value = { old_pass: '', new_pass: '' }
  } catch (error) {
    showNotification('❌ Failed to update password: ' + error.message, 'error')
  }
}

const handleStatusChange = async (status) => {
  if (currentStatus.value === status) return

  if (!shipperId.value) {
    showNotification('❌ Shipper ID not found', 'error')
    return
  }

  try {
    console.log('Updating shipper status:', { shipperId: shipperId.value, status })
    await shipperStore.updateStatus(shipperId.value, status) // ← Dùng shipperId thay vì userId
    currentStatus.value = status
    showNotification(`✅ Status changed to ${getStatusLabel(status)}`)
  } catch (error) {
    console.error('Failed to update status:', error)
    showNotification('❌ Failed to update status: ' + error.message, 'error')
  }
}

const getStatusLabel = (status) => {
  const labels = {
    available: '🟢 Available',
    busy: '🟡 Busy',
    offline: '🔴 Offline'
  }
  return labels[status] || status
}

// Simple notification
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
  max-width: 800px;
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

.profile-card {
  background: white;
  border-radius: 1rem;
  padding: 2rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid #e5e7eb;
}

.avatar {
  width: 80px;
  height: 80px;
  background: #fed7aa;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ea580c;
  flex-shrink: 0;
}

.header-info {
  flex: 1;
}

.profile-header h3 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.5rem 0;
}

.role-badge {
  display: inline-block;
  padding: 0.375rem 0.875rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: capitalize;
  margin: 0;
}

.role-badge.customer {
  background: #dbeafe;
  color: #1e40af;
}

.role-badge.shipper {
  background: #dcfce7;
  color: #15803d;
}

.role-badge.restaurant_owner {
  background: #fed7aa;
  color: #ea580c;
}

/* Shipper Status Section */
.shipper-status-section {
  background: #f9fafb;
  border: 2px solid #e5e7eb;
  border-radius: 0.75rem;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.status-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.status-header h4 {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.125rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}

.current-status {
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 700;
}

.current-status.available {
  background: #dcfce7;
  color: #15803d;
}

.current-status.busy {
  background: #fef3c7;
  color: #92400e;
}

.current-status.offline {
  background: #fee2e2;
  color: #991b1b;
}

.status-buttons {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.75rem;
}

.status-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.875rem 1rem;
  border: 2px solid #e5e7eb;
  background: white;
  border-radius: 0.75rem;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.status-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.status-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.status-btn.available {
  color: #15803d;
}

.status-btn.available:hover:not(:disabled) {
  background: #f0fdf4;
  border-color: #86efac;
}

.status-btn.available.active {
  background: #dcfce7;
  border-color: #15803d;
}

.status-btn.busy {
  color: #92400e;
}

.status-btn.busy:hover:not(:disabled) {
  background: #fffbeb;
  border-color: #fde68a;
}

.status-btn.busy.active {
  background: #fef3c7;
  border-color: #92400e;
}

.status-btn.offline {
  color: #991b1b;
}

.status-btn.offline:hover:not(:disabled) {
  background: #fef2f2;
  border-color: #fecaca;
}

.status-btn.offline.active {
  background: #fee2e2;
  border-color: #991b1b;
}

.profile-info {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 2rem;
}

.info-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  color: #6b7280;
  font-size: 0.95rem;
}

.info-row svg {
  flex-shrink: 0;
  color: #9ca3af;
}

.edit-form,
.password-section {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #e5e7eb;
}

.edit-form h4,
.password-section h4 {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.25rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 1.25rem 0;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 600;
  color: #374151;
  margin-bottom: 0.5rem;
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

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #ea580c;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #c2410c;
  transform: translateY(-1px);
}

.btn-secondary {
  background: #f3f4f6;
  color: #1f2937;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

:global(.notification.warning) {
  background: #f59e0b;
  color: white;
}

/* Responsive */
@media (max-width: 768px) {
  .container {
    padding: 1rem 0.5rem;
  }

  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .status-buttons {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column;
  }

  .form-actions .btn {
    width: 100%;
  }
}
</style>