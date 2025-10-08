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
          <div>
            <h3>{{ authStore.user.name }}</h3>
            <p class="role-badge">{{ authStore.user.role }}</p>
          </div>
        </div>

        <div class="profile-info">
          <div class="info-row">
            <Mail :size="20" />
            <span>{{ authStore.user.email }}</span>
          </div>
          <div class="info-row">
            <Phone :size="20" />
            <span>{{ authStore.user.phone_number }}</span>
          </div>
          <div class="info-row">
            <MapPin :size="20" />
            <span>{{ authStore.user.address }}</span>
          </div>
        </div>

        <button @click="showEditForm = true" class="btn btn-primary">
          Edit Profile
        </button>

        <div v-if="showEditForm" class="edit-form">
          <h4>Edit Profile</h4>
          <div class="form-group">
            <input v-model="editForm.name" placeholder="Name" class="input" />
          </div>
          <div class="form-group">
            <input v-model="editForm.email" type="email" placeholder="Email" class="input" />
          </div>
          <div class="form-group">
            <input v-model="editForm.phone_number" placeholder="Phone" class="input" />
          </div>
          <div class="form-group">
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

        <div class="password-section">
          <h4>Change Password</h4>
          <div class="form-group">
            <input v-model="passwordForm.old_pass" type="password" placeholder="Current Password" class="input" />
          </div>
          <div class="form-group">
            <input v-model="passwordForm.new_pass" type="password" placeholder="New Password" class="input" />
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
import { User, Mail, Phone, MapPin } from 'lucide-vue-next'
import Navbar from '../components/Navbar.vue'
import { useUserStore } from '../stores/user'
import { useAuthStore } from '../stores/auth'

const userStore = useUserStore()
const authStore = useAuthStore()

const showEditForm = ref(false)
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

onMounted(() => {
  if (authStore.user) {
    editForm.value = {
      name: authStore.user.name,
      email: authStore.user.email,
      phone_number: authStore.user.phone_number,
      address: authStore.user.address
    }
  }
})

const handleUpdate = async () => {
  try {
    const updated = await userStore.updateUser(authStore.user.id, editForm.value)
    console.log("Updated user:",updated)
    // Update auth store with new user data
    const updatedUser = updated.message
    authStore.user.name = updatedUser.name
    authStore.user.email = updatedUser.email
    authStore.user.phone_number = updatedUser.phone_number
    authStore.user.address = updatedUser.address
    alert('Profile updated successfully!')
    showEditForm.value = false
  } catch (error) {
    alert('Failed to update profile: ' + error.message)
  }
}

const handlePasswordChange = async () => {
  try {
    await userStore.resetPassword(passwordForm.value)
    alert('Password updated successfully!')
    passwordForm.value = { old_pass: '', new_pass: '' }
  } catch (error) {
    alert('Failed to update password: ' + error.message)
  }
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
}

.profile-header h3 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 0.25rem;
}

.role-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  background: #dbeafe;
  color: #1e40af;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: capitalize;
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
}

.edit-form,
.password-section {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #e5e7eb;
}

.edit-form h4,
.password-section h4 {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.5rem;
  font-size: 1rem;
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
</style>