<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <div class="icon-wrapper">
          <Key :size="32" color="#ea580c" />
        </div>
        <h1>Reset Password</h1>
        <p>Enter your new password</p>
      </div>

      <div v-if="!resetSuccess">
        <div class="form-group">
          <div class="input-wrapper">
            <input
                v-model="form.new_password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="New Password (8-64 characters)"
                class="input"
            />
            <button
                type="button"
                @click="showPassword = !showPassword"
                class="toggle-password"
            >
              <Eye v-if="!showPassword" :size="20" />
              <EyeOff v-else :size="20" />
            </button>
          </div>
          <p v-if="errors.password" class="error-text">{{ errors.password }}</p>
        </div>

        <div class="form-group">
          <div class="input-wrapper">
            <input
                v-model="confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                placeholder="Confirm New Password"
                class="input"
            />
            <button
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
                class="toggle-password"
            >
              <Eye v-if="!showConfirmPassword" :size="20" />
              <EyeOff v-else :size="20" />
            </button>
          </div>
          <p v-if="errors.password" class="error-text">{{ errors.password }}</p>
        </div>


        <button
            @click="handleResetPassword"
            class="btn btn-primary"
            :disabled="loading || !isValid"
        >
          {{ loading ? 'Resetting...' : 'Reset Password' }}
        </button>
      </div>

      <div v-else class="success-message">
        <CheckCircle :size="48" color="#15803d" />
        <p class="success-title">Password Reset Successful!</p>
        <p class="success-text">Your password has been changed successfully.</p>
        <router-link to="/login" class="btn btn-primary">
          Go to Login
        </router-link>
      </div>

      <div class="auth-footer">
        <router-link to="/login">← Back to Login</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, computed} from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Key, CheckCircle, Eye, EyeOff } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  reset_token: route.query.token || '',
  new_password: ''
})

const errors = ref({
  password: ''
})

const confirmPassword = ref('')
const loading = ref(false)
const resetSuccess = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)


const isValid = computed(() => {
  return form.value.new_password &&
      form.value.new_password === confirmPassword.value &&
      form.value.new_password.length >= 8 &&
      form.value.new_password.length <= 64
})

const handleResetPassword = async () => {
  errors.value.password = ''
  if (!form.value.reset_token) {
    alert('Invalid or missing reset token')
    router.push('/forgot-password')
    return
  }

  if (form.value.new_password !== confirmPassword.value) {
    errors.value.password = 'Passwords do not match'
    return
  }

  if (form.value.new_password.length < 8 || form.value.new_password.length > 64) {
    errors.value.password = 'Password must be between 8 and 64 characters'
    return
  }

  loading.value = true
  try {
    await authStore.resetForgotPassword(form.value)
    resetSuccess.value = true
  } catch (error) {
    alert('Failed to reset password: ' + (error.response?.data?.message || error.message))
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
  font-size: 1.75rem;
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

.input-wrapper {
  position: relative;
}

.input {
  width: 100%;
  padding: 0.75rem 2.5rem 0.75rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.5rem;
  font-size: 1rem;
}

.input:focus {
  outline: none;
  border-color: #ea580c;
  box-shadow: 0 0 0 3px rgba(234, 88, 12, 0.1);
}

.toggle-password {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 0.25rem;
  display: flex;
  align-items: center;
}

.toggle-password:hover {
  color: #ea580c;
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
  text-decoration: none;
  display: block;
  text-align: center;
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

.success-message {
  text-align: center;
  padding: 2rem 0;
}

.success-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #15803d;
  margin: 1rem 0 0.5rem;
}

.success-text {
  color: #1f2937;
  font-weight: 600;
  margin-bottom: 1.5rem;
}

.auth-footer {
  text-align: center;
  margin-top: 1.5rem;
}

.auth-footer a {
  color: #ea580c;
  text-decoration: none;
  font-weight: 600;
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