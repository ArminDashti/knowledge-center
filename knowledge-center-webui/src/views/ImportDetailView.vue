<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getImport, getPage, type ImportDetail, type PageDetail } from '@/api/client'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

const route = useRoute()
const router = useRouter()
const detail = ref<ImportDetail | null>(null)
const selectedPage = ref<PageDetail | null>(null)
const loading = ref(false)
const pageLoading = ref(false)
const error = ref('')

const importId = computed(() => String(route.params.id))

async function loadDetail() {
  loading.value = true
  error.value = ''
  try {
    detail.value = await getImport(importId.value)
    if (detail.value.pages.length > 0) {
      await openPage(detail.value.pages[0].id)
    } else {
      selectedPage.value = null
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load import'
  } finally {
    loading.value = false
  }
}

async function openPage(pageId: string) {
  pageLoading.value = true
  try {
    selectedPage.value = await getPage(importId.value, pageId)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load page'
  } finally {
    pageLoading.value = false
  }
}

onMounted(() => {
  void loadDetail()
})

watch(importId, () => {
  void loadDetail()
})
</script>

<template>
  <div class="space-y-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <Button variant="outline" @click="router.push('/list')">Back to list</Button>
      <Button variant="ghost" :disabled="loading" @click="loadDetail">Refresh</Button>
    </div>

    <p v-if="error" class="text-sm text-destructive">{{ error }}</p>

    <Card v-if="detail">
      <CardHeader>
        <CardTitle class="break-all text-base">{{ detail.url }}</CardTitle>
        <CardDescription>
          Host {{ detail.host }} · status
          <span class="capitalize">{{ detail.status }}</span>
          <span v-if="detail.error_message"> · {{ detail.error_message }}</span>
        </CardDescription>
      </CardHeader>
      <CardContent class="grid gap-6 md:grid-cols-[240px_1fr]">
        <div class="space-y-2">
          <p class="text-xs font-semibold uppercase tracking-wide text-muted-foreground">Pages</p>
          <button
            v-for="page in detail.pages"
            :key="page.id"
            type="button"
            class="block w-full rounded-md border border-border px-3 py-2 text-left text-sm hover:bg-accent"
            :class="selectedPage?.id === page.id ? 'bg-accent' : ''"
            @click="openPage(page.id)"
          >
            <span class="line-clamp-2 font-medium">{{ page.title || page.url }}</span>
          </button>
          <p v-if="detail.pages.length === 0" class="text-sm text-muted-foreground">
            No pages yet (still scraping or failed).
          </p>
        </div>

        <div class="min-h-48 rounded-md border border-border bg-card p-4">
          <p v-if="pageLoading" class="text-sm text-muted-foreground">Loading content…</p>
          <template v-else-if="selectedPage">
            <h3 class="mb-2 text-lg font-semibold">{{ selectedPage.title || 'Untitled' }}</h3>
            <a class="mb-4 inline-block text-sm text-primary underline" :href="selectedPage.url" target="_blank" rel="noreferrer">
              {{ selectedPage.url }}
            </a>
            <pre class="whitespace-pre-wrap break-words text-sm leading-relaxed text-foreground/90">{{ selectedPage.content_text }}</pre>
          </template>
          <p v-else class="text-sm text-muted-foreground">Select a page to view content.</p>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
