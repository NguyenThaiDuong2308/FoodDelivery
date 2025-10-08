<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <div class="icon-wrapper">
          <Lock :size="32" color="#ea580c" />
        </div>
        <h1>Forgot Password</h1>
        <p>Enter your email to receive reset instructions</p>
      </div>

      <div v-if="!emailSent">
        <div class="form-group">
          <input
              v-model="email"
              type="email"
              placeholder="Email"
              class="input"
              @keyup.enter="handleForgotPassword"
          />
        </div>

        <button @click="handleForgotPassword" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Sending...' : 'Send Reset Link' }}
        </button>
      </div>

      <div v-else class="success-message">
        <CheckCircle :size="48" color="#15803d" />
        <p class="success-title">Email Sent!</p>
        <p class="success-text">Password reset link has been sent to your email!</p>
        <p class="sub-message">Please check your inbox and follow the instructions.</p>
      </div>

      <div class="auth-footer">
        <router-link to="/login">← Back to Login</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Lock, CheckCircle } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const email = ref('')
const loading = ref(false)
const emailSent = ref(false)

const handleForgotPassword = async () => {
  if (!email.value) {
    alert('Please enter your email')
    return
  }

  loading.value = true
  try {
    await authStore.forgotPassword(email.value)
    emailSent.value = true
  } catch (error) {
    alert('Failed to send reset link: ' + (error.response?.data?.message || error.message))
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
  margin-bottom: 0.5rem;
}

.sub-message {
  color: #6b7280;
  font-size: 0.875rem;
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
</style>