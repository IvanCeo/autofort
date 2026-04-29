<template>
  <div>
    <div class="row">
      <h1 style="margin: 0;">Домашняя</h1>

      <NuxtLink to="/customers/new">
        <button class="btn">Добавить клиента</button>
      </NuxtLink>
    </div>

    <div class="card" style="margin-top: 12px;">
      <div class="card-title">Поиск клиента</div>

      <div style="margin-top: 8px;">
        <input class="input" v-model="q" placeholder="имя / фамилия / телефон" />
      </div>

      <div v-if="pending" class="muted" style="margin-top:8px;">
        Ищем...
      </div>

      <div v-else-if="error" class="muted" style="margin-top:8px;">
        Error: {{ error }}
      </div>

      <div v-else-if="q && searched && !items.length" class="muted" style="margin-top:8px;">
        Не найдено
      </div>

      <div v-else-if="!q" class="muted" style="margin-top:8px;">
        пиши чтобы искать...
      </div>
    </div>

    <div v-if="q && items.length" style="margin-top: 12px;">
      <div class="muted" style="margin-bottom: 8px;">
        Results: {{ items.length }}
      </div>

      <CustomerCard v-for="c in items" :key="c.id" :customer="c" />
    </div>
  </div>
</template>

<script setup>
const API_BASE = "/api"

const q = ref("")
const items = ref([])
const pending = ref(false)
const error = ref("")
const searched = ref(false)

let timer = null

watch(q, (val) => {
  // очистка если пусто
  if (!val) {
    if (timer) clearTimeout(timer)
    items.value = []
    pending.value = false
    error.value = ""
    searched.value = false
    return
  }

  // чтобы не искать по 1 символу (можно убрать если хочешь)
  if (val.length < 2) {
    if (timer) clearTimeout(timer)
    items.value = []
    pending.value = false
    error.value = ""
    searched.value = false
    return
  }

  // debounce 300ms
  if (timer) clearTimeout(timer)
  timer = setTimeout(async () => {
    pending.value = true
    error.value = ""
    searched.value = true

    try {
      const res = await $fetch(`${API_BASE}/customers`, {
        method: "GET",
        query: { q: val },
      })
      items.value = res?.items || []
    } catch (e) {
      error.value = e?.message || String(e)
      items.value = []
    } finally {
      pending.value = false
    }
  }, 300)
})
</script>
