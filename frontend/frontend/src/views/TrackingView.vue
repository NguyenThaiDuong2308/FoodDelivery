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

      <!-- Order Tracking Section -->
      <div v-if="authStore.user && (activeOrders && activeOrders.length > 0)" class="order-tracking-section">
        <h3 class="section-title">
          <Package :size="24" />
          Active Deliveries ({{ activeOrders.length }})
        </h3>
        <div class="order-cards">
          <div
              v-for="order in activeOrders"
              :key="order.id"
              class="order-tracking-card"
              :class="{ selected: selectedOrderId === order.id }"
              @click="selectOrder(order.id)"
          >
            <div class="order-info">
              <h4>Order #{{ order.id }}</h4>
              <span class="status-badge delivering">{{ getStatusLabel(order.status) }}</span>
            </div>
            <div class="tracking-info">
              <p><Store :size="14" /> {{ getRestaurantName(order.restaurant_id) }}</p>
              <p v-if="order.shipper_id"><User :size="14" /> Shipper #{{ order.shipper_id }}</p>
              <p><DollarSign :size="14" /> ${{ order.total_price }}</p>
            </div>
            <!-- Delivery Stage Indicator -->
            <div class="delivery-stage">
              <span :class="{ active: deliveryStages[order.id] === 'to_restaurant' }">
                🏪 To Restaurant
              </span>
              <span :class="{ active: deliveryStages[order.id] === 'to_customer' }">
                📍 To Customer
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Map Container -->
      <div class="map-container">
        <div id="map" class="map-view"></div>

        <!-- Map Legend -->
        <div class="map-legend">
          <h4>Map Legend</h4>
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
              <span class="legend-icon">🏪</span>
              <span>Restaurant</span>
            </div>
            <div class="legend-item">
              <span class="legend-icon">📍</span>
              <span>Customer</span>
            </div>
            <div class="legend-item">
              <div class="legend-line"></div>
              <span>Delivery Route</span>
            </div>
          </div>
        </div>

        <!-- Loading Overlay -->
        <div v-if="loading" class="map-loading">
          <Loader2 :size="48" class="spinner" />
          <p>Loading tracking data...</p>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="map-error">
          <AlertCircle :size="48" />
          <p>{{ error }}</p>
          <button @click="retryFetch" class="btn-retry">Try Again</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import mapboxgl from 'mapbox-gl'
import 'mapbox-gl/dist/mapbox-gl.css'
import {
  Truck,
  User,
  Activity,
  Clock,
  Loader2,
  AlertCircle,
  Package,
  Store,
  DollarSign
} from 'lucide-vue-next'
import Navbar from '../components/Navbar.vue'
import { useShipperStore } from '../stores/shipper'
import { useOrderStore } from '../stores/order'
import { useRestaurantStore } from '../stores/restaurant'
import { useAuthStore } from '../stores/auth'
import {useUserStore} from "@/stores/user";

// Mapbox Access Token
mapboxgl.accessToken = 'pk.eyJ1Ijoibmd1eWVudGhhaWR1b25nIiwiYSI6ImNtZm1kZXQ3ODAwcDgyaXE3MnprZTNnM2sifQ.VRjRwGiuFp342CTc1RXOmQ'

const API_BASE_URL = 'http://localhost:8000'
//const API_BASE_URL = 'http://192.168.237.130:8000'

const shipperStore = useShipperStore()
const orderStore = useOrderStore()
const restaurantStore = useRestaurantStore()
const authStore = useAuthStore()
const userStore = useUserStore()

const loading = ref(true)
const error = ref(null)
const map = ref(null)
const markers = ref({})
const restaurantMarkers = ref({})
const customerMarkers = ref({})
const routeLines = ref({})
const updateInterval = ref(null)
const shipperLocations = ref({})
const locationUpdates = ref({})
const selectedOrderId = ref(null)

// Cached data
const restaurantLocations = ref({})
const customerLocations = ref({})

// Mock location simulation
const shipperIntervals = ref({})
const shipperRoutes = ref({})

const deliveryStages = ref({}) // { orderId: 'to_restaurant' | 'to_customer' | 'completed' }

// Base fallback location (Hanoi center)
const HANOI_CENTER = { lat: 21.0285, lng: 105.8542 }
const ARRIVAL_THRESHOLD = 50
// Computed
const activeShippersCount = computed(() => {
  return shipperStore.shippers.filter(s =>
      s.status === 'available' || s.status === 'busy'
  ).length
})

