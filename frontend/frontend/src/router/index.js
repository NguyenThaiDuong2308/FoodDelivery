import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/login',
            name: 'Login',
            component: () => import('../views/LoginView.vue')
        },
        {
            path: '/register',
            name: 'Register',
            component: () => import('../views/RegisterView.vue')
        },
        {
          path: '/forgot-password',
          name: 'ForgotPassword',
          component: () => import('../views/ForgotPasswordView.vue')
        },
        {
            path: '/reset-forgot-password',
            name: 'ResetForgotPassword',
            component: () => import('../views/ResetForgotPasswordView.vue')
        },
        {
            path: '/',
            name: 'Home',
            component: () => import('../views/HomeView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/admin/restaurants',
            name: 'RestaurantAdmin',
            component: () => import('../views/RestaurantAdmin.vue'),
            meta: {requiresAuth: true, requiresRole: 'restaurant_admin'}

        },
        {
            path: '/menu/:id',
            name: 'Menu',
            component: () => import('../views/MenuView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/cart',
            name: 'Cart',
            component: () => import('../views/CartView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/orders',
            name: 'Orders',
            component: () => import('../views/OrdersView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/tracking',
            name: 'Tracking',
            component: () => import('../views/TrackingView.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/profile',
            name: 'Profile',
            component: () => import('../views/ProfileView.vue'),
            meta: {requiresAuth: true}
        }
    ]
})

router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()

    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
        next('/login')
    } else if ((to.name === 'Login' || to.name === 'Register') && authStore.isAuthenticated) {
        next('/')
    } else {
        next()
    }
})

export default router