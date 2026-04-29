<template>
  <NuxtLink
    :to="`/customers/${customer.id}`"
    class="customer-card-link"
  >
    <div class="card customer-card">
      <div class="row">
        <div>
          <div class="card-title">
            {{ customer.first_name }} {{ customer.last_name }}
          </div>
          <div class="muted">
            {{ customer.phone_number }}
          </div>
        </div>

        <div class="muted">
          Машины: {{ vehiclesCount }}
        </div>
      </div>

      <div v-if="vehicleModelsText" class="muted" style="margin-top: 8px;">
        {{ vehicleModelsText }}
      </div>
    </div>
  </NuxtLink>
</template>

<script setup>
const props = defineProps({
  customer: { type: Object, required: true },
  vehicleTypeById: { type: Object, required: true }, // { [id]: "Brand Model" }
})

const vehicles = computed(() => props.customer.vehicles || [])
const vehiclesCount = computed(() => vehicles.value.length)

const vehicleModelsText = computed(() => {
  if (!vehicles.value.length) return ""

  // соберём названия моделей по vehicle_type_id
  const names = vehicles.value
    .map(v => props.vehicleTypeById?.[v.vehicle_type_id] || "Unknown model")

  // уберём повторы (если 2 одинаковых авто)
  const unique = Array.from(new Set(names))

  // покажем максимум 3, чтобы не было простыни
  const shown = unique.slice(0, 3)
  const more = unique.length > 3 ? ` +${unique.length - 3}` : ""

  return `${shown.join(", ")}${more}`
})
</script>
