<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import {
  createProfile,
  deleteProfile,
  listProfiles,
  updateProfile,
  type ProfileInput,
  type ScrapeProfile,
} from '@/api/client'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

const profiles = ref<ScrapeProfile[]>([])
const loading = ref(false)
const error = ref('')
const editingId = ref<string | null>(null)

const form = reactive<ProfileInput>({
  host: '',
  title_selector: 'h1',
  content_selector: 'article, main, body',
  link_selector: '',
  exclude_selector: '',
})

function resetForm() {
  editingId.value = null
  form.host = ''
  form.title_selector = 'h1'
  form.content_selector = 'article, main, body'
  form.link_selector = ''
  form.exclude_selector = ''
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    profiles.value = await listProfiles()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load profiles'
  } finally {
    loading.value = false
  }
}

function startEdit(profile: ScrapeProfile) {
  editingId.value = profile.id
  form.host = profile.host
  form.title_selector = profile.title_selector
  form.content_selector = profile.content_selector
  form.link_selector = profile.link_selector
  form.exclude_selector = profile.exclude_selector
}

async function onSubmit() {
  error.value = ''
  try {
    if (editingId.value) {
      await updateProfile(editingId.value, { ...form })
    } else {
      await createProfile({ ...form })
    }
    resetForm()
    await load()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Save failed'
  }
}

async function onDelete(id: string) {
  error.value = ''
  try {
    await deleteProfile(id)
    if (editingId.value === id) resetForm()
    await load()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Delete failed'
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <div class="space-y-6">
    <Card>
      <CardHeader>
        <CardTitle>{{ editingId ? 'Edit scrape profile' : 'Add scrape profile' }}</CardTitle>
        <CardDescription>
          Per-domain CSS selectors. Host is matched from the import URL (e.g. <code>example.com</code>).
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form class="grid gap-4 md:grid-cols-2" @submit.prevent="onSubmit">
          <div class="space-y-2 md:col-span-2">
            <Label for="host">Host</Label>
            <Input id="host" v-model="form.host" placeholder="example.com" required />
          </div>
          <div class="space-y-2">
            <Label for="title">Title selector</Label>
            <Input id="title" v-model="form.title_selector" placeholder="h1" />
          </div>
          <div class="space-y-2">
            <Label for="content">Content selector</Label>
            <Input id="content" v-model="form.content_selector" placeholder="article, main, body" />
          </div>
          <div class="space-y-2">
            <Label for="links">Link selector (optional crawl)</Label>
            <Input id="links" v-model="form.link_selector" placeholder="a.related" />
          </div>
          <div class="space-y-2">
            <Label for="exclude">Exclude selector</Label>
            <Input id="exclude" v-model="form.exclude_selector" placeholder="nav, footer, .ads" />
          </div>
          <div class="flex gap-2 md:col-span-2">
            <Button type="submit">{{ editingId ? 'Update' : 'Create' }}</Button>
            <Button v-if="editingId" type="button" variant="outline" @click="resetForm">Cancel</Button>
          </div>
        </form>
        <p v-if="error" class="mt-3 text-sm text-destructive">{{ error }}</p>
      </CardContent>
    </Card>

    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold">Profiles</h2>
        <Button variant="outline" :disabled="loading" @click="load">Refresh</Button>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Host</TableHead>
            <TableHead>Title</TableHead>
            <TableHead>Content</TableHead>
            <TableHead></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="profile in profiles" :key="profile.id">
            <TableCell class="font-medium">{{ profile.host }}</TableCell>
            <TableCell class="max-w-[10rem] truncate">{{ profile.title_selector }}</TableCell>
            <TableCell class="max-w-[12rem] truncate">{{ profile.content_selector }}</TableCell>
            <TableCell class="space-x-2 whitespace-nowrap text-right">
              <Button size="sm" variant="outline" @click="startEdit(profile)">Edit</Button>
              <Button size="sm" variant="destructive" @click="onDelete(profile.id)">Delete</Button>
            </TableCell>
          </TableRow>
          <TableRow v-if="!loading && profiles.length === 0">
            <TableCell colspan="4" class="text-center text-muted-foreground">
              No profiles yet — imports will use fallback selectors.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
