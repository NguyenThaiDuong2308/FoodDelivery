<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <div class="icon-wrapper">
          <Truck :size="32" color="#ea580c" />
        </div>
        <h1>Delivery App</h1>
        <p>Sign in to continue</p>
      </div>

      <div class="form-group">
        <input
            v-model="form.email"
            type="email"
            placeholder="Email"
            class="input"
            @blur="validateField('email')"
        />
        <p v-if="errors.email" class="error-text">{{ errors.email }}</p>
      </div>

      <div class="form-group">
        <input
            v-model="form.password"
            type="password"
            placeholder="Password"
            class="input"
            @blur="validateField('password')"
        />
        <p v-if="errors.password" class="error-text">{{ errors.password }}</p>
      </div>

      <div class="forgot-link">
        <router-link to="/forgot-password">Forgot your password?</router-link>
      </div>
      <button @click="handleLogin" class="btn btn-primary" :disabled="loading">
        {{ loading ? 'Signing in...' : 'Sign In' }}
      </button>

      <div class="auth-footer">
        <router-link to="/register">Don't have an account? Register</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Truck } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  email: '',
  password: ''
})

const errors = ref({
      email: '',
      password: '',
    }
)
const loading = ref(false)

const validateField = (field) => {
  if (field === 'email') {
    if (!form.value.email) {
      errors.value.email = 'Email is required'
    } else if (!/^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$/.test(form.value.email)) {
      errors.value.email = 'Invalid email format'
    } else {
      errors.value.email = ''
    }
  }

  if (field === 'password') {
    if (!form.value.password) {
      errors.value.password = 'Password is required'
    } else if (form.value.password.length < 8 || form.value.password.length>64) {
      errors.value.password = 'Password must be at least 8 and no more 64 characters'
    } else {
      errors.value.password = ''
    }
  }
}

const validateForm = () => {
  if (!form.value.email || !form.value.password) {
    return false
  }
  return true
}

const handleLogin = async () => {
  if (!validateForm()) return
  loading.value = true
  try {
    await authStore.login(form.value)
    await router.push('/')
  } catch (error) {
    alert('Login failed: Invalid email or password')
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

.auth-header {
  text-align: center;
  margin-bottom: 2rem;
}

.icon-wrapper {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  background: #fed7aa;
  border-radius: 50%;
  margin-bottom: 1rem;
}

.auth-header h1 {
  font-size: 2rem;
  font-weight: 700;
  color: #1f2937;
}

.auth-header p {
  color: #6b7280;
  margin-top: 0.5rem;
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

.forgot-link {
  text-align: right;
  margin-bottom: 1rem;
}

.forgot-link a {
  color: #6b7280;
  text-decoration: none;
  font-size: 0.875rem;
}

.forgot-link a:hover {
  color: #ea580c;
  text-decoration: underline;
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

.error-text {
  color: red;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}
</style>