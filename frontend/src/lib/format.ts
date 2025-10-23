export function money(cents: number): string {
  try {
    return new Intl.NumberFormat(undefined, { style: 'currency', currency: 'USD' }).format(cents / 100)
  } catch {
    return (cents / 100).toFixed(2)
  }
}
