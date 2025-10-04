import { ref } from 'vue'
import { apiClient } from '~/core/api/client'
import type { {{.ResourceName}}, {{.ResourceName}}CreateRequest, {{.ResourceName}}UpdateRequest } from '../types'

export function use{{.PluralName}}() {
  const {{.LowerPluralName}} = ref<{{.ResourceName}}[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetch{{.PluralName}} = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await apiClient.get('/{{.LowerPluralName}}')
      {{.LowerPluralName}}.value = response.data
    } catch (e: any) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  const create{{.ResourceName}} = async (data: {{.ResourceName}}CreateRequest) => {
    loading.value = true
    error.value = null
    try {
      await apiClient.post('/{{.LowerPluralName}}', data)
      await fetch{{.PluralName}}()
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  const update{{.ResourceName}} = async (id: number, data: {{.ResourceName}}UpdateRequest) => {
    loading.value = true
    error.value = null
    try {
      await apiClient.put('/{{.LowerPluralName}}/' + id, data)
      await fetch{{.PluralName}}()
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  const delete{{.ResourceName}} = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      await apiClient.delete('/{{.LowerPluralName}}/' + id)
      await fetch{{.PluralName}}()
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    {{.LowerPluralName}},
    loading,
    error,
    fetch{{.PluralName}},
    create{{.ResourceName}},
    update{{.ResourceName}},
    delete{{.ResourceName}}
  }
}