<template>
  <div>
    <Navbar />
    <div class="container">
      <div class="header">
        <h2 class="page-title">
          <Truck :size="32" />
          Live Shipper Tracking
        </h2>
        <div class="header-info">
          <span class="info-badge">
            <Activity :size="16" />
            {{ activeShippersCount }} Active
          </span>
          <span class="info-badge">
            <Clock :size="16" />
            Updates every 5s
          </span>
        </div>
      </div>

      <!-- Map Container -->
      <div class="map-container">
        <div id="map" class="map-view"></div>

        <!-- Map Legend -->
        <div class="map-legend">
          <h4>Shipper Status</h4>
          <div class="legend-items">
            <div class="legend-item">
              <span class="legend-dot available"></span>
              <span>Available</span>
            </div>
            <div class="legend-item">
              <span class="legend-dot busy"></span>
              <span>Busy</span>
            </div>
            <div class="legend-item">
              <span class="legend-dot offline"></span>
              <span>Offline</span>
            </div>
          </div>
        </div>

        <!-- Loading Overlay -->
        <div v-if="loading" class="map-loading">
          <Loader2 :size="48" class="spinner" />
          <p>Loading shippers...</p>
        </div>
      </div>

      <!-- Shippers List -->
      <div class="shippers-section">
        <h3 class="section-title">
          <Users :size="24" />
          Active Shippers
        </h3>

        <div v-if="shipperStore.shippers.length === 0" class="empty-state">
          <Truck :size="64" />
          <p>No shippers available</p>
        </div>

        <div v-else class="shippers-grid">
          <div
              v-for="shipper in shipperStore.shippers"
              :key="shipper.id"
              class="shipper-card"
              :class="{ active: isShipperActive(shipper.status) }"
          >
            <div class="shipper-header">
              <div class="shipper-avatar" :class="shipper.status">
                <User :size="24" />
              </div>
              <div class="shipper-info">
                <h4>{{ shipper.user?.name || 'Unknown' }}</h4>
                <p class="phone">
                  <Phone :size="14" />
                  {{ shipper.user?.phone_number || 'N/A' }}
                </p>
              </div>
              <span class="status-badge" :class="shipper.status">
                {{ getStatusLabel(shipper.status) }}
              </span>
            </div>

            <div class="location-section">
              <div v-if="shipper.current_location" class="location-info">
                <p class="location-title">
                  <MapPin :size="16" />
                  Current Location
                </p>
                <div class="coordinates">
                  <span>Lat: {{ shipper.current_location.latitude?.toFixed(6) }}</span>
                  <span>Lng: {{ shipper.current_location.longitude?.toFixed(6) }}</span>
                </div>
                <p class="last-update">
                  <Clock :size="12" />
                  Updated: {{ formatTime(shipper.updatedAt) }}
                </p>
              </div>
              <div v-else class="no-location">
                <MapPin :size="20" />
                <p>Location not available</p>
              </div>
            </div>

            <button
                v-if="shipper.current_location"
                @click="focusOnShipper(shipper)"
                class="btn-focus"
            >
              <Navigation :size="16" />
              Focus on Map
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import mapboxgl from 'mapbox-gl'
import 'mapbox-gl/dist/mapbox-gl.css'
import {
  Truck,
  MapPin,
  Phone,
  User,
  Users,
  Activity,
  Clock,
  Loader2,
  Navigation
} from 'lucide-vue-next'
import Navbar from '../components/Navbar.vue'
import { useShipperStore } from '../stores/shipper'

// Mapbox Access Token - THAY ĐỔI TOKEN CỦA BẠN Ở ĐÂY
mapboxgl.accessToken = 'pk.eyJ1Ijoibmd1eWVudGhhaWR1b25nIiwiYSI6ImNtZm1kZXQ3ODAwcDgyaXE3MnprZTNnM2sifQ.VRjRwGiuFp342CTc1RXOmQ'

const shipperStore = useShipperStore()
const loading = ref(true)
const map = ref(null)
const markers = ref({})
const updateInterval = ref(null)

// Computed
const activeShippersCount = computed(() => {
  return shipperStore.shippers.filter(s =>
      s.status === 'available' || s.status === 'busy'
  ).length
})

