<template>
  <div>
    <div class="row">
      <h1 style="margin: 0;">Клиенты</h1>
      <NuxtLink to="/customers/new">
        <button class="btn">Добавить клиента</button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="muted" style="margin-top: 12px;">Загрузка...</div>
    <div v-else-if="error" class="muted" style="margin-top: 12px;">
      Error: {{ error.message }}
    </div>

    <div v-else style="margin-top: 12px;">
      <div v-if="typesPending" class="muted" style="margin-bottom: 8px;">
        Загружаем машины...
      </div>
      <div v-else-if="typesError" class="muted" style="margin-bottom: 8px;">
        Vehicle types error: {{ typesError.message }}
      </div>

      <div v-if="customers.length">
        <CustomerCard
          v-for="c in customers"
          :key="c.id"
          :customer="c"
          :vehicleTypeById="vehicleTypeById"
        />
      </div>

      <div v-else class="muted">
        Нет клиентов
      </div>
    </div>
  </div>
</template>

<script setup async>
const API_BASE = "/api"

// customers
const { data, pending, error } = await useFetch(`${API_BASE}/customers`)
const customers = computed(() => data.value?.items || [])

// vehicle types (brand/model dictionary)
const { data: vehicleTypes, pending: typesPending, error: typesError } = await useFetch(
  `${API_BASE}/vehicle-types`
)

const vehicleTypeById = computed(() => {
  const map = {}
  for (const t of (vehicleTypes.value?.items || [])) {
    map[t.id] = `${t.brand} ${t.model}`
  }
  return map
})
</script>
