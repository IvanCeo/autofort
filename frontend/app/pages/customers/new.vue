<template>
  <div>
    <NuxtLink to="/">← Back</NuxtLink>

    <h1>Create customer</h1>

    <div v-if="error" class="muted" style="white-space: pre-wrap; margin-top: 8px;">
      Error: {{ error }}
    </div>

    <div class="card" style="margin-top: 12px;">
      <div class="kv">
        <div class="key">First name</div>
        <input class="input" v-model="form.first_name" />

        <div class="key">Last name</div>
        <input class="input" v-model="form.last_name" />

        <div class="key">Phone</div>
        <PhoneInput v-model="form.phone_number" />
      </div>

      <div class="card-actions" style="margin-top: 10px;">
        <button class="btn" @click="create" :disabled="pending || !canCreate">
          {{ pending ? "Creating..." : "Create" }}
        </button>

        <div v-if="!canCreate" class="muted" style="margin-top: 8px;">
          Enter full phone number
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const API_BASE = "/api"

const form = ref({
  first_name: "",
  last_name: "",
  phone_number: "", // "+7XXXXXXXXXX"
})

const pending = ref(false)
const error = ref("")

function digitsOnly(s) {
  return String(s || "").replace(/\D/g, "")
}

const canCreate = computed(() => {
  // номер полный, когда цифр ровно 11 и начинается с 7
  const d = digitsOnly(form.value.phone_number)
  const phoneOk = d.length === 11 && d.startsWith("7")
  return Boolean(form.value.first_name && form.value.last_name && phoneOk)
})

function extractCreatedId(res) {
  // swagger: {"created":"uuid"} (может прийти строкой)
  if (!res) return ""
  if (typeof res === "string") {
    try {
      return JSON.parse(res)?.created || ""
    } catch {
      return ""
    }
  }
  return res?.created || ""
}

async function create() {
  error.value = ""
  if (!canCreate.value) return

  pending.value = true
  try {
    const res = await $fetch(`${API_BASE}/customers`, {
      method: "POST",
      body: {
        first_name: form.value.first_name,
        last_name: form.value.last_name,
        phone_number: form.value.phone_number,
      },
    })

    const id = extractCreatedId(res)
    if (!id) {
      error.value = "Created, but cannot read id from response: " + JSON.stringify(res)
      return
    }

    await navigateTo(`/customers/${id}`)
  } catch (e) {
    error.value = e?.data ? JSON.stringify(e.data) : (e?.message || String(e))
  } finally {
    pending.value = false
  }
}
</script>
