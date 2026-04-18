#!/usr/bin/env bash
# Short interactive installer for gitmap on Linux / macOS.
#
# Prompts for an install folder (with a sensible default), then delegates
# to the canonical gitmap/scripts/install.sh with that path.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v3/main/install-quick.sh | bash
#   ./install-quick.sh
#   ./install-quick.sh --dir /opt/gitmap

set -euo pipefail

REPO="alimtvnetwork/gitmap-v3"
INSTALLER_URL="https://raw.githubusercontent.com/${REPO}/main/gitmap/scripts/install.sh"
DEFAULT_DIR="${HOME}/.local/bin"

INSTALL_DIR=""
VERSION=""

while [ $# -gt 0 ]; do
    case "$1" in
        --dir)     INSTALL_DIR="$2"; shift 2 ;;
        --version) VERSION="$2";     shift 2 ;;
        -h|--help)
            sed -n '2,12p' "$0"
            exit 0
            ;;
        *)
            printf '  Unknown argument: %s\n' "$1" >&2
            exit 1
            ;;
    esac
done

prompt_dir() {
    printf '\n'
    printf '  \033[36mgitmap quick installer\033[0m\n'
    printf '  \033[90m---------------------\033[0m\n'
    printf '  Choose install folder. Press Enter to accept the default.\n'
    printf '  \033[90mDefault: %s\033[0m\n' "${DEFAULT_DIR}"
    printf '  Install path: '

    # Read from the controlling terminal so it works under `curl | bash`.
    if [ -r /dev/tty ]; then
        IFS= read -r answer < /dev/tty || answer=""
    else
        IFS= read -r answer || answer=""
    fi

    if [ -z "${answer}" ]; then
        echo "${DEFAULT_DIR}"
    else
        echo "${answer}"
    fi
}

if [ -z "${INSTALL_DIR}" ]; then
    INSTALL_DIR="$(prompt_dir)"
fi

printf '\n  \033[32mInstalling gitmap to: %s\033[0m\n\n' "${INSTALL_DIR}"

ARGS=(--dir "${INSTALL_DIR}")
if [ -n "${VERSION}" ]; then
    ARGS+=(--version "${VERSION}")
fi

curl -fsSL "${INSTALLER_URL}" | bash -s -- "${ARGS[@]}"
