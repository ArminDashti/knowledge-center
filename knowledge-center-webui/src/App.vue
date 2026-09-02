<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class="mx-auto flex min-h-screen max-w-5xl flex-col px-4 py-6">
    <header class="mb-8 flex flex-wrap items-center justify-between gap-4 border-b border-border pb-4">
      <div>
        <p class="text-xs uppercase tracking-[0.2em] text-muted-foreground">knowledge center</p>
        <h1 class="text-2xl font-semibold tracking-tight">Website knowledge base</h1>
      </div>
      <nav class="flex flex-wrap gap-2">
        <RouterLink
          to="/list"
          class="rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground"
          active-class="!bg-primary !text-primary-foreground"
        >
          List
        </RouterLink>
        <RouterLink
          to="/import"
          class="rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground"
          active-class="!bg-primary !text-primary-foreground"
        >
          Import
        </RouterLink>
        <RouterLink
          to="/scraping"
          class="rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground"
          active-class="!bg-primary !text-primary-foreground"
        >
          Scraping
        </RouterLink>
      </nav>
    </header>
    <main class="flex-1">
      <RouterView />
    </main>
  </div>
</template>
