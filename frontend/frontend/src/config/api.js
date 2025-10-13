export const API_BASE_URL = 'http://localhost:8000'

export const API_ENDPOINTS = {
    // Auth
    LOGIN: '/auth/login',
    REGISTER: '/auth/register',
    LOGOUT: '/auth/logout',
    RESET_PASSWORD: 'auth/reset-password',
    FORGOT_PASSWORD: 'auth/forgot-password',
    RESET_FORGOT_PASSWORD: 'auth/reset-forgot-password',
    REFRESH_TOKEN: '/auth/refresh-token',
    // User
    USERS: '/user',
    USER_BY_ID: (id) => `/user/${id}`,

    // Restaurant
    RESTAURANTS: '/restaurant',
    RESTAURANT_BY_ID: (id) => `/restaurant/${id}`,
    RESTAURANT_MENU: (id) => `/restaurant/${id}/menu`,
    MENU_ITEM: (restaurantID, menuItemID) => `/restaurant/${restaurantID}/menu/${menuItemID}`,

    // Order
    ORDERS: '/order',
    ORDER_BY_ID: (id) => `/order/${id}`,
    CUSTOMER_ORDERS: (id) => `/order/customer/${id}`,
    RESTAURANT_ORDERS: (id) => `/order/restaurant/${id}`,
    SHIPPER_ORDERS: (id) => `/order/shipper/${id}`,
    // Shipper
    SHIPPERS: '/shipper',
    SHIPPER_BY_ID: (id) => `/shipper/${id}`,
    SHIPPER_LOCATION: (id) => `/shipper/${id}/location`
}