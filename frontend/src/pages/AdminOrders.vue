<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-xl font-semibold">Admin: Orders</h2>
      <span class="text-xs px-2 py-1 rounded-full"
            :class="sseConnected ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'">
        {{ sseConnected ? 'Live' : 'Offline' }}
      </span>
    </div>
    <div class="flex flex-wrap gap-2 mb-3 text-sm">
      <span class="px-2 py-1 rounded border">All: {{ orders.length }}</span>
      <span class="px-2 py-1 rounded border">Pending: {{ countBy('PENDING') }}</span>
      <span class="px-2 py-1 rounded border">Paid: {{ countBy('PAID') }}</span>
      <span class="px-2 py-1 rounded border">Failed: {{ countBy('FAILED') }}</span>
    </div>
    <div class="mb-4 grid gap-3 sm:grid-cols-6">
      <input v-model="q" type="text" placeholder="Search by order ID or user ID" class="sm:col-span-2 border rounded px-3 py-2 w-full" />
      <select v-model="statusFilter" class="border rounded px-3 py-2 w-full">
        <option value="">All Statuses</option>
        <option value="PENDING">PENDING</option>
        <option value="PAID">PAID</option>
        <option value="FAILED">FAILED</option>
      </select>
      <input v-model="dateFrom" type="date" class="border rounded px-3 py-2 w-full" />
      <input v-model="dateTo" type="date" class="border rounded px-3 py-2 w-full" />
      <div class="flex items-center gap-2">
        <button @click="load" class="border rounded px-3 py-2">Refresh</button>
        <button @click="setLast7Days" class="border rounded px-3 py-2">Last 7 days</button>
        <button @click="setLast24h" class="border rounded px-3 py-2">Recent 24h</button>
        <button @click="exportCsv" class="border rounded px-3 py-2">Export CSV</button>
      </div>
    </div>
    <div v-if="loading" class="space-y-3">
      <div v-for="n in 4" :key="n" class="bg-white border rounded p-3 animate-pulse">
        <div class="h-4 bg-gray-200 w-1/3 mb-2 rounded"></div>
        <div class="h-3 bg-gray-100 w-1/4 rounded"></div>
      </div>
    </div>
    <div v-else-if="!filtered.length" class="text-gray-600">No results. Adjust filters or <button class="underline" @click="clearFilters">clear all</button>.</div>
    <div class="space-y-3">
      <div v-for="o in filtered" :key="o.id" class="bg-white border rounded p-3">
        <div class="flex items-center justify-between">
          <div>
            <div class="font-medium">Order #{{ o.id }}</div>
            <div class="text-sm text-gray-600">User: {{ o.user_id }} | Status: {{ o.status }}</div>
          </div>
          <div class="font-semibold">Total: {{ money(o.total_cents) }}</div>
        </div>
        <div class="mt-2 flex items-center gap-2 text-sm">
          <div class="text-gray-500">{{ fmtDate(o.created_at) }}</div>
          <button @click="copyId(o.id)" class="border rounded px-2 py-1">Copy ID</button>
        </div>
        <ul class="mt-2 list-disc list-inside text-sm">
          <li v-for="it in o.items" :key="it.id">Product {{ it.product_id }} x {{ it.qty }} @ {{ money(it.price_cents) }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed } from 'vue'
import { api } from '../lib/api'
import { toast } from '../lib/toast'
import { money } from '../lib/format'

const orders = ref<any[]>([])
let timer: any = null
const lastOk = ref(true)
let es: EventSource | null = null
const apiBase = import.meta.env.VITE_API_BASE || 'http://localhost:8080'
const sseConnected = ref(false)
const q = ref('')
const statusFilter = ref('')
const dateFrom = ref('')
const dateTo = ref('')
const loading = ref(false)

const filtered = computed(() => {
  const f = orders.value.filter((o: any) => {
    const qOk = !q.value || (String(o.id).includes(q.value) || String(o.user_id||'').includes(q.value))
    const sOk = !statusFilter.value || o.status === statusFilter.value
    const created = o.created_at ? new Date(o.created_at) : null
    let dOk = true
    if (dateFrom.value && created) { dOk = dOk && created >= new Date(dateFrom.value) }
    if (dateTo.value && created) { const end = new Date(dateTo.value); end.setHours(23,59,59,999); dOk = dOk && created <= end }
    return qOk && sOk && dOk
  })
  return f
})

function fmtDate(d: any): string {
  try { return new Date(d).toLocaleString() } catch { return '' }
}

async function copyId(id: number) {
  try { await navigator.clipboard.writeText(String(id)); toast('Order ID copied','success') } catch {}
}

function toDateInputValue(dt: Date) {
  const y = dt.getFullYear()
  const m = String(dt.getMonth()+1).padStart(2,'0')
  const d = String(dt.getDate()).padStart(2,'0')
  return `${y}-${m}-${d}`
}

function setLast7Days() {
  const end = new Date()
  const start = new Date()
  start.setDate(end.getDate()-6)
  dateFrom.value = toDateInputValue(start)
  dateTo.value = toDateInputValue(end)
}

function setLast24h() {
  const end = new Date()
  const start = new Date(end.getTime() - 24*60*60*1000)
  dateFrom.value = toDateInputValue(start)
  dateTo.value = toDateInputValue(end)
}

function exportCsv() {
  const rows = [[
    'id','user_id','status','total_cents','created_at'
  ]]
  for (const o of filtered.value as any[]) {
    rows.push([
      o.id,
      o.user_id ?? '',
      o.status,
      o.total_cents,
      o.created_at || ''
    ])
  }
  const csv = rows.map(r => r.map(v => {
    const s = String(v ?? '')
    if (s.includes(',') || s.includes('"') || s.includes('\n')) return '"'+s.replace(/"/g,'""')+'"'
    return s
  }).join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'orders.csv'
  a.click()
  URL.revokeObjectURL(url)
}

function clearFilters() {
  q.value = ''
  statusFilter.value = ''
  dateFrom.value = ''
  dateTo.value = ''
}

function countBy(status: string): number {
  return orders.value.filter((o:any) => o.status === status).length
}

async function load() {
  try {
    if (!orders.value.length) loading.value = true
    orders.value = await api.adminOrders()
    if (!lastOk.value) { toast('Admin orders reconnected','success'); lastOk.value = true }
  } catch (e) {
    if (lastOk.value) { toast('Failed to load admin orders','error'); lastOk.value = false }
  } finally { loading.value = false }
}

onMounted(() => {
  load()
  // SSE subscription for admin orders
  const token = localStorage.getItem('token')
  if (token) {
    try {
      es = new EventSource(`${apiBase}/api/admin/orders/stream?token=${encodeURIComponent(token)}`)
      es.onopen = () => { sseConnected.value = true }
      es.onmessage = () => { load() }
      es.onerror = () => { sseConnected.value = false /* fallback keeps running */ }
    } catch {}
  }
  // Fallback polling
  timer = setInterval(load, 7000)
})

onUnmounted(() => { if (timer) clearInterval(timer); if (es) es.close(); sseConnected.value = false })
</script>

