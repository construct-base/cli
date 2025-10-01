# Construct CRUD Generator - Detailed Plan

## Overview
Create a full-stack CRUD generator that scaffolds both Go backend and Vue frontend code from a single command.

## Command Syntax
```bash
construct generate Post title:string content:text published:bool
construct g Post title:string content:text published:bool  # alias
```

## Field Type Support

### Basic Types
- `string` → Go: `string`, Vue: `string`
- `text` → Go: `string` (long text), Vue: `string` (textarea)
- `int` → Go: `int`, Vue: `number`
- `uint` → Go: `uint`, Vue: `number`
- `float` / `float64` → Go: `float64`, Vue: `number`
- `bool` / `boolean` → Go: `bool`, Vue: `boolean`
- `date` → Go: `time.Time`, Vue: `Date`
- `datetime` / `timestamp` → Go: `time.Time`, Vue: `Date`

### Special Types
- `email` → `string` with email validation
- `url` → `string` with URL validation
- `image` → File upload with image validation
- `file` → Generic file upload

### Relationships (Future)
- `category_id:uint` → Auto-detects BelongsTo relationship
- `user_id:uint` → Auto-detects BelongsTo relationship

## Backend Generation (Go)

### File Structure
```
api/posts/
├── controller.go   # HTTP handlers (CRUD endpoints)
├── service.go      # Business logic
├── module.go       # Module registration
└── validator.go    # Input validation rules

api/models/
└── post.go         # GORM model
```

### Generated Files

#### 1. `api/models/post.go`
```go
package models

import (
    "base/core/types"
    "time"
)

type Post struct {
    types.Model
    Title     string    `json:"title" gorm:"type:varchar(255);not null"`
    Content   string    `json:"content" gorm:"type:text"`
    Published bool      `json:"published" gorm:"default:false"`
}
```

#### 2. `api/posts/service.go`
```go
package posts

import (
    "base/api/models"
    "gorm.io/gorm"
)

type Service struct {
    db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
    return &Service{db: db}
}

func (s *Service) GetAll() ([]models.Post, error)
func (s *Service) GetByID(id uint) (*models.Post, error)
func (s *Service) Create(data CreateRequest) (*models.Post, error)
func (s *Service) Update(id uint, data UpdateRequest) (*models.Post, error)
func (s *Service) Delete(id uint) error
```

#### 3. `api/posts/controller.go`
```go
package posts

import (
    "base/core/router"
    "strconv"
)

type Controller struct {
    service *Service
}

func NewController(service *Service) *Controller {
    return &Controller{service: service}
}

// @Summary List posts
// @Tags posts
// @Produce json
// @Success 200 {array} models.Post
// @Router /posts [get]
func (c *Controller) List(ctx *router.Context) error

// @Summary Get post
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post
// @Router /posts/{id} [get]
func (c *Controller) Get(ctx *router.Context) error

// @Summary Create post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body CreateRequest true "Post data"
// @Success 201 {object} models.Post
// @Router /posts [post]
func (c *Controller) Create(ctx *router.Context) error

// @Summary Update post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body UpdateRequest true "Post data"
// @Success 200 {object} models.Post
// @Router /posts/{id} [put]
func (c *Controller) Update(ctx *router.Context) error

// @Summary Delete post
// @Tags posts
// @Param id path int true "Post ID"
// @Success 204
// @Router /posts/{id} [delete]
func (c *Controller) Delete(ctx *router.Context) error
```

#### 4. `api/posts/validator.go`
```go
package posts

type CreateRequest struct {
    Title     string `json:"title" validate:"required,max=255"`
    Content   string `json:"content"`
    Published bool   `json:"published"`
}

type UpdateRequest struct {
    Title     *string `json:"title,omitempty" validate:"omitempty,max=255"`
    Content   *string `json:"content,omitempty"`
    Published *bool   `json:"published,omitempty"`
}
```

#### 5. `api/posts/module.go`
```go
package posts

import (
    "base/core/module"
    "base/core/router"
    "gorm.io/gorm"
)

type Module struct {
    controller *Controller
}

func Init(deps module.Dependencies) module.Module {
    service := NewService(deps.DB)
    controller := NewController(service)

    return &Module{controller: controller}
}

func (m *Module) Register(router *router.Router) {
    group := router.Group("/posts")

    group.GET("", m.controller.List)
    group.GET("/:id", m.controller.Get)
    group.POST("", m.controller.Create)
    group.PUT("/:id", m.controller.Update)
    group.DELETE("/:id", m.controller.Delete)
}

func (m *Module) Migrate(db *gorm.DB) error {
    return db.AutoMigrate(&models.Post{})
}
```

### Auto-Registration in `api/init.go`
Automatically inject:
```go
import "base/api/posts"

modules["posts"] = posts.Init(deps)
```

## Frontend Generation (Vue)

### File Structure
```
vue/structures/posts/
├── index.vue       # Main CRUD page with table + forms
├── composable.ts   # API integration & state
└── types.ts        # TypeScript types
```

### Generated Files

