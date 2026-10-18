<template>
  <div class="relative rounded-2xl overflow-hidden border border-white/10 bg-white/5">
    <div ref="mapContainer" class="w-full h-[500px]"></div>
  </div>
</template>

<script setup lang="ts">
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

interface Place {
  id: number
  name: string
  lat: number
  lng: number
  description: string
  emoji: string
  order_index: number
}

const props = defineProps<{
  places: Place[]
}>()

const mapContainer = ref<HTMLElement | null>(null)
let map: L.Map | null = null

const createIcon = (emoji: string, index: number) => {
  return L.divIcon({
    className: 'custom-marker',
    html: `
      <div class="marker-wrapper">
        <div class="marker-pin"><span>${emoji}</span></div>
        <div class="marker-number">${index}</div>
      </div>
    `,
    iconSize: [40, 50],
    iconAnchor: [20, 50],
    popupAnchor: [0, -50],
  })
}

onMounted(() => {
  if (!mapContainer.value || !props.places.length) return

  const avgLat = props.places.reduce((sum, p) => sum + p.lat, 0) / props.places.length
  const avgLng = props.places.reduce((sum, p) => sum + p.lng, 0) / props.places.length

  map = L.map(mapContainer.value, {
    attributionControl: true,
  }).setView([avgLat, avgLng], 6)

  map.attributionControl.setPrefix(false)

  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors',
    maxZoom: 19,
  }).addTo(map)

  const coords: [number, number][] = []

  props.places
    .slice()
    .sort((a, b) => a.order_index - b.order_index)
    .forEach((place, index) => {
      const coord: [number, number] = [place.lat, place.lng]
      coords.push(coord)

      const marker = L.marker(coord, {
        icon: createIcon(place.emoji || '📍', index + 1),
      }).addTo(map!)

      marker.bindPopup(`
        <div style="font-family: Inter, sans-serif; min-width: 180px;">
          <div style="font-size: 24px; margin-bottom: 4px;">${place.emoji || '📍'}</div>
          <div style="font-weight: 600; font-size: 16px; margin-bottom: 4px;">${place.name}</div>
          <div style="color: #94a3b8; font-size: 13px;">${place.description || ''}</div>
        </div>
      `)
    })

  if (coords.length > 1) {
    L.polyline(coords, {
      color: '#00dc82',
      weight: 3,
      opacity: 0.7,
      dashArray: '10, 10',
    }).addTo(map)
  }

  if (coords.length > 0) {
    const bounds = L.latLngBounds(coords)
    map.fitBounds(bounds, { padding: [50, 50] })
  }
})

onBeforeUnmount(() => {
  if (map) {
    map.remove()
    map = null
  }
})
</script>

<style>
.custom-marker {
  background: transparent;
  border: none;
}

.marker-wrapper {
  position: relative;
  width: 40px;
  height: 50px;
}

.marker-pin {
  width: 40px;
  height: 40px;
  border-radius: 50% 50% 50% 0;
  background: linear-gradient(135deg, #00dc82, #00b86b);
  transform: rotate(-45deg);
  box-shadow: 0 4px 12px rgba(0, 220, 130, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  position: absolute;
  top: 0;
  left: 0;
}

.marker-pin span {
  transform: rotate(45deg);
  font-size: 18px;
}

.marker-number {
  position: absolute;
  top: -5px;
  right: -5px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #020420;
  color: #00dc82;
  font-size: 11px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid #00dc82;
  font-family: Inter, sans-serif;
}

.leaflet-popup-content-wrapper {
  background: #0a0e1f;
  color: #fff;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.leaflet-popup-tip {
  background: #0a0e1f;
}

.leaflet-popup-content {
  margin: 12px;
}

.leaflet-popup-close-button {
  color: #94a3b8 !important;
}
</style>
