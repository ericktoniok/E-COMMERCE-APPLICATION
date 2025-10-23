import { ref } from 'vue'

const base = import.meta.env.VITE_API_BASE || 'http://localhost:8080'
export const token = ref<string | null>(null)

export function setToken(t: string | null) {
  token.value = t
}

async function request(path: string, opts: RequestInit = {}) {
  const headers: Record<string, string> = { 'Content-Type': 'application/json', ...(opts.headers as any) }
  if (token.value) headers['Authorization'] = `Bearer ${token.value}`
  const res = await fetch(`${base}${path}`, { ...opts, headers })
  if (!res.ok) throw new Error(await res.text())
  return res.headers.get('content-type')?.includes('application/json') ? res.json() : res.text()
}

export const api = {
  // auth
  register: (email: string, password: string) => request('/api/auth/register', { method: 'POST', body: JSON.stringify({ email, password }) }),
  login: (email: string, password: string) => request('/api/auth/login', { method: 'POST', body: JSON.stringify({ email, password }) }),

  // products
  products: () => request('/api/products'),
  createProduct: (p: any) => request('/api/products', { method: 'POST', body: JSON.stringify(p) }),
  updateProduct: (id: number, p: any) => request(`/api/products/${id}`, { method: 'PUT', body: JSON.stringify(p) }),
  deleteProduct: (id: number) => request(`/api/products/${id}`, { method: 'DELETE' }),

  // orders
  checkout: (items: Array<{product_id:number; qty:number}>) => request('/api/cart/checkout', { method: 'POST', body: JSON.stringify({ items }) }),
  myOrders: () => request('/api/orders/me'),
  adminOrders: () => request('/api/admin/orders/'),
}
