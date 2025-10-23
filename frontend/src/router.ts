import { createRouter, createWebHistory } from 'vue-router'
import Home from './pages/Home.vue'
import Login from './pages/Login.vue'
import Register from './pages/Register.vue'
import Products from './pages/Products.vue'
import Orders from './pages/Orders.vue'
import Cart from './pages/Cart.vue'
import Checkout from './pages/Checkout.vue'
import AdminProducts from './pages/AdminProducts.vue'
import AdminOrders from './pages/AdminOrders.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: Home },
    { path: '/login', name: 'login', component: Login },
    { path: '/register', name: 'register', component: Register },
    { path: '/products', name: 'products', component: Products },
    { path: '/orders', name: 'orders', component: Orders, meta: { requiresAuth: true } },
    { path: '/cart', name: 'cart', component: Cart },
    { path: '/checkout', name: 'checkout', component: Checkout, meta: { requiresAuth: true } },
    { path: '/admin/products', name: 'admin-products', component: AdminProducts, meta: { requiresAdmin: true } },
    { path: '/admin/orders', name: 'admin-orders', component: AdminOrders, meta: { requiresAdmin: true } },
  ],
})

router.beforeEach((to: any) => {
  const token = localStorage.getItem('token')
  const role = localStorage.getItem('role')
  if (to.meta?.requiresAuth && !token) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta?.requiresAdmin && role !== 'admin') {
    return { name: 'home' }
  }
})

export default router
