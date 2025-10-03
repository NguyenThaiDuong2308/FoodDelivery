export const API_BASE_URL = 'http://192.168.237.130:8000'

export const API_ENDPOINTS = {
    // Auth
    LOGIN: '/auth/login',
    REGISTER: '/auth/register',
    LOGOUT: '/auth/logout',

    // User
    USERS: '/user',
    USER_BY_ID: (id) => `/user/${id}`,

    // Restaurant
    RESTAURANTS: '/restaurant',
    RESTAURANT_BY_ID: (id) => `/restaurant/${id}`,
    RESTAURANT_MENU: (id) => `/restaurant/${id}/menu`,

    // Order
    ORDERS: '/order',
    ORDER_BY_ID: (id) => `/order/${id}`,
    CUSTOMER_ORDERS: (id) => `/order/customer/${id}`,

    // Shipper
    SHIPPERS: '/shipper',
    SHIPPER_BY_ID: (id) => `/shipper/${id}`,
    SHIPPER_LOCATION: (id) => `/shipper/${id}/location`
}