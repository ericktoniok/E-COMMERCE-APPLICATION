type ToastType = 'success' | 'error' | 'info'

export function toast(message: string, type: ToastType = 'info', timeout = 3000) {
  const detail = { id: Date.now() + Math.random(), message, type, timeout }
  window.dispatchEvent(new CustomEvent('app:toast', { detail }))
}