// Initialize Map
onMounted(async () => {
  try {
    // Initial fetch
    await shipperStore.fetchShippers()

    // Initialize Mapbox
    initMap()

    // Update shippers every 5 seconds
    updateInterval.value = setInterval(async () => {
      await shipperStore.fetchShippers()
      updateMarkers()
    }, 5000)

  } catch (error) {
    console.error('Failed to initialize tracking:', error)
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  if (updateInterval.value) {
    clearInterval(updateInterval.value)
  }
  if (map.value) {
    map.value.remove()
  }
})

// Watch shippers changes
watch(() => shipperStore.shippers, (newShippers) => {
  if (map.value) {
    updateMarkers()
  }
}, { deep: true })

// Initialize Mapbox Map
function initMap() {
  map.value = new mapboxgl.Map({
    container: 'map',
    style: 'mapbox://styles/mapbox/streets-v12',
    center: [105.8342, 21.0278], // Hanoi, Vietnam
    zoom: 12
  })

  // Add navigation controls
  map.value.addControl(new mapboxgl.NavigationControl(), 'top-right')

  // Add fullscreen control
  map.value.addControl(new mapboxgl.FullscreenControl(), 'top-right')

  map.value.on('load', () => {
    updateMarkers()
  })
}

// Update markers on map
function updateMarkers() {
  if (!map.value) return

  const activeShippers = shipperStore.shippers.filter(
      s => s.current_location && (s.status === 'available' || s.status === 'busy')
  )

  // Remove markers for shippers that no longer exist or are offline
  Object.keys(markers.value).forEach(id => {
    if (!activeShippers.find(s => s.id === id)) {
      markers.value[id].remove()
      delete markers.value[id]
    }
  })

  // Add or update markers
  activeShippers.forEach(shipper => {
    const { latitude, longitude } = shipper.current_location

    if (!latitude || !longitude) return

    if (markers.value[shipper.id]) {
      // Update existing marker position
      markers.value[shipper.id].setLngLat([longitude, latitude])
    } else {
      // Create new marker
      const el = createMarkerElement(shipper)
      const marker = new mapboxgl.Marker(el)
          .setLngLat([longitude, latitude])
          .setPopup(createPopup(shipper))
          .addTo(map.value)

      markers.value[shipper.id] = marker
    }
  })

  // Fit bounds to show all markers
  if (activeShippers.length > 0) {
    const bounds = new mapboxgl.LngLatBounds()
    activeShippers.forEach(shipper => {
      const { latitude, longitude } = shipper.current_location
      bounds.extend([longitude, latitude])
    })
    map.value.fitBounds(bounds, { padding: 50, maxZoom: 15 })
  }
}

// Create custom marker element
function createMarkerElement(shipper) {
  const el = document.createElement('div')
  el.className = `custom-marker ${shipper.status}`
  el.innerHTML = `
    <div class="marker-icon">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M14 18V6a2 2 0 0 0-2-2H4a2 2 0 0 0-2 2v11a1 1 0 0 0 1 1h2"></path>
        <path d="M15 18H9"></path>
        <path d="M19 18h2a1 1 0 0 0 1-1v-3.28a1 1 0 0 0-.684-.948l-1.923-.641a1 1 0 0 1-.578-.579l-1.388-4.165A1 1 0 0 0 16.48 7H14"></path>
        <circle cx="8" cy="18" r="2"></circle>
        <circle cx="16" cy="18" r="2"></circle>
      </svg>
    </div>
    <div class="marker-label">${shipper.user?.name || 'Shipper'}</div>
  `
  return el
}

// Create popup content
function createPopup(shipper) {
  const content = `
    <div class="popup-content">
      <h4>${shipper.user?.name || 'Unknown'}</h4>
      <p><strong>Status:</strong> <span class="status-${shipper.status}">${getStatusLabel(shipper.status)}</span></p>
      <p><strong>Phone:</strong> ${shipper.user?.phone_number || 'N/A'}</p>
      <p><strong>Location:</strong><br/>
        ${shipper.current_location.latitude.toFixed(6)}, ${shipper.current_location.longitude.toFixed(6)}
      </p>
    </div>
  `
  return new mapboxgl.Popup({ offset: 25 }).setHTML(content)
}

// Focus on specific shipper
function focusOnShipper(shipper) {
  if (!map.value || !shipper.current_location) return

  const { latitude, longitude } = shipper.current_location
  map.value.flyTo({
    center: [longitude, latitude],
    zoom: 15,
    duration: 1500
  })

  // Open popup
  if (markers.value[shipper.id]) {
    markers.value[shipper.id].togglePopup()
  }
}

// Helper functions
function isShipperActive(status) {
  return status === 'available' || status === 'busy'
}

function getStatusLabel(status) {
  const labels = {
    available: '🟢 Available',
    busy: '🟡 Busy',
    offline: '🔴 Offline'
  }
  return labels[status] || status
}

function formatTime(timestamp) {
  if (!timestamp) return 'Unknown'
  const date = new Date(timestamp)
  const now = new Date()
  const diff = Math.floor((now - date) / 1000) // seconds

  if (diff < 60) return `${diff}s ago`
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return date.toLocaleDateString()
}
</script>

<style scoped>
.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 2rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}

.header-info {
  display: flex;
  gap: 1rem;
}

