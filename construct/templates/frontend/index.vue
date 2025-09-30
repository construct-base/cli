<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { use{{.PluralName}} } from './composable'
import type { {{.ResourceName}} } from './types'

const { {{.LowerPluralName}}, loading, fetch{{.PluralName}}, create{{.ResourceName}}, update{{.ResourceName}}, delete{{.ResourceName}} } = use{{.PluralName}}()

onMounted(() => {
  fetch{{.PluralName}}()
})

const columns = [
  { key: 'id', label: 'ID' },
  {{range .Fields}}{ key: '{{.Name}}', label: '{{.Label}}' },
  {{end}}{ key: 'created_at', label: 'Created' },
  {
    key: 'actions',
    label: 'Actions'
  }
]

const showAddModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const selectedItem = ref<{{.ResourceName}} | null>(null)

const handleEdit = (item: {{.ResourceName}}) => {
  selectedItem.value = item
  showEditModal.value = true
}

const handleDelete = (item: {{.ResourceName}}) => {
  selectedItem.value = item
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  if (selectedItem.value) {
    try {
      await delete{{.ResourceName}}(selectedItem.value.id)
      showDeleteModal.value = false
      selectedItem.value = null
    } catch (error) {
      console.error('Failed to delete {{.LowerResourceName}}:', error)
    }
  }
}
</script>

<template>
  <UDashboardPanel>
    <template #header>
      <UDashboardNavbar title="{{.PluralName}}">
        <template #right>
          <UButton @click="showAddModal = true" icon="i-lucide-plus">
            Add {{.ResourceName}}
          </UButton>
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <UTable
        :rows="{{.LowerPluralName}}"
        :columns="columns"
        :loading="loading"
      >
        {{range .Fields}}{{if .IsBool}}<template #{{.Name}}-data="{ row }">
          <UBadge :color="row.{{.Name}} ? 'green' : 'gray'">
            {{`{{ row.`}}{{.Name}}{{` ? '`}}{{.TrueLabel}}{{`' : '`}}{{.FalseLabel}}{{`' }}`}}
          </UBadge>
        </template>

        {{end}}{{end}}<template #actions-data="{ row }">
          <div class="flex gap-2">
            <UButton
              size="xs"
              color="primary"
              variant="ghost"
              icon="i-lucide-pencil"
              @click="handleEdit(row)"
            />
            <UButton
              size="xs"
              color="red"
              variant="ghost"
              icon="i-lucide-trash"
              @click="handleDelete(row)"
            />
          </div>
        </template>
      </UTable>
    </template>
  </UDashboardPanel>

  <!-- Add Modal -->
  <UModal v-model="showAddModal" title="Add {{.ResourceName}}">
    <div class="p-4">
      <p class="text-sm text-gray-500">Form implementation needed</p>
    </div>
  </UModal>

  <!-- Edit Modal -->
  <UModal v-model="showEditModal" title="Edit {{.ResourceName}}">
    <div class="p-4">
      <p class="text-sm text-gray-500">Form implementation needed</p>
    </div>
  </UModal>

  <!-- Delete Modal -->
  <UModal v-model="showDeleteModal" title="Delete {{.ResourceName}}">
    <div class="p-4 space-y-4">
      <p class="text-sm">
        Are you sure you want to delete this {{.LowerResourceName}}? This action cannot be undone.
      </p>
      <div class="flex justify-end gap-2">
        <UButton
          color="gray"
          variant="ghost"
          @click="showDeleteModal = false"
        >
          Cancel
        </UButton>
        <UButton
          color="red"
          @click="confirmDelete"
          :loading="loading"
        >
          Delete
        </UButton>
      </div>
    </div>
  </UModal>
</template>