const activeOrders = computed(() => {
  if (!authStore.user) return []

  const orders = orderStore.orders || []
  const restaurants = restaurantStore.restaurants || []
  const shippers = shipperStore.shippers || []

  return orders.filter(order => {
    if (order.status !== 'delivering') return false

    const role = authStore.user.role
    if (role === 'customer') {
      return order.customer_id === authStore.user.id
    } else if (role === 'restaurant_admin') {
      const userRestaurant = restaurants.find(r => r.manager_id === authStore.user.id)
      return userRestaurant && order.restaurant_id === userRestaurant.id
    } else if (role === 'shipper') {
      const shipperRecord = shippers.find(s => s.user_id === authStore.user.id)
      return shipperRecord && order.shipper_id === shipperRecord.id
    }
    return false
  })
})

// Geocode address using Mapbox Geocoding API
async function geocodeAddress(address) {
  if (!address) return null

  try {
    const response = await fetch(
        `https://api.mapbox.com/geocoding/v5/mapbox.places/${encodeURIComponent(address)}.json?access_token=${mapboxgl.accessToken}&country=VN&proximity=105.8342,21.0278`
    )
    const data = await response.json()

    if (data.features && data.features.length > 0) {
      const [lng, lat] = data.features[0].center
      return { lat, lng, address: data.features[0].place_name }
    }
  } catch (err) {
    console.error('Geocoding error:', err)
  }
  return null
}

// Get restaurant location from API
async function getRestaurantLocation(restaurantId) {
  if (restaurantLocations.value[restaurantId]) {
    return restaurantLocations.value[restaurantId]
  }

  try {
    const restaurant = await restaurantStore.fetchRestaurantById(restaurantId)

    if (restaurant && restaurant.address) {
      const location = await geocodeAddress(restaurant.address)
      if (location) {
        restaurantLocations.value[restaurantId] = {
          ...location,
          name: restaurant.name
        }
        return restaurantLocations.value[restaurantId]
      }
    }
  } catch (err) {
    console.error(`Failed to get restaurant ${restaurantId} location:`, err)
  }

  // Fallback
  return { ...HANOI_CENTER, name: `Restaurant #${restaurantId}` }
}

// Get customer location from API
async function getCustomerLocation(customerId) {
  if (customerLocations.value[customerId]) {
    return customerLocations.value[customerId]
  }

  try {
    const locationData = await userStore.fetchUserLocation(customerId)
    if (locationData && locationData.address) {
      const location = await geocodeAddress(locationData.address)
      if (location) {
        customerLocations.value[customerId] = {
          ...location,
          name: `Customer #${customerId}`
        }
        return customerLocations.value[customerId]
      }
    }
  } catch (err) {
    console.error(`Failed to get customer ${customerId} location:`, err)
  }

  return { ...HANOI_CENTER, name: `Customer #${customerId}` }
}


// Get restaurant name
function getRestaurantName(restaurantId) {
  return restaurantLocations.value[restaurantId]?.name || `Restaurant #${restaurantId}`
}

// Fetch route from Mapbox Directions API
async function getMapboxRoute(start, end) {
  try {
    const response = await fetch(
        `https://api.mapbox.com/directions/v5/mapbox/driving/${start.lng},${start.lat};${end.lng},${end.lat}?geometries=geojson&access_token=${mapboxgl.accessToken}`
    )
    const data = await response.json()

    if (data.routes && data.routes.length > 0) {
      const route = data.routes[0]
      const coordinates = route.geometry.coordinates

      // Convert to lat/lng format
      return coordinates.map(coord => ({
        lng: coord[0],
        lat: coord[1]
      }))
    }
  } catch (err) {
    console.error('Failed to fetch route:', err)
  }
  return null
}

// Parse WKT
function parseWKT(wkt) {
  if (!wkt) return null

  const match = wkt.match(/POINT\(([0-9.-]+)\s+([0-9.-]+)\)/)
  if (match) {
    return {
      longitude: parseFloat(match[1]),
      latitude: parseFloat(match[2]),
      wkt: wkt
    }
  }
  return null
}

// Send location to backend
async function sendLocationToBackend(shipperId, location) {
  try {
    const response = await fetch(`${API_BASE_URL}/shipper/${shipperId}/location`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        latitude: location.latitude,
        longitude: location.longitude
      })
    })

    if (response.ok) {
      return true
    }
    return false
  } catch (err) {
    console.error(`Failed to send location for shipper ${shipperId}:`, err)
    return false
  }
}

