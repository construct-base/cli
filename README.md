# Construct CLI

The official command-line tool for the [Construct Framework](https://github.com/construct-base/core).

Construct is a modern full-stack framework combining **Vue 3** (frontend) and **Base Go** (backend) into one cohesive system.

## Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/construct-base/cli/main/install.sh | bash
```

This will:
- Download the latest CLI binary for your platform
- Install to `~/.base/bin/construct`
- Add to your PATH automatically

### Manual Installation

Download the binary for your platform from the [latest release](https://github.com/construct-base/cli/releases/latest):

**macOS (Apple Silicon):**
```bash
curl -L https://github.com/construct-base/cli/releases/latest/download/construct-darwin-arm64.tar.gz | tar xz
mkdir -p ~/.base/bin
mv construct-darwin-arm64 ~/.base/bin/construct
chmod +x ~/.base/bin/construct
echo 'export PATH="$HOME/.base/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**macOS (Intel):**
```bash
curl -L https://github.com/construct-base/cli/releases/latest/download/construct-darwin-amd64.tar.gz | tar xz
mkdir -p ~/.base/bin
mv construct-darwin-amd64 ~/.base/bin/construct
chmod +x ~/.base/bin/construct
echo 'export PATH="$HOME/.base/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**Linux (amd64):**
```bash
curl -L https://github.com/construct-base/cli/releases/latest/download/construct-linux-amd64.tar.gz | tar xz
mkdir -p ~/.base/bin
mv construct-linux-amd64 ~/.base/bin/construct
chmod +x ~/.base/bin/construct
echo 'export PATH="$HOME/.base/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Build from Source

```bash
git clone https://github.com/construct-base/cli.git
cd cli
go build -o construct main.go
mv construct ~/.base/bin/
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

## Verify Installation

```bash
construct --version
```

## Development

```bash
# Clone the repository
git clone https://github.com/construct-base/cli.git
cd cli

# Install dependencies
go mod download

# Build
go build -o construct main.go

# Test
./construct --help
```

## Related Projects

- [construct-core](https://github.com/construct-base/core) - The Construct framework itself
- [Documentation](https://github.com/construct-base/cli) - Full documentation

## License

MIT