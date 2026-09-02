<script setup lang="ts">
import {
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogOverlay,
  AlertDialogPortal,
  AlertDialogRoot,
  AlertDialogTitle,
} from 'reka-ui'

const open = defineModel<boolean>('open', { default: false })

const emit = defineEmits<{ confirm: [] }>()
</script>

<template>
  <AlertDialogRoot v-model:open="open">
    <AlertDialogPortal>
      <AlertDialogOverlay class="fixed inset-0 z-50 bg-black/40" />
      <AlertDialogContent
        class="fixed left-1/2 top-1/2 z-50 grid w-full max-w-md -translate-x-1/2 -translate-y-1/2 gap-4 rounded-lg border border-border bg-card p-6 shadow-lg"
      >
        <AlertDialogTitle class="text-lg font-semibold">
          <slot name="title">Confirm</slot>
        </AlertDialogTitle>
        <AlertDialogDescription class="text-sm text-muted-foreground">
          <slot name="description" />
        </AlertDialogDescription>
        <div class="flex justify-end gap-2">
          <AlertDialogCancel
            class="inline-flex h-9 items-center rounded-md border border-border bg-card px-4 text-sm hover:bg-accent"
          >
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction
            class="inline-flex h-9 items-center rounded-md bg-destructive px-4 text-sm text-destructive-foreground hover:opacity-90"
            @click="emit('confirm')"
          >
            Delete
          </AlertDialogAction>
        </div>
      </AlertDialogContent>
    </AlertDialogPortal>
  </AlertDialogRoot>
</template>
