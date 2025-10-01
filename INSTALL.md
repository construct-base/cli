# Construct Installation Guide

## Quick Install (Recommended)

Install the Construct CLI with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/construct-base/cli/main/install.sh | bash
```

This will:
- ✅ Detect your OS and architecture automatically
- ✅ Download the latest CLI binary
- ✅ Install to `~/.base/bin/construct`
- ✅ Add to your PATH automatically
- ✅ Support for macOS (Intel/Apple Silicon) and Linux (amd64/arm64)

## Verify Installation

```bash
construct --version
```

## Create Your First Project

```bash
construct new my-blog
cd my-blog
construct dev
```

Visit `http://localhost:3100` to see your app running!

## What Gets Installed

- **Binary**: `~/.base/bin/construct`
- **PATH**: Automatically added to your shell profile (`~/.zshrc`, `~/.bashrc`, or `~/.config/fish/config.fish`)

## Manual Installation

If you prefer to install manually, see the [CLI README](https://github.com/construct-base/cli#installation) for platform-specific instructions.

## Updating

To update to the latest version, simply run the install script again:

```bash
curl -fsSL https://raw.githubusercontent.com/construct-base/cli/main/install.sh | bash
```

## Uninstalling

```bash
rm ~/.base/bin/construct
```

Then remove the PATH export from your shell profile:
- **Bash**: Edit `~/.bashrc`
- **Zsh**: Edit `~/.zshrc`
- **Fish**: Edit `~/.config/fish/config.fish`

## Troubleshooting

### Command not found after installation

Reload your shell profile:

```bash
# For bash
source ~/.bashrc

# For zsh
source ~/.zshrc

# For fish
source ~/.config/fish/config.fish
```

Or restart your terminal.

### Permission denied

Make sure the binary is executable:

```bash
chmod +x ~/.base/bin/construct
```

### Download fails

Check your internet connection and make sure you can access GitHub:

```bash
curl -I https://github.com/construct-base/cli
```

## System Requirements

- **OS**: macOS or Linux (Windows support coming soon)
- **Architecture**: amd64 (x86_64) or arm64 (aarch64)
- **Tools**: curl, tar
- **For development**: Go 1.21+, Node.js 18+, Git

## Next Steps

- [CLI Documentation](https://github.com/construct-base/cli)
- [Core Framework](https://github.com/construct-base/core)
- [Quick Start Guide](https://github.com/construct-base/cli#quick-start)