// Start route-based tracking
async function startRouteBasedTracking(shipper, order) {
  if (shipperIntervals.value[shipper.id]) return

  // Initialize delivery stage
  if (!deliveryStages.value[order.id]) {
    deliveryStages.value[order.id] = 'to_restaurant'
  }

  const restaurantLoc = await getRestaurantLocation(order.restaurant_id)
  const customerLoc = await getCustomerLocation(order.customer_id)

  const currentLoc = shipperLocations.value[shipper.id] || HANOI_CENTER

  // Get route to restaurant first
  let routeWaypoints = await getMapboxRoute(currentLoc, restaurantLoc)

  if (!routeWaypoints || routeWaypoints.length === 0) {
    console.error('Failed to get route, falling back to random movement')
    startRandomLocationUpdates(shipper)
    return
  }

  shipperRoutes.value[shipper.id] = {
    waypoints: routeWaypoints,
    currentIndex: 0,
    orderId: order.id,
    restaurantLoc,
    customerLoc,
    stage: 'to_restaurant'
  }

  // Set initial location
  const initialPoint = routeWaypoints[0]
  shipperLocations.value[shipper.id] = {
    latitude: initialPoint.lat,
    longitude: initialPoint.lng,
    wkt: `POINT(${initialPoint.lng} ${initialPoint.lat})`
  }
  locationUpdates.value[shipper.id] = new Date()
  sendLocationToBackend(shipper.id, shipperLocations.value[shipper.id])

  // Update every 5 seconds
  shipperIntervals.value[shipper.id] = setInterval(async () => {
    const route = shipperRoutes.value[shipper.id]
    if (!route) return

    // Move to next waypoint
    route.currentIndex++

    // Check if reached current destination
    if (route.currentIndex >= route.waypoints.length) {
      if (route.stage === 'to_restaurant') {
        // Reached restaurant, now go to customer
        console.log(`Shipper ${shipper.id} reached restaurant for order ${order.id}`)

        deliveryStages.value[order.id] = 'to_customer'
        route.stage = 'to_customer'

        // Get route to customer
        const routeToCustomer = await getMapboxRoute(route.restaurantLoc, route.customerLoc)
        if (routeToCustomer && routeToCustomer.length > 0) {
          route.waypoints = routeToCustomer
          route.currentIndex = 0
        } else {
          console.error('Failed to get route to customer')
          return
        }
      } else if (route.stage === 'to_customer') {
        // Reached customer, mark as completed
        console.log(`Shipper ${shipper.id} reached customer for order ${order.id}`)

        deliveryStages.value[order.id] = 'completed'

        // Auto-update order status to 'completed' if user is customer
        if (authStore.user?.role === 'customer' && order.customer_id === authStore.user.id) {
          try {
            await orderStore.updateOrderStatus(order.id, 'completed')
            console.log(`Order ${order.id} automatically marked as completed`)

            showNotification('✅ Your order has been delivered!', 'success')
          } catch (err) {
            console.error('Failed to auto-update order status:', err)
          }
        }

        // Stop tracking this shipper
        stopShipperLocationUpdates(shipper.id)
        return
      }
    }

    const currentPoint = route.waypoints[route.currentIndex]

    shipperLocations.value[shipper.id] = {
      latitude: currentPoint.lat,
      longitude: currentPoint.lng,
      wkt: `POINT(${currentPoint.lng} ${currentPoint.lat})`
    }
    locationUpdates.value[shipper.id] = new Date()
    sendLocationToBackend(shipper.id, shipperLocations.value[shipper.id])
  }, 5000)
}

// Start random location updates
function startRandomLocationUpdates(shipper) {
  if (shipperIntervals.value[shipper.id]) return

  const radiusMeters = 500
  const radiusInDegrees = radiusMeters / 111000

  const initialLocation = {
    latitude: HANOI_CENTER.lat + (Math.random() - 0.5) * radiusInDegrees * 2,
    longitude: HANOI_CENTER.lng + (Math.random() - 0.5) * radiusInDegrees * 2
  }

  shipperLocations.value[shipper.id] = {
    ...initialLocation,
    wkt: `POINT(${initialLocation.longitude} ${initialLocation.latitude})`
  }
  locationUpdates.value[shipper.id] = new Date()
  sendLocationToBackend(shipper.id, initialLocation)

  shipperIntervals.value[shipper.id] = setInterval(() => {
    const newLocation = {
      latitude: HANOI_CENTER.lat + (Math.random() - 0.5) * radiusInDegrees * 2,
      longitude: HANOI_CENTER.lng + (Math.random() - 0.5) * radiusInDegrees * 2
    }
    shipperLocations.value[shipper.id] = {
      ...newLocation,
      wkt: `POINT(${newLocation.longitude} ${newLocation.latitude})`
    }
    locationUpdates.value[shipper.id] = new Date()
    sendLocationToBackend(shipper.id, newLocation)
  }, 5000)
}

