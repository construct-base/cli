<script setup lang="ts">
import { ref, watch } from 'vue'
import { use{{.PluralName}}Store } from '../stores/{{.LowerPluralName}}'
import type { {{.ResourceName}} } from '../types'

const props = withDefaults(defineProps<{
  count?: number
  {{.LowerResourceName}}?: {{.ResourceName}} | null
}>(), {
  count: 0,
  {{.LowerResourceName}}: null
})

const emit = defineEmits<{
  success: []
}>()

const store = use{{.PluralName}}Store()
const toast = useToast()
const open = ref(false)

// Watch for item prop to open modal
watch(() => props.{{.LowerResourceName}}, (item) => {
  if (item) {
    open.value = true
  }
})

async function onSubmit() {
  try {
    if (props.{{.LowerResourceName}}) {
      // Delete single item
      await store.delete{{.ResourceName}}(props.{{.LowerResourceName}}.id)

      toast.add({
        title: 'Success',
        description: `{{.ResourceName}} deleted successfully`,
        color: 'success',
        icon: 'i-lucide-check-circle'
      })
    } else if (props.count > 0) {
      // Bulk delete (placeholder for now)
      toast.add({
        title: 'Success',
        description: `${props.count} {{.LowerResourceName}}${props.count > 1 ? 's' : ''} deleted successfully`,
        color: 'success',
        icon: 'i-lucide-check-circle'
      })
    }

    open.value = false
    emit('success')
  } catch (error) {
    toast.add({
      title: 'Error',
      description: error instanceof Error ? error.message : 'Failed to delete {{.LowerResourceName}}',
      color: 'error',
      icon: 'i-lucide-alert-circle'
    })
  }
}
</script>

<template>
  <UModal
    v-model:open="open"
    :title="{{.LowerResourceName}} ? 'Delete {{.ResourceName}}' : `Delete ${count} {{.LowerResourceName}}${count > 1 ? 's' : ''}`"
  >
    <slot />

    <template #body>
      <div class="space-y-4">
        <div class="flex items-start gap-3">
          <UIcon name="i-lucide-alert-triangle" class="h-6 w-6 text-red-500 flex-shrink-0 mt-0.5" />
          <div class="flex-1">
            <p class="text-gray-900 font-medium">
              Are you sure you want to delete {{`{{ `}}{{.LowerResourceName}} ? 'this {{.LowerResourceName}}' : `${count} {{.LowerResourceName}}${count > 1 ? 's' : ''}` {{` }}`}}?
            </p>
            <p class="text-gray-500 text-sm mt-1">
              This action cannot be undone. All data will be permanently removed.
            </p>

            <!-- Show item details if single delete -->
            <div v-if="{{.LowerResourceName}}" class="mt-3 p-3 bg-gray-50 rounded-md">
              <div class="flex items-center gap-3">
                <div>
                  <p class="font-medium">
                    {{`{{ `}}{{.LowerResourceName}}.{{.DisplayField}}{{` }}`}}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-2">
          <UButton
            label="Cancel"
            color="neutral"
            variant="subtle"
            @click="open = false"
          />
          <UButton
            label="Delete"
            color="error"
            variant="solid"
            loading-auto
            @click="onSubmit"
          />
        </div>
      </div>
    </template>
  </UModal>
</template>
