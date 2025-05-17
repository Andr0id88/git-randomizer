#!/bin/bash
# One-liner install script for gitr (git-randomizer)

# --- Configuration ---
REPO="Andr0id88/git-randomizer"
BINARY_NAME="gitr" # The name you want the installed binary to have
INSTALL_DIR_CANDIDATES=(
    "$HOME/.local/bin"  # Common on Linux, good practice
    "$HOME/bin"         # Older common practice
    "/usr/local/bin"    # Common system-wide, might need sudo
)
# --- End Configuration ---

set -e # Exit on error
set -o pipefail # Ensure pipes fail correctly

# Function to detect OS and Arch
get_os_arch() {
    local os_name arch_name
    case "$(uname -s)" in
        Linux*)     os_name=linux;;
        Darwin*)    os_name=darwin;;
        *)          echo "Unsupported OS: $(uname -s)"; exit 1;;
    esac

    case "$(uname -m)" in
        x86_64)     arch_name=amd64;;
        aarch64)    arch_name=arm64;; # Linux ARM64
        arm64)      arch_name=arm64;; # macOS Apple Silicon
        *)          echo "Unsupported architecture: $(uname -m)"; exit 1;;
    esac
    echo "${os_name}-${arch_name}"
}

# Function to get the latest release tag from GitHub API
get_latest_release_tag() {
    # Uses curl and jq (common utils). Falls back if jq is not available.
    local latest_tag
    if command -v jq >/dev/null 2>&1; then
        latest_tag=$(curl --silent "https://api.github.com/repos/$1/releases/latest" | jq -r .tag_name)
    else
        # Fallback if jq is not installed (less robust, relies on GitHub API response order)
        latest_tag=$(curl --silent "https://api.github.com/repos/$1/releases" | grep -oP '"tag_name": "\K[^"]+' | head -n 1)
    fi

    if [ -z "$latest_tag" ]; then
        echo "Could not fetch the latest release tag for $1. Please check the repository or network."
        exit 1
    fi
    echo "$latest_tag"
}


OS_ARCH=$(get_os_arch)
OS_NAME=$(echo "$OS_ARCH" | cut -d'-' -f1)
ARCH_NAME=$(echo "$OS_ARCH" | cut -d'-' -f2)

echo "Detected OS: $OS_NAME, Arch: $ARCH_NAME"

LATEST_TAG=$(get_latest_release_tag "$REPO")
echo "Latest release tag: $LATEST_TAG"

DOWNLOAD_FILENAME="${BINARY_NAME}-${OS_NAME}-${ARCH_NAME}"
if [ "$OS_NAME" = "windows" ]; then # Though this script is primarily for Unix-like
    DOWNLOAD_FILENAME="${DOWNLOAD_FILENAME}.exe"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${DOWNLOAD_FILENAME}"

echo "Downloading $BINARY_NAME from $DOWNLOAD_URL"

# Create a temporary directory for download
TMP_DIR=$(mktemp -d)
trap 'rm -rf -- "$TMP_DIR"' EXIT # Clean up temp dir on exit

curl -Lo "${TMP_DIR}/${BINARY_NAME}" "$DOWNLOAD_URL"
if [ $? -ne 0 ]; then
    echo "Error: Download failed. Please check the URL or your network connection."
    exit 1
fi

echo "Download complete. Making it executable..."
chmod +x "${TMP_DIR}/${BINARY_NAME}"

# Find a suitable installation directory
INSTALL_DIR=""
for dir_candidate in "${INSTALL_DIR_CANDIDATES[@]}"; do
    if [ -d "$dir_candidate" ] && [ -w "$dir_candidate" ]; then
        INSTALL_DIR="$dir_candidate"
        break
    elif [ ! -d "$dir_candidate" ]; then # Try to create if it doesn't exist (user-owned paths)
        if [[ "$dir_candidate" == "$HOME/"* ]]; then # Only try to create if it's in $HOME
            echo "Attempting to create directory: $dir_candidate"
            if mkdir -p "$dir_candidate"; then
                if [ -w "$dir_candidate" ]; then
                    INSTALL_DIR="$dir_candidate"
                    break
                fi
            else
                echo "Could not create $dir_candidate."
            fi
        fi
    fi
done

# Fallback for /usr/local/bin if no user-writable path found, might need sudo
if [ -z "$INSTALL_DIR" ]; then
    if [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    else # Attempt with sudo if /usr/local/bin exists but isn't writable by user
        if command -v sudo >/dev/null 2>&1 && [ -d "/usr/local/bin" ]; then
            echo "No user-writable install directory found in your PATH."
            echo "Attempting to install to /usr/local/bin using sudo."
            echo "Please enter your password if prompted."
            if sudo mkdir -p "/usr/local/bin" && sudo test -w "/usr/local/bin"; then
                 sudo mv "${TMP_DIR}/${BINARY_NAME}" "/usr/local/bin/${BINARY_NAME}"
                 echo ""
                 echo "$BINARY_NAME installed successfully to /usr/local/bin/${BINARY_NAME}"
                 echo "You may need to open a new terminal or run 'hash -r' (bash) or 'rehash' (zsh) for the shell to find it."
                 exit 0 # Successfully installed with sudo
            else
                echo "Failed to install to /usr/local/bin even with sudo."
            fi
        fi
    fi
fi


if [ -z "$INSTALL_DIR" ]; then
    echo "Error: Could not find a writable installation directory from the candidates: ${INSTALL_DIR_CANDIDATES[*]}"
    echo "Please create one of these directories or ensure it's writable, then try again."
    echo "Alternatively, you can manually move '${TMP_DIR}/${BINARY_NAME}' to a directory in your PATH."
    exit 1
fi

echo "Installing $BINARY_NAME to ${INSTALL_DIR}/${BINARY_NAME}"
mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"

# Check if the installation directory is in PATH
if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    echo ""
    echo "Warning: Installation directory '${INSTALL_DIR}' is not in your PATH."
    echo "You will need to add it to your PATH to run '${BINARY_NAME}' directly."
    echo "You can do this by adding the following line to your shell configuration file"
    echo "(e.g., ~/.bashrc, ~/.zshrc, ~/.profile, or ~/.config/fish/config.fish):"
    echo ""
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    echo ""
    echo "After adding it, open a new terminal or source your config file (e.g., 'source ~/.bashrc')."
else
    echo ""
    echo "$BINARY_NAME installed successfully to ${INSTALL_DIR}/${BINARY_NAME}"
    echo "You may need to open a new terminal or run 'hash -r' (bash) or 'rehash' (zsh) for the shell to find it."
fi

echo "Installation complete!"
