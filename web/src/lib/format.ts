export function formatBytes(bytes: number | undefined): string {
  if (bytes == null || bytes === 0) {
    return '—'
  }
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let value = bytes
  let unitIndex = 0
  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024
    unitIndex += 1
  }
  return `${value.toFixed(unitIndex === 0 ? 0 : 1)} ${units[unitIndex]}`
}

export function bytesToGb(bytes: number | undefined): string {
  if (bytes == null) {
    return '—'
  }
  return `${(bytes / 1024 ** 3).toFixed(1)} GB`
}
