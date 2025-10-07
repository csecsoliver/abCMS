#!/bin/bash
set -e

# Use Poetry's environment and freeze dependencies
poetry run pip freeze > requirements.txt

# Create a temporary build directory
BUILD_DIR="./pex_build"
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# Copy ALL files to build directory (flat structure with subdirs)
cp *.py "$BUILD_DIR/"
cp *.html "$BUILD_DIR/" 2>/dev/null || true
cp *.tpl "$BUILD_DIR/" 2>/dev/null || true
cp .env "$BUILD_DIR/" 2>/dev/null || touch "$BUILD_DIR/.env"

# Copy resources directory
cp -r ./resources "$BUILD_DIR/"

# # Copy data directories (not necessary tho)
# cp -r users "$BUILD_DIR/" 2>/dev/null || mkdir -p "$BUILD_DIR/users"
# cp -r deployments "$BUILD_DIR/" 2>/dev/null || mkdir -p "$BUILD_DIR/deployments"
# cp -r md_blog_content "$BUILD_DIR/" 2>/dev/null || mkdir -p "$BUILD_DIR/md_blog_content"

# Build the PEX with the build directory as the source
poetry run pex -D "$BUILD_DIR" -r requirements.txt -o app.pex --scie lazy -e main

# Clean up
rm -rf "$BUILD_DIR"
rm requirements.txt

echo "Build complete! Run with: ./app"
