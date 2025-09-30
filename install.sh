#!/bin/bash

# Construct CLI Installation Script
# Installs the construct CLI to ~/.base/construct and adds it to PATH

set -e

echo "🚀 Installing Construct CLI..."
echo ""

# Detect OS and Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Map architecture names
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        echo "   Supported: amd64, arm64"
        exit 1
        ;;
esac

# Map OS names
case "$OS" in
    linux)
        PLATFORM="linux"
        ;;
    darwin)
        PLATFORM="darwin"
        ;;
    *)
        echo "❌ Unsupported operating system: $OS"
        echo "   Supported: linux, darwin (macOS)"
        exit 1
        ;;
esac

echo "📋 System detected:"
echo "   OS: $PLATFORM"
echo "   Architecture: $ARCH"
echo ""

# Get latest release version from GitHub
echo "🔍 Fetching latest release..."
LATEST_VERSION=$(curl -s https://api.github.com/repos/construct-base/cli/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "❌ Failed to fetch latest release version"
    echo "   Please check your internet connection and try again"
    exit 1
fi

echo "   Latest version: $LATEST_VERSION"
echo ""

# Construct download URL
BINARY_NAME="construct-${PLATFORM}-${ARCH}"
DOWNLOAD_URL="https://github.com/construct-base/cli/releases/download/${LATEST_VERSION}/${BINARY_NAME}.tar.gz"

echo "📥 Downloading Construct CLI..."
echo "   URL: $DOWNLOAD_URL"
echo ""

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download binary
if ! curl -L -o "$TMP_DIR/construct.tar.gz" "$DOWNLOAD_URL"; then
    echo "❌ Failed to download CLI"
    echo "   URL: $DOWNLOAD_URL"
    exit 1
fi

# Extract binary
echo "📦 Extracting..."
tar -xzf "$TMP_DIR/construct.tar.gz" -C "$TMP_DIR"

# Create installation directory
INSTALL_DIR="$HOME/.base"
BIN_DIR="$INSTALL_DIR/bin"

echo "📁 Creating installation directory..."
mkdir -p "$BIN_DIR"

# Move binary
echo "🔧 Installing binary..."
mv "$TMP_DIR/$BINARY_NAME" "$BIN_DIR/construct"
chmod +x "$BIN_DIR/construct"

echo ""
echo "✅ Construct CLI installed successfully!"
echo ""
echo "📍 Installation location: $BIN_DIR/construct"
echo ""

# Detect shell and update PATH
SHELL_NAME=$(basename "$SHELL")
PROFILE_FILE=""

case "$SHELL_NAME" in
    bash)
        if [ -f "$HOME/.bashrc" ]; then
            PROFILE_FILE="$HOME/.bashrc"
        elif [ -f "$HOME/.bash_profile" ]; then
            PROFILE_FILE="$HOME/.bash_profile"
        fi
        ;;
    zsh)
        PROFILE_FILE="$HOME/.zshrc"
        ;;
    fish)
        PROFILE_FILE="$HOME/.config/fish/config.fish"
        ;;
    *)
        echo "⚠️  Unknown shell: $SHELL_NAME"
        ;;
esac

# Check if PATH is already configured
if echo "$PATH" | grep -q "$BIN_DIR"; then
    echo "✓ PATH is already configured"
else
    echo "⚙️  Configuring PATH..."

    if [ -n "$PROFILE_FILE" ]; then
        # Add to profile
        if [ "$SHELL_NAME" = "fish" ]; then
            echo "set -gx PATH $BIN_DIR \$PATH" >> "$PROFILE_FILE"
        else
            echo "" >> "$PROFILE_FILE"
            echo "# Construct CLI" >> "$PROFILE_FILE"
            echo "export PATH=\"$BIN_DIR:\$PATH\"" >> "$PROFILE_FILE"
        fi

        echo "   ✓ Added to $PROFILE_FILE"
        echo ""
        echo "📝 To use construct now, run:"
        echo "   source $PROFILE_FILE"
        echo ""
        echo "   Or restart your terminal"
    else
        echo "   ⚠️  Could not detect profile file"
        echo ""
        echo "📝 Please add this to your shell profile manually:"
        echo "   export PATH=\"$BIN_DIR:\$PATH\""
    fi
fi

echo ""
echo "🎉 Installation complete!"
echo ""
echo "🚀 Quick start:"
echo "   construct new my-blog      # Create a new project"
echo "   cd my-blog"
echo "   construct dev              # Start development"
echo ""
echo "📖 For more info, visit: https://github.com/construct-base/cli"
echo ""
