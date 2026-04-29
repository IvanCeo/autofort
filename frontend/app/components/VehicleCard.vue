<template>
  <div class="card">
    <!-- VIEW MODE -->
    <div v-if="!isEditing">
      <div class="row">
        <div>
          <div class="card-title">
            <!-- Brand + Model -->
            {{ vehicleTypeTitle }}
          </div>
          <div class="muted">
            Plate: {{ vehicle.gov_number || "—" }}
          </div>
        </div>

        <div class="card-actions">
          <button class="btn" @click="startEdit">Edit</button>
          <button class="btn" @click="print" :disabled="printing">
            {{ printing ? "Printing..." : "Print" }}
          </button>
        </div>
      </div>

      <div class="kv" style="margin-top:8px;">
        <div class="key">VIN</div>
        <div>{{ vehicle.vin || "-" }}</div>

        <div class="key">Mileage</div>
        <div>{{ vehicle.mileage ?? "-" }}</div>
      </div>

      <div v-if="actionError" class="muted" style="margin-top:8px;">
        Error: {{ actionError }}
      </div>
    </div>

    <!-- EDIT MODE -->
    <div v-else>
      <div class="row">
        <div class="card-title">{{ vehicleTypeTitle }}</div>
      </div>

      <div class="kv" style="margin-top:8px;">
        <div class="key">Gov number</div>
        <input class="input" v-model="form.gov_number" />

        <div class="key">VIN</div>
        <input class="input" v-model="form.vin" />

        <div class="key">Mileage</div>
        <input class="input" v-model="form.mileage" />
      </div>

      <div class="card-actions" style="margin-top:10px;">
        <button class="btn" @click="save" :disabled="saving">
          {{ saving ? "Saving..." : "Save" }}
        </button>
        <button class="btn" @click="cancel" :disabled="saving">
          Cancel
        </button>
      </div>

      <div v-if="actionError" class="muted" style="margin-top:8px;">
        Error: {{ actionError }}
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  vehicle: { type: Object, required: true },
  apiBase: { type: String, required: true }, // "http://localhost:8080/api"
  vehicleTypeById: { type: Object, required: true }, // { [id]: "Brand Model" }
})

const emit = defineEmits(["updated"])

const isEditing = ref(false)
const saving = ref(false)
const printing = ref(false)
const actionError = ref("")

const form = ref({
  gov_number: "",
  vin: "",
  mileage: "",
})

const vehicleTypeTitle = computed(() => {
  const id = props.vehicle.vehicle_type_id
  return props.vehicleTypeById?.[id] || "Unknown vehicle type"
})

function startEdit() {
  actionError.value = ""
  form.value.gov_number = props.vehicle.gov_number || ""
  form.value.vin = props.vehicle.vin || ""
  form.value.mileage = props.vehicle.mileage ?? ""
  isEditing.value = true
}

function cancel() {
  isEditing.value = false
  actionError.value = ""
}

async function save() {
  saving.value = true
  actionError.value = ""

  try {
    await $fetch(`${props.apiBase}/vehicles/${props.vehicle.id}`, {
      method: "PATCH",
      body: {
        gov_number: form.value.gov_number,
        vin: form.value.vin,
        mileage: form.value.mileage === "" ? null : Number(form.value.mileage),
      },
    })

    isEditing.value = false
    emit("updated")
  } catch (e) {
    actionError.value = e?.data ? JSON.stringify(e.data) : (e?.message || String(e))
  } finally {
    saving.value = false
  }
}

async function print() {
  printing.value = true
  actionError.value = ""

  try {
    const blob = await $fetch(
      `${props.apiBase}/vehicles/${props.vehicle.id}/work-order`,
      { method: "GET", responseType: "blob" }
    )

    const url = URL.createObjectURL(blob)
    window.open(url, "_blank")
    setTimeout(() => URL.revokeObjectURL(url), 60_000)
  } catch (e) {
    actionError.value = e?.message || String(e)
  } finally {
    printing.value = false
  }
}
</script>
