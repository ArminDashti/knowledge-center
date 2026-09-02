<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { listImports, type ImportItem } from '@/api/client'
import { Button } from '@/components/ui/button'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const router = useRouter()
const items = ref<ImportItem[]>([])
const loading = ref(false)
const error = ref('')
let timer: number | undefined

async function load() {
  loading.value = true
  error.value = ''
  try {
    items.value = await listImports()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load imports'
  } finally {
    loading.value = false
  }
}

function statusClass(status: string) {
  switch (status) {
    case 'success':
      return 'bg-emerald-100 text-emerald-800'
    case 'failed':
      return 'bg-red-100 text-red-800'
    case 'running':
      return 'bg-amber-100 text-amber-800'
    default:
      return 'bg-slate-100 text-slate-700'
  }
}

function openItem(id: string) {
  void router.push(`/list/${id}`)
}

onMounted(() => {
  void load()
  timer = window.setInterval(() => {
    void load()
  }, 4000)
})

onUnmounted(() => {
  if (timer) window.clearInterval(timer)
})
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-3">
      <div>
        <h2 class="text-lg font-semibold">Imported websites</h2>
        <p class="text-sm text-muted-foreground">Click a row to browse scraped pages and content.</p>
      </div>
      <Button variant="outline" :disabled="loading" @click="load">Refresh</Button>
    </div>

    <p v-if="error" class="text-sm text-destructive">{{ error }}</p>

    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>URL</TableHead>
          <TableHead>Host</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Created</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow
          v-for="item in items"
          :key="item.id"
          class="cursor-pointer hover:bg-accent/60"
          @click="openItem(item.id)"
        >
          <TableCell class="max-w-xs truncate font-medium">{{ item.url }}</TableCell>
          <TableCell>{{ item.host }}</TableCell>
          <TableCell>
            <span class="rounded-md px-2 py-1 text-xs font-medium capitalize" :class="statusClass(item.status)">
              {{ item.status }}
            </span>
          </TableCell>
          <TableCell class="whitespace-nowrap text-muted-foreground">
            {{ new Date(item.created_at).toLocaleString() }}
          </TableCell>
        </TableRow>
        <TableRow v-if="!loading && items.length === 0">
          <TableCell colspan="4" class="text-center text-muted-foreground">No imports yet.</TableCell>
        </TableRow>
      </TableBody>
    </Table>
  </div>
</template>
