<template>
  <div>
    <div class="row">
      <NuxtLink to="/customers">← Назад</NuxtLink>

      <!-- <button class="btn" @click="openAddVehicle" :disabled="pending">
        Добавить авто
      </button> -->
    </div>

    <h1>Клиент</h1>

    <div v-if="pending" class="muted" style="margin-top:12px;">Загрузка...</div>
    <div v-else-if="error" class="muted" style="margin-top:12px;">
      Error: {{ error.message }}
    </div>

    <div v-else>
      <!-- CUSTOMER CARD -->
      <div class="card">
        <div class="row" style="flex-direction: column;">
          <!-- <div> -->
            <div class="card-title">
              {{ customerData.first_name }} {{ customerData.last_name }}
            </div>
            <div class="muted">{{ customerData.phone_number }}</div>
          <!-- </div> -->

          <button class="btn" @click="startEdit" v-if="!isEditing">
            Изменить данные
          </button>
          <button class="btn" @click="openAddVehicle" :disabled="pending">
            Добавить авто
          </button> 
        </div>

        <!-- EDIT CUSTOMER -->
        <div v-if="isEditing" style="margin-top: 12px;">
          <div class="kv">
            <div class="key">Имя</div>
            <input class="input" v-model="form.first_name" />

            <div class="key">Фамилия</div>
            <input class="input" v-model="form.last_name" />

            <div class="key">Телефон</div>
            <PhoneInput v-model="form.phone_number" />
          </div>

          <div class="card-actions" style="margin-top:10px;">
            <button class="btn" @click="saveCustomer" :disabled="saving">
              {{ saving ? "Сохранение..." : "Сохранить" }}
            </button>
            <button class="btn" @click="cancelEdit" :disabled="saving">
              Отмена
            </button>
          </div>

          <div v-if="saveError" class="muted" style="margin-top:8px;">
            Error: {{ saveError }}
          </div>
        </div>
      </div>

      <!-- ADD VEHICLE FORM -->
      <div v-if="addingVehicle" class="card">
        <div class="row">
          <b>Добавить авто</b>
          <button class="btn" @click="closeAddVehicle" :disabled="vehicleSaving">
            Отмена
          </button>
        </div>

        <div class="kv" style="margin-top: 8px;">
          <div class="key">Гос. номер</div>
          <input class="input" v-model="vehicleForm.gov_number" />

          <div class="key">VIN</div>
          <input class="input" v-model="vehicleForm.vin" />

          <div class="key">Тип авто</div>
          <select class="input" v-model="vehicleForm.vehicle_type_id">
            <option value="">-- Выбрать --</option>
            <option
              v-for="t in (vehicleTypes?.items || [])"
              :key="t.id"
              :value="t.id"
            >
              {{ t.brand }} {{ t.model }}
            </option>
          </select>

          <div class="key">Пробег</div>
          <input class="input" v-model="vehicleForm.mileage" />
        </div>

        <div v-if="typesPending" class="muted" style="margin-top:8px;">
          Загрузка машин...
        </div>
        <div v-else-if="typesError" class="muted" style="margin-top:8px;">
          Vehicle types error: {{ typesError.message }}
        </div>

        <div class="card-actions" style="margin-top:10px;">
          <button class="btn" @click="createVehicle" :disabled="vehicleSaving">
            {{ vehicleSaving ? "Добавление..." : "Добавить" }}
          </button>
        </div>

        <div v-if="vehicleError" class="muted" style="margin-top:8px;">
          Error: {{ vehicleError }}
        </div>
      </div>

      <!-- VEHICLES -->
      <div v-if="vehicles.length">
        <h2>Машины</h2>

        <VehicleCard
          v-for="v in vehicles"
          :key="v.id"
          :vehicle="v"
          :apiBase="API_BASE"
          :vehicleTypeById="vehicleTypeById"
          @updated="refresh()"
        />
      </div>

      <div v-else class="muted">
        Пока нет машин
      </div>
    </div>
  </div>
</template>

<script setup async>
const route = useRoute()
const API_BASE = "/api"

/* vehicle types */
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

/* customer response: { customer: {...}, vehicles: [...] } */
const { data: customerResp, pending, error, refresh } = await useFetch(
  `${API_BASE}/customers/${route.params.id}`
)

const customerData = computed(() => customerResp.value?.customer || {})

/**
 * Берём машины:
 * 1) customerResp.vehicles (если есть)
 * 2) иначе customerResp.customer.vehicles
 */
const vehicles = computed(() => {
  const top = customerResp.value?.vehicles
  if (Array.isArray(top) && top.length) return top

  const nested = customerResp.value?.customer?.vehicles
  if (Array.isArray(nested) && nested.length) return nested

  return []
})

/* edit customer */
const isEditing = ref(false)
const saving = ref(false)
const saveError = ref("")

const form = ref({ first_name: "", last_name: "", phone_number: "" })

function startEdit() {
  saveError.value = ""
  form.value.first_name = customerData.value.first_name || ""
  form.value.last_name = customerData.value.last_name || ""
  form.value.phone_number = customerData.value.phone_number || ""
  isEditing.value = true
}

function cancelEdit() {
  isEditing.value = false
  saveError.value = ""
}

async function saveCustomer() {
  saving.value = true
  saveError.value = ""
  try {
    await $fetch(`${API_BASE}/customers/${route.params.id}`, {
      method: "PATCH",
      body: {
        first_name: form.value.first_name,
        last_name: form.value.last_name,
        phone_number: form.value.phone_number,
      },
    })
    await refresh()
    isEditing.value = false
  } catch (e) {
    saveError.value = e?.data ? JSON.stringify(e.data) : (e?.message || String(e))
  } finally {
    saving.value = false
  }
}

/* add vehicle */
const addingVehicle = ref(false)
const vehicleSaving = ref(false)
const vehicleError = ref("")

const vehicleForm = ref({
  gov_number: "",
  vin: "",
  vehicle_type_id: "",
  mileage: "",
})

function openAddVehicle() {
  vehicleError.value = ""
  vehicleForm.value.gov_number = ""
  vehicleForm.value.vin = ""
  vehicleForm.value.vehicle_type_id = ""
  vehicleForm.value.mileage = ""
  addingVehicle.value = true
}

function closeAddVehicle() {
  addingVehicle.value = false
  vehicleError.value = ""
}

async function createVehicle() {
  vehicleError.value = ""

  if (!vehicleForm.value.vehicle_type_id) {
    vehicleError.value = "Select vehicle type"
    return
  }

  vehicleSaving.value = true
  try {
    await $fetch(`${API_BASE}/customers/${route.params.id}/vehicles`, {
      method: "POST",
      body: {
        gov_number: vehicleForm.value.gov_number,
        vin: vehicleForm.value.vin,
        vehicle_type_id: vehicleForm.value.vehicle_type_id,
        mileage: vehicleForm.value.mileage === "" ? 0 : Number(vehicleForm.value.mileage),
      },
    })

    addingVehicle.value = false
    await refresh()
  } catch (e) {
    vehicleError.value = e?.data ? JSON.stringify(e.data) : (e?.message || String(e))
  } finally {
    vehicleSaving.value = false
  }
}
</script>
