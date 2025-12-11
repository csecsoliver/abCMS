#!/bin/bash
set -e

# setup.sh - Environment setup script for abCMS
# Replicates the environment setup from shell.nix for Ubuntu/Debian systems.

echo "Starting setup for abCMS..."

# Helper function for sudo
SUDO=""
if [ "$EUID" -ne 0 ]; then
    if command -v sudo >/dev/null 2>&1; then
        SUDO="sudo"
    else
        echo "This script requires root privileges to install packages. Please run as root or ensure sudo is available."
        exit 1
    fi
fi

# Detect OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$NAME
fi

if [[ "$OS" == "Ubuntu" ]] || [[ "$OS" == "Debian GNU/Linux" ]]; then
    echo "Installing system dependencies for $OS..."

    # Update package list
    $SUDO apt-get update

    # Install basic tools
    $SUDO apt-get install -y wget gnupg ca-certificates lsb-release build-essential git curl pkg-config

    # Install OpenResty
    if ! command -v openresty >/dev/null 2>&1; then
        echo "Installing OpenResty..."
        # Add OpenResty GPG key
        wget -O - https://openresty.org/package/pubkey.gpg | $SUDO gpg --dearmor --yes -o /usr/share/keyrings/openresty.gpg
        # Add OpenResty repository
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/openresty.gpg] http://openresty.org/package/ubuntu $(lsb_release -sc) main" | $SUDO tee /etc/apt/sources.list.d/openresty.list

        $SUDO apt-get update
        $SUDO apt-get install -y openresty
    else
        echo "OpenResty is already installed."
    fi

    # Install other system dependencies
    # shell.nix requirements:
    # imagemagick, libuuid, libargon2, openssl, sqlite, lua5.1, luarocks, tmux
    echo "Installing libraries..."
    $SUDO apt-get install -y \
        imagemagick libmagickwand-dev \
        uuid-dev \
        libargon2-dev \
        libssl-dev \
        sqlite3 libsqlite3-dev \
        lua5.1 liblua5.1-0-dev \
        luarocks \
        tmux

else
    echo "Warning: This script assumes Ubuntu/Debian. You may need to install dependencies manually for your OS ($OS)."
    echo "Required: OpenResty, ImageMagick (dev), UUID (dev), Argon2 (dev), OpenSSL (dev), SQLite (dev), Lua 5.1 (dev), LuaRocks, Tmux."
    read -p "Press Enter to continue trying setup with existing tools, or Ctrl+C to abort..."
fi

# Setup Lua Environment (Local)
echo "Setting up local Lua environment..."

# Initialize/Prepare local luarocks tree
mkdir -p logs

# Add local bin to PATH for the duration of this script
export PATH="$HOME/.luarocks/bin:$PATH"

# Check if we can run luarocks with Lua 5.1
if ! luarocks --lua-version=5.1 list >/dev/null 2>&1; then
    echo "Error: LuaRocks cannot list packages for Lua 5.1."
    echo "Ensure 'lua5.1' and 'liblua5.1-0-dev' are installed correctly."
    # Attempt to continue? No, likely will fail.
    # On some systems luarocks might default to 5.1 if only 5.1 is installed, but we explicity use flag.
    echo "Checking if we can config it..."
fi

# Function to install rock
install_rock() {
    PACKAGE=$1
    shift
    echo "Installing Lua rock: $PACKAGE"
    luarocks --local --lua-version=5.1 install $PACKAGE "$@"
}

# Install Lua packages from shellHook
install_rock sendmail
install_rock luasec
install_rock lume
install_rock moonscript
install_rock bcrypt
install_rock Lua-curl
install_rock markdown

# Magick setup
# Attempt to detect paths using pkg-config if available
if pkg-config --exists MagickWand; then
    MAGICK_INCDIR=$(pkg-config --cflags-only-I MagickWand | sed 's/-I//g' | awk '{print $1}')
    if [ -n "$MAGICK_INCDIR" ]; then
        echo "Found MagickWand headers at: $MAGICK_INCDIR"
        install_rock magick MAGICK_INCDIR="$MAGICK_INCDIR"
    else
        install_rock magick
    fi
else
    # Fallback or standard install
    install_rock magick
fi

# Lapis (needs OpenSSL)
install_rock lapis CRYPTO_DIR=/usr OPENSSL_DIR=/usr

# Argon2
install_rock argon2 ARGON2_DIR=/usr

# UUID
install_rock lua-uuid UUID_DIR=/usr

# SQLite3
install_rock lsqlite3 SQLITE_DIR=/usr

echo ""
echo "----------------------------------------------------------------"
echo "Setup Complete!"
echo ""
echo "To start working, run the following command in your shell:"
echo "  eval \"\$(luarocks path --bin)\""
echo ""
echo "Also set these environment variables (add to ~/.bashrc or similar):"
echo "  export LAPIS_OPENRESTY=\$(which openresty)"
# Try to find libMagickWand shared object to suggest LD_LIBRARY_PATH/MAGICK_WAND_LIBRARY
WAND_LIB=""
if pkg-config --exists MagickWand; then
    LIB_DIR=$(pkg-config --libs-only-L MagickWand | sed 's/-L//g' | awk '{print $1}')
    if [ -n "$LIB_DIR" ]; then
        WAND_LIB=$(find "$LIB_DIR" -name "libMagickWand*.so" | head -n 1)
    fi
fi
if [ -z "$WAND_LIB" ]; then
    # Guess
    WAND_LIB=$(find /usr/lib -name "libMagickWand*.so" 2>/dev/null | head -n 1)
fi

if [ -n "$WAND_LIB" ]; then
    echo "  export MAGICK_WAND_LIBRARY=\"$WAND_LIB\""
fi
echo "----------------------------------------------------------------"
