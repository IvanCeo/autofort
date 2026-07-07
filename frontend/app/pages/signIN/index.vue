<template>
    <div class="login-container">
        <div class="card">
            <h2>Вход в систему</h2>
            <form @submit.prevent="login">
            <input type="email" v-model="form.login" placeholder="email" required />
            <input type="password" v-model="form.password" placeholder="пароль" required />
            <button type="submit" @click="signIN" :disabled="pending">
                {{  pending ? "Успех":"Войти" }}
            </button>
            </form>
        </div>
    </div>
</template>

<script setup>

definePageMeta({
    layout: 'auth'
})

const API_BASE = "/api"

const form = ref({
    login: "",
    password: "",
})

const pending = ref(false) // нужен чтобы передать состояние для стилей
const error = ref("")

function canSignIN() {
    const credentialsOK = form.value.login.length > 0 && form.value.password.length
    return Boolean(credentialsOK)
}

async function signIN() {
    error.value = ""

    if (!canSignIN()) return
    pending.value = true

    try {
        const res = await $fetch(`${API_BASE}/signIN`, {
            method: "POST",
            body: {
                login: form.value.login,
                password: form.value.password,
            },
        })

        // тут надо обработать res - в заголовках должен быть bearer

        await navigateTo(`/`)
    } catch (e) {
        error.value = e?.data ? JSON.stringify(e.data) : (e?.message || String(e))
    } finally {
        pending.value = false
    }
}
</script>

<style scoped>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    background-color: #f8f9fa;
}

.card {
    background: white;
    padding: 32px;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    width: 320px;
    text-align: center;
}

h2 {
    margin-bottom: 24px;
    font-size: 24px;
    color: #3c4043;
}

input {
    width: calc(100%-24px);
    padding: 12px;
    margin: 10px 0;
    border:  1px solid #dadce0;
    border-radius: 4px;
    font-size: 16px;
}

button {
    width: calc(100%-24px);
    padding: 12px;
    margin-top: 16px;
    background-color: #4285f4;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    cursor: pointer;
}

button:hover {
    background-color: #357ae8;
}
</style>