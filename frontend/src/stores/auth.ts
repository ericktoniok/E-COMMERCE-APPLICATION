import { defineStore } from 'pinia'
import { api, setToken, token as tokenRef } from '../lib/api'

export const useAuth = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') as string | null,
    role: (localStorage.getItem('role') as string | null) || null,
    email: (localStorage.getItem('email') as string | null) || null,
  }),
  actions: {
    async register(email: string, password: string) {
      const res = await api.register(email, password)
      this.token = res.token; this.role = res.role; this.email = email
      localStorage.setItem('token', this.token!)
      localStorage.setItem('role', this.role!)
      localStorage.setItem('email', this.email!)
      setToken(this.token)
    },
    async login(email: string, password: string) {
      const res = await api.login(email, password)
      this.token = res.token; this.role = res.role; this.email = email
      localStorage.setItem('token', this.token!)
      localStorage.setItem('role', this.role!)
      localStorage.setItem('email', this.email!)
      setToken(this.token)
    },
    logout() {
      this.token = null; this.role = null; this.email = null
      localStorage.removeItem('token')
      localStorage.removeItem('role')
      localStorage.removeItem('email')
      setToken(null)
    },
    init() {
      setToken(this.token)
    }
  }
})
