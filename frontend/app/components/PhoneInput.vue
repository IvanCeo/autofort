<template>
  <input
    class="input"
    :value="display"
    inputmode="tel"
    placeholder="+7 (___) ___-__-__"
    @input="onInput"
  />
</template>

<script setup>
const props = defineProps({
  modelValue: { type: String, default: "" }, // хранится как "+7XXXXXXXXXX" (частично/полностью)
})

const emit = defineEmits(["update:modelValue"])

function digitsOnly(s) {
  return String(s || "").replace(/\D/g, "")
}

function toPlus7Limited(raw) {
  let d = digitsOnly(raw)

  if (!d) return ""

  // 8XXXXXXXXXX -> 7XXXXXXXXXX
  if (d.startsWith("8")) d = "7" + d.slice(1)

  // если ввели без кода (например 999...) -> добавим 7
  if (!d.startsWith("7")) d = "7" + d

  // оставляем максимум 11 цифр (7 + 10 цифр номера)
  d = d.slice(0, 11)

  return "+" + d
}

function formatPretty(e164) {
  const d = digitsOnly(e164)
  if (!d) return ""

  const rest = d.startsWith("7") ? d.slice(1) : d.slice(0, 10)

  const a = rest.slice(0, 3)
  const b = rest.slice(3, 6)
  const c = rest.slice(6, 8)
  const e = rest.slice(8, 10)

  let out = "+7"
  if (a) out += ` (${a}`
  if (a.length === 3) out += ")"
  if (b) out += ` ${b}`
  if (c) out += `-${c}`
  if (e) out += `-${e}`
  return out
}

const display = computed(() => formatPretty(props.modelValue))

function onInput(ev) {
  const val = ev?.target?.value || ""
  emit("update:modelValue", toPlus7Limited(val))
}
</script>
