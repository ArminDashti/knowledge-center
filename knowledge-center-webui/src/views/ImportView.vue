<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { createImport } from '@/api/client'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

const router = useRouter()
const url = ref('')
const loading = ref(false)
const error = ref('')
const successId = ref('')

async function onSubmit() {
  error.value = ''
  successId.value = ''
  loading.value = true
  try {
    const item = await createImport(url.value.trim())
    successId.value = item.id
    url.value = ''
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Import failed'
  } finally {
    loading.value = false
  }
}

function goToList() {
  if (successId.value) {
    void router.push(`/list/${successId.value}`)
  } else {
    void router.push('/list')
  }
}
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Import a website</CardTitle>
      <CardDescription>
        Paste a URL. The API scrapes it using the host profile (or built-in fallbacks) and stores the content.
      </CardDescription>
    </CardHeader>
    <CardContent>
      <form class="flex flex-col gap-4" @submit.prevent="onSubmit">
        <div class="space-y-2">
          <Label for="url">Website URL</Label>
          <Input id="url" v-model="url" type="url" placeholder="https://example.com/article" :disabled="loading" />
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <Button type="submit" :disabled="loading || !url.trim()">
            {{ loading ? 'Submitting…' : 'Start import' }}
          </Button>
          <Button v-if="successId" type="button" variant="outline" @click="goToList">
            View import
          </Button>
        </div>
        <p v-if="successId" class="text-sm text-muted-foreground">
          Import queued. Status updates appear on the list page.
        </p>
        <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
      </form>
    </CardContent>
  </Card>
</template>
