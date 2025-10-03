<template>
  <div>
    <Navbar />
    <div class="container">
      <h2 class="page-title">Live Shipper Tracking</h2>

      <div class="map-container">
        <div class="map-view">
          <div class="map-background"></div>
          <div class="map-overlay">
            <MapPin :size="64" color="#ea580c" class="map-icon" />
            <p class="map-text">Map View - Shipper Locations</p>
            <p class="map-subtext">Real-time updates every 5 seconds</p>
          </div>

          <div
              v-for="(shipper, index) in shipperStore.shippers"
              :key="shipper.id"
              class="shipper-marker"
              :style="{
              left: `${20 + (index * 15) % 60}%`,
              top: `${30 + (index * 20) % 40}%`
            }"
          >
            <Truck :size="16" />
            {{ shipper.name }}
          </div>
        </div>
      </div>

      <div class="shippers-grid">
        <div
            v-for="shipper in shipperStore.shippers"
            :key="shipper.id"
            class="shipper-card"
        >
          <div class="shipper-header">
            <div>
              <h3>{{ shipper.name }}</h3>
              <p class="phone">{{ shipper.phone_number }}</p>
            </div>
            <span class="status-badge" :class="shipper.status">
              {{ shipper.status }}
            </span>
          </div>

          <div v-if="shipper.current_location" class="location-info">
            <p class="location-title">
              <MapPin :size="16" />
              Current Location
            </p>
            <p class="coordinates">Lat: {{ shipper.current_location.latitude?.toFixed(6) }}</p>
            <p class="coordinates">Lng: {{ shipper.current_location.longitude?.toFixed(6) }}</p>
            <p class="update-info">Updates every 5 seconds</p>
          </div>

          <div v-else class="location-info">
            <p class="no-location">Location not available</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { MapPin, Truck } from 'lucide-vue-next'
import Navbar from '../components/Navbar.vue'
import { useShipperStore } from '../stores/shipper'
import { useAuthStore } from '../stores/auth'
import { shipperService } from '../services/shipper.service'

const shipperStore = useShipperStore()
const authStore = useAuthStore()
const updateInterval = ref(null)

onMounted(() => {
  shipperStore.fetchShippers()

  // Update shippers every 5 seconds
  updateInterval.value = setInterval(() => {
    shipperStore.fetchShippers()
  }, 5000)

  // If user is a shipper, update their location every 5 seconds
  if (authStore.user?.role === 'shipper') {
    const locationInterval = setInterval(() => {
      const lat = 21.0285 + (Math.random() - 0.5) * 0.1
      const lng = 105.8542 + (Math.random() - 0.5) * 0.1
      shipperService.updateLocation(authStore.user.id, {
        latitude: lat,
        longitude: lng
      })
    }, 5000)

    onUnmounted(() => {
      clearInterval(locationInterval)
    })
  }
})

onUnmounted(() => {
  if (updateInterval.value) {
    clearInterval(updateInterval.value)
  }
})
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

.map-container {
  background: white;
  border-radius: 1rem;
  padding: 1.5rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  margin-bottom: 2rem;
}

.map-view {
  height: 400px;
  border-radius: 0.75rem;
  position: relative;
  overflow: hidden;
}

.map-background {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #dbeafe 0%, #dcfce7 100%);
}

.map-overlay {
  position: relative;
  z-index: 10;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.map-icon {
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.map-text {
  color: #6b7280;
  font-size: 1.125rem;
  margin-top: 1rem;
}

.map-subtext {
  color: #9ca3af;
  font-size: 0.875rem;
  margin-top: 0.5rem;
}

.shipper-marker {
  position: absolute;
  background: #ea580c;
  color: white;
  padding: 0.5rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(234, 88, 12, 0.4);
  display: flex;
  align-items: center;
  gap: 0.25rem;
  transition: all 1s ease-in-out;
  z-index: 20;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 4px 12px rgba(234, 88, 12, 0.4);
  }
  50% {
    box-shadow: 0 4px 20px rgba(234, 88, 12, 0.6);
  }
}

.shippers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.shipper-card {
  background: white;
  border-radius: 1rem;
  padding: 1.5rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  transition: all 0.3s;
}

.shipper-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.15);
}

.shipper-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.shipper-header h3 {
  font-size: 1.125rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 0.25rem;
}

.phone {
  color: #6b7280;
  font-size: 0.875rem;
}

.status-badge {
  padding: 0.375rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: capitalize;
}

.status-badge.available {
  background: #dcfce7;
  color: #15803d;
}

.status-badge.busy {
  background: #fee2e2;
  color: #991b1b;
}

.status-badge.unavailable {
  background: #f3f4f6;
  color: #6b7280;
}

.location-info {
  background: #f9fafb;
  border-radius: 0.5rem;
  padding: 1rem;
}

.location-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #1f2937;
  font-weight: 600;
  margin-bottom: 0.75rem;
}

.coordinates {
  color: #6b7280;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
}

.update-info {
  color: #9ca3af;
  font-size: 0.75rem;
  margin-top: 0.5rem;
}

.no-location {
  color: #9ca3af;
  font-size: 0.875rem;
  text-align: center;
}
</style>