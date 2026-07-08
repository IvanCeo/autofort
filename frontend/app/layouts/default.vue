<template>
  <div class="app-layout">
    <!-- sidebar -->
    <aside class="sidebar">
      <div class="sidebar-title">Autofort</div>

      <nav class="sidebar-nav">
        <NuxtLink to="/" class="sidebar-link">Доманшняя</NuxtLink>
        <NuxtLink to="/customers" class="sidebar-link">Клиенты</NuxtLink>
      </nav>
    </aside>

    <!-- main + footer -->
    <div class="app-right">
      <main class="app-main">
        <div class="container">
          <slot />
        </div>
      </main>

      <!-- footer -->
      <footer class="footer">
        <span class="muted">
          © 2026 · made by
          <a href="https://t.me/iva0520" target="_blank" rel="noopener">
            @iva0520
          </a>
        </span>
      </footer>
    </div>
  </div>
</template>