#### 1. `vue/structures/posts/types.ts`
```typescript
export interface Post {
  id: number
  title: string
  content: string
  published: boolean
  created_at: string
  updated_at: string
  deleted_at?: string | null
}

export interface PostCreateRequest {
  title: string
  content: string
  published: boolean
}

export interface PostUpdateRequest {
  title?: string
  content?: string
  published?: boolean
}

export interface PostListResponse {
  data: Post[]
  total: number
}
```

#### 2. `vue/structures/posts/composable.ts`
```typescript
import { ref } from 'vue'
import type { Post, PostCreateRequest, PostUpdateRequest } from './types'
import { useApi } from '@/core/composables/useApi'

export function usePosts() {
  const api = useApi()
  const posts = ref<Post[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchPosts = async () => {
    loading.value = true
    try {
      const response = await api.get<Post[]>('/posts')
      posts.value = response.data
    } catch (e) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  const getPost = async (id: number) => {
    return await api.get<Post>(\`/posts/\${id}\`)
  }

  const createPost = async (data: PostCreateRequest) => {
    const response = await api.post<Post>('/posts', data)
    await fetchPosts()
    return response.data
  }

  const updatePost = async (id: number, data: PostUpdateRequest) => {
    const response = await api.put<Post>(\`/posts/\${id}\`, data)
    await fetchPosts()
    return response.data
  }

  const deletePost = async (id: number) => {
    await api.delete(\`/posts/\${id}\`)
    await fetchPosts()
  }

  return {
    posts,
    loading,
    error,
    fetchPosts,
    getPost,
    createPost,
    updatePost,
    deletePost
  }
}
```

#### 3. `vue/structures/posts/index.vue`
```vue
<template>
  <div class="posts-page">
    <div class="header">
      <h1>Posts</h1>
      <UButton @click="showCreateModal = true">Create Post</UButton>
    </div>

    <!-- Table -->
    <UTable
      :rows="posts"
      :columns="columns"
      :loading="loading"
    >
      <template #actions-data="{ row }">
        <UButton size="xs" @click="editPost(row)">Edit</UButton>
        <UButton size="xs" color="red" @click="confirmDelete(row)">Delete</UButton>
      </template>
    </UTable>

    <!-- Create/Edit Modal -->
    <UModal v-model="showCreateModal">
      <UForm :state="form" @submit="handleSubmit">
        <UFormGroup label="Title" required>
          <UInput v-model="form.title" />
        </UFormGroup>

        <UFormGroup label="Content">
          <UTextarea v-model="form.content" />
        </UFormGroup>

        <UFormGroup label="Published">
          <UCheckbox v-model="form.published" />
        </UFormGroup>

        <UButton type="submit" :loading="loading">
          {{ isEditing ? 'Update' : 'Create' }}
        </UButton>
      </UForm>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { usePosts } from './composable'
import type { Post } from './types'

const { posts, loading, fetchPosts, createPost, updatePost, deletePost } = usePosts()

const showCreateModal = ref(false)
const isEditing = ref(false)
const form = ref({
  title: '',
  content: '',
  published: false
})

const columns = [
  { key: 'id', label: 'ID' },
  { key: 'title', label: 'Title' },
  { key: 'published', label: 'Published' },
  { key: 'created_at', label: 'Created' },
  { key: 'actions', label: 'Actions' }
]

onMounted(() => {
  fetchPosts()
})

const handleSubmit = async () => {
  if (isEditing.value) {
    await updatePost(currentId.value, form.value)
  } else {
    await createPost(form.value)
  }
  showCreateModal.value = false
  resetForm()
}
</script>
```

### Auto-Registration in `vue/core/main.ts`
Automatically inject route:
```typescript
{ path: 'posts', component: () => import('@/structures/posts/index.vue') }
```

## Implementation Steps

### Phase 1: Core Generator Infrastructure
1. Create generator command structure
2. Parse field definitions
3. Template engine setup
4. File system operations

### Phase 2: Backend Generation
1. Model generation with GORM tags
2. Service layer with CRUD methods
3. Controller with HTTP handlers
4. Validator with validation rules
5. Module with router registration
6. Auto-inject into `api/init.go`

### Phase 3: Frontend Generation
1. TypeScript types matching Go models
2. Composable with API methods
3. Vue component with table + forms
4. Auto-inject route into `vue/core/main.ts`

### Phase 4: Polish & Features
1. Support for relationships
2. File upload fields
3. Validation rules
4. Pagination support
5. Search/filter support
6. Soft delete support

## Templates Location
```
cli/construct/templates/
├── backend/
│   ├── model.tmpl
│   ├── service.tmpl
│   ├── controller.tmpl
│   ├── validator.tmpl
│   └── module.tmpl
└── frontend/
    ├── types.tmpl
    ├── composable.tmpl
    └── index.tmpl
```

## Success Criteria
- ✅ One command generates full-stack CRUD
- ✅ Backend compiles without errors
- ✅ Frontend has no TypeScript errors
- ✅ API endpoints work immediately
- ✅ UI displays and functions correctly
- ✅ Routes auto-registered
- ✅ Migrations run automatically

## Example Usage
```bash
# Simple blog post
construct g Post title:string content:text published:bool

# E-commerce product
construct g Product name:string description:text price:float stock:uint featured_image:image

# User profile
construct g Profile user_id:uint bio:text avatar:image website:url
```

Next: Implement the generator!