.info-badge {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: white;
  padding: 0.75rem 1.25rem;
  border-radius: 0.75rem;
  font-weight: 600;
  color: #1f2937;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

/* Map Container */
.map-container {
  position: relative;
  background: white;
  border-radius: 1rem;
  padding: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  margin-bottom: 2rem;
}

.map-view {
  height: 500px;
  border-radius: 0.75rem;
  overflow: hidden;
}

.map-legend {
  position: absolute;
  top: 2rem;
  left: 2rem;
  background: white;
  padding: 1rem 1.25rem;
  border-radius: 0.75rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  z-index: 10;
}

.map-legend h4 {
  font-size: 0.875rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.75rem 0;
}

.legend-items {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  color: #6b7280;
}

.legend-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.legend-dot.available {
  background: #22c55e;
}

.legend-dot.busy {
  background: #f59e0b;
}

.legend-dot.offline {
  background: #ef4444;
}

.map-loading {
  position: absolute;
  inset: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 0.75rem;
  z-index: 100;
}

.spinner {
  animation: spin 1s linear infinite;
  color: #ea580c;
  margin-bottom: 1rem;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Shippers Section */
.shippers-section {
  margin-top: 2rem;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 1.5rem;
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
}

.empty-state svg {
  color: #d1d5db;
  margin-bottom: 1rem;
}

.empty-state p {
  color: #6b7280;
  margin: 0;
}

.shippers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
}

.shipper-card {
  background: white;
  border-radius: 1rem;
  padding: 1.5rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  transition: all 0.3s;
  border: 2px solid transparent;
}

.shipper-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
}

.shipper-card.active {
  border-color: #22c55e;
}

.shipper-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.shipper-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.shipper-avatar.available {
  background: #dcfce7;
  color: #15803d;
}

.shipper-avatar.busy {
  background: #fef3c7;
  color: #92400e;
}

.shipper-avatar.offline {
  background: #fee2e2;
  color: #991b1b;
}

.shipper-info {
  flex: 1;
}

.shipper-info h4 {
  font-size: 1.125rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.25rem 0;
}

.phone {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: #6b7280;
  font-size: 0.875rem;
  margin: 0;
}

.status-badge {
  padding: 0.375rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 700;
  white-space: nowrap;
}

.status-badge.available {
  background: #dcfce7;
  color: #15803d;
}

.status-badge.busy {
  background: #fef3c7;
  color: #92400e;
}

.status-badge.offline {
  background: #fee2e2;
  color: #991b1b;
}

.location-section {
  margin-bottom: 1rem;
}

.location-info {
  background: #f9fafb;
  border-radius: 0.75rem;
  padding: 1rem;
}

.location-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #1f2937;
  font-weight: 600;
  font-size: 0.9rem;
  margin: 0 0 0.75rem 0;
}

.coordinates {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  color: #6b7280;
  font-size: 0.875rem;
  font-family: 'Courier New', monospace;
}

.last-update {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: #9ca3af;
  font-size: 0.75rem;
  margin-top: 0.5rem;
}

.no-location {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  color: #9ca3af;
  text-align: center;
}

.no-location svg {
  margin-bottom: 0.5rem;
}

.no-location p {
  margin: 0;
  font-size: 0.875rem;
}

.btn-focus {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: #ea580c;
  color: white;
  border: none;
  border-radius: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-focus:hover {
  background: #c2410c;
  transform: scale(1.02);
}

/* Custom Marker Styles */
:global(.custom-marker) {
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
}

:global(.custom-marker .marker-icon) {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
  animation: pulse 2s infinite;
}

:global(.custom-marker.available .marker-icon) {
  background: #22c55e;
  color: white;
}

:global(.custom-marker.busy .marker-icon) {
  background: #f59e0b;
  color: white;
}

:global(.custom-marker .marker-label) {
  margin-top: 0.25rem;
  padding: 0.25rem 0.5rem;
  background: rgba(0,0,0,0.75);
  color: white;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 600;
  white-space: nowrap;
}

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 4px 12px rgba(0,0,0,0.3);
  }
  50% {
    box-shadow: 0 4px 20px rgba(0,0,0,0.5);
  }
}

/* Popup Styles */
:global(.mapboxgl-popup-content) {
  padding: 1rem;
  border-radius: 0.75rem;
  box-shadow: 0 8px 24px rgba(0,0,0,0.2);
}

:global(.popup-content h4) {
  font-size: 1.125rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.75rem 0;
}

:global(.popup-content p) {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0.5rem 0;
}

:global(.popup-content .status-available) {
  color: #15803d;
  font-weight: 600;
}

:global(.popup-content .status-busy) {
  color: #92400e;
  font-weight: 600;
}

/* Responsive */
@media (max-width: 768px) {
  .container {
    padding: 1rem 0.5rem;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
  }

  .map-view {
    height: 350px;
  }

  .map-legend {
    top: 1rem;
    left: 1rem;
    padding: 0.75rem 1rem;
  }

  .shippers-grid {
    grid-template-columns: 1fr;
  }
}
</style>