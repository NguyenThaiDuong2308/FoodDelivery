<template>
  <div>
    <Navbar />
    <div class="container">
      <h2 class="page-title">User Management</h2>

      <div v-if="userStore.loading" class="loading">
        Loading users...
      </div>

      <div v-else-if="userStore.users.length === 0" class="empty-state">
        <Users :size="64" />
        <p>No users found</p>
      </div>

      <div v-else class="users-table">
        <table>
          <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Email</th>
            <th>Phone</th>
            <th>Role</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="user in userStore.users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.name }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.phone_number || 'N/A' }}</td>
            <td>
                <span class="role-badge" :class="user.role">
                  {{ user.role }}
                </span>
            </td>
            <td>

            </td>
          </tr>
          </tbody>
        </table>
      </div>

      <!-- User Detail Modal -->
      <div v-if="showModal" class="modal" @click.self="closeModal">
        <div class="modal-content">
          <div class="modal-header">
            <h3>User Details</h3>
            <button @click="closeModal" class="close-btn">
              <X :size="24" />
            </button>
          </div>
          <div class="modal-body">
            <div class="detail-row">
              <strong>ID:</strong>
              <span>{{ selectedUser.id }}</span>
            </div>
            <div class="detail-row">
              <strong>Name:</strong>
              <span>{{ selectedUser.name }}</span>
            </div>
            <div class="detail-row">
              <strong>Email:</strong>
              <span>{{ selectedUser.email }}</span>
            </div>
            <div class="detail-row">
              <strong>Phone:</strong>
              <span>{{ selectedUser.phone_number || 'N/A' }}</span>
            </div>
            <div class="detail-row">
              <strong>Address:</strong>
              <span>{{ selectedUser.address || 'N/A' }}</span>
            </div>
            <div class="detail-row">
              <strong>Role:</strong>
              <span class="role-badge" :class="selectedUser.role">
                {{ selectedUser.role }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Users, Eye, Trash2, X } from 'lucide-vue-next'
import Navbar from '../components/Navbar.vue'
import { useUserStore } from '../stores/user'
import { useAuthStore } from '../stores/auth'

const userStore = useUserStore()
const authStore = useAuthStore()

const showModal = ref(false)
const selectedUser = ref({})

onMounted(() => {
  userStore.fetchUsers()
})

const viewUser = (user) => {
  selectedUser.value = user
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const showNotification = (message, type = 'success') => {
  const notification = document.createElement('div')
  notification.className = `notification ${type}`
  notification.textContent = message
  document.body.appendChild(notification)

  setTimeout(() => notification.classList.add('show'), 10)

  setTimeout(() => {
    notification.classList.remove('show')
    setTimeout(() => document.body.removeChild(notification), 300)
  }, 3000)
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
  padding: 4rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.empty-state svg {
  color: #d1d5db;
  margin-bottom: 1rem;
}

.users-table {
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  overflow: hidden;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  background: #f9fafb;
}

th {
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid #e5e7eb;
}

td {
  padding: 1rem;
  border-bottom: 1px solid #e5e7eb;
}

tr:hover {
  background: #f9fafb;
}

.role-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: capitalize;
}

/* Modal */
.modal {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 1rem;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25rem;
  color: #6b7280;
}

.close-btn:hover {
  color: #1f2937;
}

.modal-body {
  padding: 1.5rem;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem 0;
  border-bottom: 1px solid #f3f4f6;
}

.detail-row:last-child {
  border-bottom: none;
}

/* Notification */

@media (max-width: 768px) {
  .users-table {
    overflow-x: auto;
  }

  table {
    min-width: 600px;
  }
}
</style>