// Stop tracking
function stopShipperLocationUpdates(shipperId) {
  if (shipperIntervals.value[shipperId]) {
    clearInterval(shipperIntervals.value[shipperId])
    delete shipperIntervals.value[shipperId]
    delete shipperRoutes.value[shipperId]
  }
}

// Fetch shipper location
async function fetchShipperLocation(shipperId) {
  try {
    const response = await fetch(`${API_BASE_URL}/shipper/${shipperId}/location`)

    if (!response.ok) {
      if (response.status === 404) {
        return null
      }
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const data = await response.json()

    if (data && data.location) {
      const parsedLocation = parseWKT(data.location)

      if (parsedLocation) {
        shipperLocations.value[shipperId] = parsedLocation
        locationUpdates.value[shipperId] = new Date()
        return parsedLocation
      }
    }
  } catch (err) {
    console.error(`Failed to fetch location for shipper ${shipperId}:`, err)
  }
  return null
}

// Select order
async function selectOrder(orderId) {
  selectedOrderId.value = orderId
  const order = activeOrders.value.find(o => o.id === orderId)
  if (order && order.shipper_id) {
    await focusOnOrder(order)
  }
}

// Focus on order route
async function focusOnOrder(order) {
  if (!map.value) return

  const restaurantLoc = await getRestaurantLocation(order.restaurant_id)
  const customerLoc = await getCustomerLocation(order.customer_id)
  const shipperLoc = shipperLocations.value[order.shipper_id]

  if (shipperLoc) {
    const bounds = new mapboxgl.LngLatBounds()
    bounds.extend([restaurantLoc.lng, restaurantLoc.lat])
    bounds.extend([customerLoc.lng, customerLoc.lat])
    bounds.extend([shipperLoc.longitude, shipperLoc.latitude])

    map.value.fitBounds(bounds, { padding: 100 })
  }
}

// Draw route on map
async function drawRouteOnMap(orderId, restaurantLoc, customerLoc, shipperLoc) {
  if (!map.value || !map.value.isStyleLoaded()) return

  const routeId = `route-${orderId}`

  // Remove existing
  if (map.value.getLayer(routeId)) {
    map.value.removeLayer(routeId)
  }
  if (map.value.getSource(routeId)) {
    map.value.removeSource(routeId)
  }

  const stage = deliveryStages.value[orderId]
  let targetLoc = stage === 'to_customer' ? customerLoc : restaurantLoc

  // Get route from shipper to target
  const routeWaypoints = await getMapboxRoute(
      { lat: shipperLoc.latitude, lng: shipperLoc.longitude },
      targetLoc
  )

  if (routeWaypoints && routeWaypoints.length > 0) {
    const coordinates = routeWaypoints.map(wp => [wp.lng, wp.lat])

    map.value.addSource(routeId, {
      type: 'geojson',
      data: {
        type: 'Feature',
        properties: {},
        geometry: {
          type: 'LineString',
          coordinates: coordinates
        }
      }
    })

    map.value.addLayer({
      id: routeId,
      type: 'line',
      source: routeId,
      layout: {
        'line-join': 'round',
        'line-cap': 'round'
      },
      paint: {
        'line-color': '#ea580c',
        'line-width': 4
      }
    })
  }

  // Add restaurant marker
  if (!restaurantMarkers.value[`restaurant-${orderId}`]) {
    restaurantMarkers.value[`restaurant-${orderId}`] = new mapboxgl.Marker({
      color: '#f59e0b',
      scale: 1.2
    })
        .setLngLat([restaurantLoc.lng, restaurantLoc.lat])
        .setPopup(new mapboxgl.Popup().setHTML(`
        <div class="custom-popup">
          <h4>🏪 ${restaurantLoc.name}</h4>
          <p>${restaurantLoc.address || ''}</p>
        </div>
      `))
        .addTo(map.value)
  }

  // Add customer marker
  if (!customerMarkers.value[`customer-${orderId}`]) {
    customerMarkers.value[`customer-${orderId}`] = new mapboxgl.Marker({
      color: '#22c55e',
      scale: 1.2
    })
        .setLngLat([customerLoc.lng, customerLoc.lat])
        .setPopup(new mapboxgl.Popup().setHTML(`
        <div class="custom-popup">
          <h4>📍 ${customerLoc.name}</h4>
          <p>${customerLoc.address || ''}</p>
        </div>
      `))
        .addTo(map.value)
  }

  routeLines.value[orderId] = routeId
}

// Show notification
function showNotification(message, type = 'success') {
  const notification = document.createElement('div')
  notification.className = `notification ${type}`
  notification.textContent = message
  document.body.appendChild(notification)

  setTimeout(() => notification.classList.add('show'), 10)

  setTimeout(() => {
    notification.classList.remove('show')
    setTimeout(() => document.body.removeChild(notification), 300)
  }, 5000)
}

// Initialize
onMounted(async () => {
  try {
    await shipperStore.fetchShippers()

    if (authStore.user) {
      const role = authStore.user.role
      if (role === 'customer') {
        await orderStore.fetchCustomerOrders(authStore.user.id)
      } else if (role === 'restaurant_admin') {
        await restaurantStore.fetchRestaurants()
        const userRestaurant = restaurantStore.restaurants.find(r => r.manager_id === authStore.user.id)
        if (userRestaurant) {
          await orderStore.fetchRestaurantOrders(userRestaurant.id)
        }
      } else if (role === 'shipper') {
        const shipperRecord = shipperStore.shippers.find(s => s.user_id === authStore.user.id)
        if (shipperRecord) {
          await orderStore.fetchShipperOrders(shipperRecord.id)
        }
      }
    }

    initMap()

    // Start tracking
    for (const shipper of shipperStore.shippers) {
      if (shipper.status === 'available' || shipper.status === 'busy') {
        const shipperOrder = orderStore.orders.find(o =>
            o.shipper_id === shipper.id && o.status === 'delivering'
        )

        if (shipperOrder) {
          await startRouteBasedTracking(shipper, shipperOrder)
        } else {
          startRandomLocationUpdates(shipper)
        }
      }
    }

    updateInterval.value = setInterval(async () => {
      await shipperStore.fetchShippers()

      if (authStore.user) {
        const role = authStore.user.role
        if (role === 'customer') {
          await orderStore.fetchCustomerOrders(authStore.user.id)
        } else if (role === 'restaurant_admin') {
          const userRestaurant = restaurantStore.restaurants.find(r => r.manager_id === authStore.user.id)
          if (userRestaurant) {
            await orderStore.fetchRestaurantOrders(userRestaurant.id)
          }
        } else if (role === 'shipper') {
          const shipperRecord = shipperStore.shippers.find(s => s.user_id === authStore.user.id)
          if (shipperRecord) {
            await orderStore.fetchShipperOrders(shipperRecord.id)
          }
        }
      }

      updateMarkers()
    }, 30000)

  } catch (err) {
    console.error('Failed to initialize tracking:', err)
    error.value = 'Failed to load tracking data'
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  Object.keys(shipperIntervals.value).forEach(shipperId => {
    stopShipperLocationUpdates(parseInt(shipperId))
  })

  if (updateInterval.value) {
    clearInterval(updateInterval.value)
  }
  if (map.value) {
    map.value.remove()
  }
})

async function retryFetch() {
  error.value = null
  loading.value = true
  try {
    await shipperStore.fetchShippers()
    if (map.value) {
      updateMarkers()
    } else {
      initMap()
    }
  } catch (err) {
    error.value = 'Failed to load tracking data'
  } finally {
    loading.value = false
  }
}

function initMap() {
  map.value = new mapboxgl.Map({
    container: 'map',
    style: 'mapbox://styles/mapbox/streets-v12',
    center: [105.8342, 21.0278],
    zoom: 12
  })

  map.value.addControl(new mapboxgl.NavigationControl(), 'top-right')
  map.value.addControl(new mapboxgl.FullscreenControl(), 'top-right')

  map.value.on('load', () => {
    updateMarkers()

    setInterval(() => {
      updateMarkers()
      updateRoutes()
    }, 2000)
  })
}

function updateMarkers() {
  if (!map.value) return

  const activeShippers = shipperStore.shippers.filter(s => {
    const hasLocation = shipperLocations.value[s.id]
    const isActive = s.status === 'available' || s.status === 'busy'
    return hasLocation && isActive
  })

  Object.keys(markers.value).forEach(id => {
    if (!activeShippers.find(s => s.id == id)) {
      markers.value[id].remove()
      delete markers.value[id]
    }
  })

  activeShippers.forEach(shipper => {
    const location = shipperLocations.value[shipper.id]
    if (!location || !location.latitude || !location.longitude) return

    const { latitude, longitude } = location

    if (markers.value[shipper.id]) {
      markers.value[shipper.id].setLngLat([longitude, latitude])
    } else {
      const el = createMarkerElement(shipper)
      const marker = new mapboxgl.Marker(el)
          .setLngLat([longitude, latitude])
          .addTo(map.value)

      markers.value[shipper.id] = marker
    }
  })
}

async function updateRoutes() {
  for (const order of activeOrders.value) {
    if (order.shipper_id && shipperLocations.value[order.shipper_id]) {
      const restaurantLoc = await getRestaurantLocation(order.restaurant_id)
      const customerLoc = await getCustomerLocation(order.customer_id)
      const shipperLoc = shipperLocations.value[order.shipper_id]

      await drawRouteOnMap(order.id, restaurantLoc, customerLoc, shipperLoc)
    }
  }
}

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
  `
  return el
}

function getStatusLabel(status) {
  const labels = {
    available: '🟢 Available',
    busy: '🟡 Busy',
    offline: '🔴 Offline',
    delivering: '🚚 Delivering',
    created: '📝 Created',
    completed: '✅ Completed',
    cancelled: '❌ Cancelled'
  }
  return labels[status] || status
}
</script>

<style scoped>

.delivery-stage {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid #e5e7eb;
}

.delivery-stage span {
  flex: 1;
  text-align: center;
  padding: 0.5rem;
  background: #f3f4f6;
  color: #9ca3af;
  border-radius: 0.5rem;
  font-size: 0.75rem;
  font-weight: 600;
  transition: all 0.3s;
}

.delivery-stage span.active {
  background: #fef3c7;
  color: #92400e;
  box-shadow: 0 2px 4px rgba(234, 88, 12, 0.2);
}

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

/* Order Tracking Section */
.order-tracking-section {
  margin-bottom: 2rem;
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

.order-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1rem;
}

.order-tracking-card {
  background: white;
  border-radius: 0.75rem;
  padding: 1.25rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
}

.order-tracking-card:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  transform: translateY(-2px);
}

.order-tracking-card.selected {
  border-color: #ea580c;
  box-shadow: 0 4px 16px rgba(234, 88, 12, 0.3);
}

.order-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.order-info h4 {
  margin: 0;
  font-size: 1.125rem;
  color: #1f2937;
  font-weight: 700;
}

.status-badge {
  padding: 0.375rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: capitalize;
}

.status-badge.delivering {
  background: #fef3c7;
  color: #92400e;
}

.tracking-info {
  font-size: 0.875rem;
  color: #6b7280;
}

.tracking-info p {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0.5rem 0;
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
  height: 600px;
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

.legend-icon {
  font-size: 1.25rem;
}

.legend-line {
  width: 30px;
  height: 3px;
  background: #ea580c;
  border-radius: 2px;
}

.map-loading,
.map-error {
  position: absolute;
  inset: 0;
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 0.75rem;
  z-index: 100;
  padding: 2rem;
}

.map-error svg {
  color: #ef4444;
  margin-bottom: 1rem;
}

.map-error p {
  color: #6b7280;
  margin-bottom: 1.5rem;
  font-size: 1.125rem;
}

.btn-retry {
  padding: 0.75rem 1.5rem;
  background: #ea580c;
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-retry:hover {
  background: #c2410c;
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

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 4px 12px rgba(0,0,0,0.3);
  }
  50% {
    box-shadow: 0 4px 20px rgba(0,0,0,0.5);
  }
}

/* Popup Styles */
:global(.custom-popup h4) {
  font-size: 1rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 0.5rem 0;
}

:global(.custom-popup p) {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0;
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

  .header-info {
    width: 100%;
    flex-direction: column;
  }

  .info-badge {
    width: 100%;
    justify-content: center;
  }

  .map-view {
    height: 400px;
  }

  .map-legend {
    top: 1rem;
    left: 1rem;
    padding: 0.75rem 1rem;
  }

  .order-cards {
    grid-template-columns: 1fr;
  }
}
</style>