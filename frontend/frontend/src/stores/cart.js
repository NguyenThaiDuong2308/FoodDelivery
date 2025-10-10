import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useCartStore = defineStore('cart', () => {
    const items = ref([])
    const restaurantId = ref(null)

    const itemCount = computed(() => items.value.length)
    const totalAmount = computed(() => {
        return items.value.reduce((sum, item) => sum + (item.price * item.quantity), 0)
    })

    function addItem(menuItem) {
        const existingItem = items.value.find(i => i.menu_item_id === menuItem.id)
        if (existingItem) {
            existingItem.quantity++
        } else {
            items.value.push({
                menu_item_id: menuItem.id,
                name: menuItem.name,
                price: menuItem.price,
                quantity: 1
            })
        }
    }

    function removeItem(index) {
        items.value.splice(index, 1)
    }

    function updateQuantity(index, quantity) {
        if (quantity <= 0) {
            removeItem(index)
        } else {
            items.value[index].quantity = quantity
        }
    }

    function clearCart() {
        items.value = []
        restaurantId.value = null
    }

    function setRestaurant(id) {
        restaurantId.value = Number(id)
    }

    return {
        items,
        restaurantId,
        itemCount,
        totalAmount,
        addItem,
        removeItem,
        updateQuantity,
        clearCart,
        setRestaurant
    }
})