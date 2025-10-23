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

const routes = [
  { path: '/', name: 'home', component: Home },
  { path: '/login', name: 'login', component: Login },
  { path: '/register', name: 'register', component: Register },
  { path: '/products', name: 'products', component: Products },
  { path: '/orders', name: 'orders', component: Orders },
  { path: '/cart', name: 'cart', component: Cart },
  { path: '/checkout', name: 'checkout', component: Checkout },
  { path: '/admin/products', name: 'admin-products', component: AdminProducts },
  { path: '/admin/orders', name: 'admin-orders', component: AdminOrders },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})
