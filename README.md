# Construct CLI

The official command-line tool for the [Construct Framework](https://github.com/construct-go/core).

Construct is a modern full-stack framework combining **Vue 3** (frontend) and **Base Go** (backend) into one cohesive system.

## Installation

```bash
# Install globally (macOS/Linux)
curl -fsSL https://raw.githubusercontent.com/construct-go/cli/main/install.sh | sh

# Or build from source
go install github.com/construct-go/cli@latest
```

## Quick Start

```bash
# Create a new project
construct new my-blog

# Navigate to project
cd my-blog

# Generate a CRUD structure
construct g Post title:string content:text published:bool

# Start development servers
construct dev
```

## Commands

### `construct new [project]`
Create a new Construct project by downloading the latest framework.

```bash
construct new my-app
```

### `construct generate [resource] [fields...]`
Aliases: `g`, `gen`

Generate full-stack CRUD scaffolding (Go backend + Vue frontend).

```bash
# Basic structure
construct g Post title:string content:text published:bool

# With relationships
construct g Article title:string category_id:uint author_id:uint

# All field types
construct g Product name:string price:float stock:uint featured:bool
```

**Supported field types:**
- `string`, `text` - Text fields
- `int`, `uint` - Integer fields
- `float`, `float64` - Decimal fields
- `bool`, `boolean` - Boolean fields
- `date`, `datetime`, `time` - Date/time fields

**What gets generated:**
- **Backend** (`api/{resource}/`): service.go, controller.go, module.go, validator.go
- **Model** (`app/models/`): {resource}.go
- **Frontend** (`vue/structures/{resource}/`): index.vue, composable.ts, types.ts
- **Auto-registration**: Module added to `api/init.go`

### `construct dev`
Start development servers for both Go (port 8100) and Vue (port 3100).

```bash
construct dev
```

### `construct build`
Build the application for production. Creates a `dist/` directory with:
- Compiled Go binary
- Built Vue SPA
- Runtime directories (storage/, logs/)

```bash
construct build
```

### `construct start`
Start the production server from the `dist/` directory.

```bash
construct start
```

## Architecture

Construct uses **HMVC (Hierarchical Model-View-Controller)** on both sides:

**Backend (Go):**
```
api/{resource}/
├── service.go      # Business logic
├── controller.go   # HTTP handlers
├── module.go       # Module registration
└── validator.go    # Input validation
```

**Frontend (Vue):**
```
vue/structures/{resource}/
├── index.vue       # Main page with CRUD
├── composable.ts   # API integration
└── types.ts        # TypeScript types
```

## Example Workflow

```bash
# 1. Create project
construct new blog

# 2. Navigate to project
cd blog

# 3. Generate blog post structure
construct g Post title:string content:text published:bool category_id:uint

# 4. Generate category structure
construct g Category name:string description:text

# 5. Start development
construct dev

# 6. Visit http://localhost:3100
# - Frontend: Vue dev server
# - Backend API: http://localhost:8100/api
```

## Development

```bash
# Clone the repository
git clone https://github.com/construct-go/cli.git
cd cli

# Install dependencies
go mod download

# Build
go build -o construct main.go

# Test
./construct --help
```

## Related Projects

- [construct-core](https://github.com/construct-go/core) - The Construct framework itself
- [base-core](https://github.com/base-go/base-core) - The Go backend framework

## License

MIT