<template>
  <div class="auth-container">
    <div class="auth-card">
      <h2>Create Account</h2>

      <div class="form-group">
        <input
            v-model="form.name"
            type="text"
            placeholder="Full Name"
            class="input"
        />
      </div>

      <div class="form-group">
        <input
            v-model="form.email"
            type="email"
            placeholder="Email"
            class="input"
        />
      </div>

      <div class="form-group">
        <input
            v-model="form.phone_number"
            type="tel"
            placeholder="Phone Number"
            class="input"
        />
      </div>

      <div class="form-group">
        <input
            v-model="form.address"
            type="text"
            placeholder="Address"
            class="input"
        />
      </div>

      <div class="form-group">
        <input
            v-model="form.password"
            type="password"
            placeholder="Password"
            class="input"
        />
      </div>

      <div class="form-group">
        <select v-model="form.role" class="input">
          <option value="customer">Customer</option>
          <option value="shipper">Shipper</option>
          <option value="restaurant_admin">Restaurant Admin</option>
        </select>
      </div>

      <button @click="handleRegister" class="btn btn-primary" :disabled="loading">
        {{ loading ? 'Creating Account...' : 'Register' }}
      </button>

      <div class="auth-footer">
        <router-link to="/login">Already have an account? Login</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  name: '',
  email: '',
  phone_number: '',
  address: '',
  password: '',
  role: 'customer'
})

const loading = ref(false)

const handleRegister = async () => {
  loading.value = true
  try {
    await authStore.register(form.value)
    alert('Registration successful! Please login.')
    router.push('/login')
  } catch (error) {
    alert('Registration failed: ' + (error.response?.data?.message || error.message))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #fb923c 0%, #ef4444 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.auth-card {
  background: white;
  border-radius: 1rem;
  padding: 2rem;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}

.auth-card h2 {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 1.5rem;
  text-align: center;
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

.btn {
  width: 100%;
  padding: 0.75rem;
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

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.auth-footer {
  text-align: center;
  margin-top: 1.5rem;
}

.auth-footer a {
  color: #ea580c;
  text-decoration: none;
}

.auth-footer a:hover {
  text-decoration: underline;
}
</style>