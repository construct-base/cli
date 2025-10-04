<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import * as z from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { use{{.PluralName}}Store } from '../stores/{{.LowerPluralName}}'
import type { {{.ResourceName}} } from '../types'

const props = defineProps<{
  {{.LowerResourceName}}?: {{.ResourceName}} | null
}>()

const emit = defineEmits<{
  success: []
}>()

const store = use{{.PluralName}}Store()
const toast = useToast()

const open = ref(false)
const isEditing = computed(() => !!props.{{.LowerResourceName}})

// Validation schema
const schema = z.object({
{{range .Fields}}{{if eq .Type "string"}}  {{.Name}}: z.string().min(1, '{{.Label}} is required'),
{{else if eq .Type "text"}}  {{.Name}}: z.string().optional(),
{{else if eq .Type "int"}}  {{.Name}}: z.number().int(),
{{else if eq .Type "uint"}}  {{.Name}}: z.number().int().positive(),
{{else if eq .Type "float"}}  {{.Name}}: z.number(),
{{else if eq .Type "bool"}}  {{.Name}}: z.boolean().optional(),
{{else if eq .Type "email"}}  {{.Name}}: z.string().email('Invalid email address'),
{{else if eq .Type "url"}}  {{.Name}}: z.string().url('Invalid URL').optional(),
{{else}}  {{.Name}}: z.string().optional(),
{{end}}{{end}}})

type Schema = z.output<typeof schema>

const state = reactive<Partial<Schema>>({
{{range .Fields}}  {{.Name}}: {{.ZeroValue}},
{{end}}})

// Watch for prop changes to populate form when editing
watch(() => props.{{.LowerResourceName}}, (item) => {
  if (item) {
{{range .Fields}}    state.{{.Name}} = item.{{.Name}}
{{end}}    open.value = true
  }
}, { immediate: true })

// Watch open state to reset form when closed
watch(open, (isOpen) => {
  if (!isOpen && !props.{{.LowerResourceName}}) {
    resetForm()
  }
})

function resetForm() {
{{range .Fields}}  state.{{.Name}} = {{.ZeroValue}}
{{end}}}

async function onSubmit(event: FormSubmitEvent<Schema>) {
  try {
    if (isEditing.value && props.{{.LowerResourceName}}) {
      // Update existing item
      await store.update{{.ResourceName}}(props.{{.LowerResourceName}}.id, {
{{range .Fields}}        {{.Name}}: event.data.{{.Name}},
{{end}}      })

      toast.add({
        title: 'Success',
        description: `{{.ResourceName}} updated successfully`,
        color: 'success',
        icon: 'i-lucide-check-circle'
      })
    } else {
      // Create new item
      await store.create{{.ResourceName}}({
{{range .Fields}}        {{.Name}}: event.data.{{.Name}}!,
{{end}}      })

      toast.add({
        title: 'Success',
        description: `New {{.LowerResourceName}} added successfully`,
        color: 'success',
        icon: 'i-lucide-check-circle'
      })
    }

    open.value = false
    emit('success')
    resetForm()
  } catch (error) {
    toast.add({
      title: 'Error',
      description: error instanceof Error ? error.message : 'Failed to save {{.LowerResourceName}}',
      color: 'error',
      icon: 'i-lucide-alert-circle'
    })
  }
}
</script>

<template>
  <UModal
    v-model:open="open"
    :title="isEditing ? 'Edit {{.ResourceName}}' : 'New {{.ResourceName}}'"
    :description="isEditing ? 'Update {{.LowerResourceName}} information' : 'Add a new {{.LowerResourceName}} to the system'"
  >
    <UButton
      v-if="!{{.LowerResourceName}}"
      label="New {{.LowerResourceName}}"
      icon="i-lucide-plus"
    />

    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="onSubmit"
      >
{{range .Fields}}{{if eq .Type "text"}}        <UFormField label="{{.Label}}" name="{{.Name}}">
          <UTextarea v-model="state.{{.Name}}" rows="4" class="w-full" />
        </UFormField>
{{else if eq .Type "bool"}}        <UFormField label="{{.Label}}" name="{{.Name}}">
          <UCheckbox v-model="state.{{.Name}}" />
        </UFormField>
{{else if eq .Type "email"}}        <UFormField label="{{.Label}}" placeholder="email@example.com" name="{{.Name}}" required>
          <UInput v-model="state.{{.Name}}" type="email" class="w-full" />
        </UFormField>
{{else if eq .Type "url"}}        <UFormField label="{{.Label}}" placeholder="https://example.com" name="{{.Name}}">
          <UInput v-model="state.{{.Name}}" type="url" class="w-full" />
        </UFormField>
{{else if or (eq .Type "int") (eq .Type "uint") (eq .Type "float")}}        <UFormField label="{{.Label}}" name="{{.Name}}" required>
          <UInput v-model="state.{{.Name}}" type="number" class="w-full" />
        </UFormField>
{{else}}        <UFormField label="{{.Label}}" name="{{.Name}}" required>
          <UInput v-model="state.{{.Name}}" class="w-full" />
        </UFormField>
{{end}}{{end}}
        <div class="flex justify-end gap-2">
          <UButton
            label="Cancel"
            color="neutral"
            variant="subtle"
            @click="open = false"
          />
          <UButton
            :label="isEditing ? 'Update' : 'Create'"
            color="primary"
            variant="solid"
            type="submit"
            loading-auto
          />
        </div>
      </UForm>
    </template>
  </UModal>
</template>
