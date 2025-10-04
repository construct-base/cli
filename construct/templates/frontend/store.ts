import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { use{{.PluralName}} } from '../composables/use{{.PluralName}}'
import type { {{.ResourceName}}, {{.ResourceName}}CreateRequest, {{.ResourceName}}UpdateRequest, QueryParams } from '../types'

export const use{{.PluralName}}Store = defineStore('{{.LowerPluralName}}', () => {
  // Get the composable with API operations
  const {{.LowerPluralName}}Api = use{{.PluralName}}()

  // State
  const {{.LowerPluralName}} = ref<{{.ResourceName}}[]>([])
  const selected{{.ResourceName}} = ref<{{.ResourceName}} | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const pagination = ref({
    total: 0,
    page: 1,
    page_size: 10,
    total_pages: 1
  })

  // Search
  const searchQuery = ref('')

  // Getters
  const total{{.PluralName}} = computed(() => pagination.value.total)
  const has{{.PluralName}} = computed(() => {{.LowerPluralName}}.value.length > 0)
  const isLoading = computed(() => loading.value)

  // Filtered items based on search
  const filtered{{.PluralName}} = computed(() => {
    let filtered = {{.LowerPluralName}}.value

    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      filtered = filtered.filter((item: {{.ResourceName}}) => {
        // Search across all string fields
        {{range .Fields}}{{if eq .TypeScriptType "string"}}if (item.{{.Name}} && item.{{.Name}}.toLowerCase().includes(query)) return true
        {{end}}{{end}}return false
      })
    }

    return filtered
  })

  // Actions - using the composable internally
  const fetch{{.PluralName}} = async (params?: QueryParams): Promise<void> => {
    loading.value = true
    error.value = null

    try {
      const result = await {{.LowerPluralName}}Api.fetch{{.PluralName}}(params)
      {{.LowerPluralName}}.value = result.{{.LowerPluralName}}
      pagination.value = result.pagination
    } catch (err: unknown) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch {{.LowerPluralName}}'
    } finally {
      loading.value = false
    }
  }

  const fetch{{.ResourceName}} = async (id: number): Promise<{{.ResourceName}} | null> => {
    loading.value = true
    error.value = null

    try {
      const item = await {{.LowerPluralName}}Api.fetch{{.ResourceName}}(id)
      selected{{.ResourceName}}.value = item
      return item
    } catch (err: unknown) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch {{.LowerResourceName}}'
      return null
    } finally {
      loading.value = false
    }
  }

  const create{{.ResourceName}} = async (data: {{.ResourceName}}CreateRequest): Promise<{{.ResourceName}} | null> => {
    loading.value = true
    error.value = null

    try {
      const newItem = await {{.LowerPluralName}}Api.create{{.ResourceName}}(data)
      {{.LowerPluralName}}.value.push(newItem)
      pagination.value.total += 1
      return newItem
    } catch (err: unknown) {
      error.value = err instanceof Error ? err.message : 'Failed to create {{.LowerResourceName}}'
      return null
    } finally {
      loading.value = false
    }
  }

  const update{{.ResourceName}} = async (id: number, data: {{.ResourceName}}UpdateRequest): Promise<{{.ResourceName}} | null> => {
    loading.value = true
    error.value = null

    try {
      const updatedItem = await {{.LowerPluralName}}Api.update{{.ResourceName}}(id, data)
      const index = {{.LowerPluralName}}.value.findIndex(item => item.id === id)
      if (index !== -1) {
        {{.LowerPluralName}}.value[index] = updatedItem
      }
      if (selected{{.ResourceName}}.value?.id === id) {
        selected{{.ResourceName}}.value = updatedItem
      }
      return updatedItem
    } catch (err: unknown) {
      error.value = err instanceof Error ? err.message : 'Failed to update {{.LowerResourceName}}'
      return null
    } finally {
      loading.value = false
    }
  }

  const delete{{.ResourceName}} = async (id: number): Promise<boolean> => {
    loading.value = true
    error.value = null

    try {
      await {{.LowerPluralName}}Api.delete{{.ResourceName}}(id)
      {{.LowerPluralName}}.value = {{.LowerPluralName}}.value.filter(item => item.id !== id)
      if (selected{{.ResourceName}}.value?.id === id) {
        selected{{.ResourceName}}.value = null
      }
      pagination.value.total -= 1
      return true
    } catch (err: unknown) {
      error.value = err instanceof Error ? err.message : 'Failed to delete {{.LowerResourceName}}'
      return false
    } finally {
      loading.value = false
    }
  }

  // Helper actions
  const setSearchQuery = (query: string) => {
    searchQuery.value = query
  }

  const setPage = async (page: number): Promise<void> => {
    await fetch{{.PluralName}}({ page, page_size: pagination.value.page_size })
  }

  const setPerPage = async (perPage: number): Promise<void> => {
    await fetch{{.PluralName}}({ page: 1, page_size: perPage })
  }

  const clearError = () => {
    error.value = null
  }

  const clearSelected{{.ResourceName}} = () => {
    selected{{.ResourceName}}.value = null
  }

  const clearFilters = () => {
    searchQuery.value = ''
  }

  return {
    // State
    {{.LowerPluralName}},
    selected{{.ResourceName}},
    loading,
    error,
    pagination,
    searchQuery,

    // Getters
    total{{.PluralName}},
    has{{.PluralName}},
    isLoading,
    filtered{{.PluralName}},

    // Actions
    fetch{{.PluralName}},
    fetch{{.ResourceName}},
    create{{.ResourceName}},
    update{{.ResourceName}},
    delete{{.ResourceName}},
    setSearchQuery,
    setPage,
    setPerPage,
    clearError,
    clearSelected{{.ResourceName}},
    clearFilters
  }
})
