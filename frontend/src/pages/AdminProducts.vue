<template>
  <div>
    <h2 class="text-xl font-semibold mb-4">Admin: Products</h2>
    <form @submit.prevent="create" class="bg-white border rounded p-3 mb-4 grid grid-cols-1 md:grid-cols-5 gap-2">
      <input v-model="form.name" placeholder="Name" class="input" required />
      <input v-model.number="form.price_cents" type="number" placeholder="Price (cents)" class="input" required />
      <input v-model.number="form.stock" type="number" placeholder="Stock" class="input" required />
      <input v-model="form.description" placeholder="Description" class="input md:col-span-2" />
      <button class="btn-primary md:col-span-5" :disabled="creating">{{ creating ? 'Creating…' : 'Create' }}</button>
      <p v-if="err" class="text-red-600 text-sm md:col-span-5">{{ err }}</p>
    </form>

    <div class="space-y-3">
      <div v-for="p in products" :key="p.id" class="bg-white border rounded p-3 flex items-center justify-between">
        <div class="flex-1">
          <div class="font-medium">#{{ p.id }} - {{ p.name }}</div>
          <div class="text-sm text-gray-600">{{ p.description }}</div>
          <div class="text-sm">Price: {{ money(p.price_cents) }} | Stock: {{ p.stock }}</div>
          <div class="mt-2 flex items-center gap-2">
            <input type="file" @change="e=>upload(p.id, (e.target as HTMLInputElement).files?.[0]||null)" />
          </div>
        </div>
        <div class="flex items-center gap-2">
          <button class="btn" :disabled="savingId===p.id" @click="save(p)">{{ savingId===p.id ? 'Saving…' : 'Save' }}</button>
          <button class="btn-danger" :disabled="deletingId===p.id" @click="remove(p.id)">{{ deletingId===p.id ? 'Deleting…' : 'Delete' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { api } from '../lib/api'
import { toast } from '../lib/toast'
import { money } from '../lib/format'

const products = ref<any[]>([])
const err = ref('')
const form = reactive({ name:'', price_cents:0, stock:0, description:'' })
const creating = ref(false)
const savingId = ref<number | null>(null)
const deletingId = ref<number | null>(null)

async function load() { products.value = await api.products() }

async function create() {
  err.value = ''
  creating.value = true
  try {
    await api.createProduct(form)
    form.name=''; form.price_cents=0; form.stock=0; form.description=''
    await load()
    toast('Product created','success')
  } catch (e:any) {
    err.value = 'Create failed'
    toast('Create failed','error')
  }
  creating.value = false
}
async function save(p:any) {
  savingId.value = p.id
  try { await api.updateProduct(p.id, p); await load(); toast('Product saved','success') } catch { toast('Save failed','error') }
  savingId.value = null
}
async function remove(id:number) {
  deletingId.value = id
  try { await api.deleteProduct(id); await load(); toast('Product deleted','success') } catch { toast('Delete failed','error') }
  deletingId.value = null
}
async function upload(id:number, file: File|null) { if (!file) return; try { await api.uploadProductImage(id, file); await load(); toast('Image uploaded','success') } catch { toast('Image upload failed','error') }
}

onMounted(load)
</script>
<style scoped>
.input { @apply border rounded px-3 py-2 w-full; }
.btn { @apply bg-gray-200 px-3 py-1 rounded; }
.btn-primary { @apply bg-blue-600 text-white px-3 py-1 rounded; }
.btn-danger { @apply bg-red-600 text-white px-3 py-1 rounded; }
</style>
