#!/usr/bin/env bash
# Re-exec under bash if invoked via sh/dash (which lack pipefail, local, etc.)
if [ -z "${BASH_VERSION:-}" ]; then
    if command -v bash >/dev/null 2>&1; then
        case "${0##*/}" in
            sh|dash|ash|ksh|mksh)
                exec bash -s -- "$@"
                ;;
        esac

        exec bash "$0" "$@"
    else
        printf '\033[31m  Error: bash is required but not found. Install bash first.\033[0m\n' >&2
        exit 1
    fi
fi
# ─────────────────────────────────────────────────────────────────────
# gitmap installer for Linux and macOS
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v2/main/gitmap/scripts/install.sh | bash
#
# Options:
#   --version <tag>    Install a specific version (e.g. v2.55.0). Default: latest.
#   --dir <path>       Target directory. Default: ~/.local/bin
#   --arch <arch>      Force architecture (amd64, arm64). Default: auto-detect.
#   --no-path          Skip adding install directory to PATH.
#
# Examples:
#   curl -fsSL .../install.sh | bash
#   curl -fsSL .../install.sh | bash -s -- --version v2.55.0
#   ./install.sh --dir /opt/gitmap --arch arm64
# ─────────────────────────────────────────────────────────────────────

set -euo pipefail

REPO="alimtvnetwork/gitmap-v2"
BINARY_NAME="gitmap"
TMP_DIR=""
PATH_SHELL=""
PATH_TARGET=""
PATH_LINE=""
PATH_STATUS=""
PATH_RELOAD=""

cleanup() {
    if [ -n "${TMP_DIR}" ] && [ -d "${TMP_DIR}" ]; then
        rm -rf "${TMP_DIR}"
    fi
}
trap cleanup EXIT

# ── Logging helpers ─────────────────────────────────────────────────

step()  { printf '  \033[36m%s\033[0m\n' "$*" >&2; }
ok()    { printf '  \033[32m%s\033[0m\n' "$*" >&2; }
err()   { printf '  \033[31m%s\033[0m\n' "$*" >&2; }

# ── Detect OS ───────────────────────────────────────────────────────

detect_os() {
    local uname_out
    uname_out="$(uname -s)"
    case "${uname_out}" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*)
            err "Windows detected. Use the PowerShell installer instead:"
            err "  irm https://raw.githubusercontent.com/${REPO}/main/gitmap/scripts/install.ps1 | iex"
            exit 1
            ;;
        *)
            err "Unsupported OS: ${uname_out}"
            exit 1
            ;;
    esac
}

# ── Detect architecture ────────────────────────────────────────────

detect_arch() {
    local arch_flag="$1"
    if [ -n "${arch_flag}" ]; then
        echo "${arch_flag}"
        return
    fi

    local machine
    machine="$(uname -m)"
    case "${machine}" in
        x86_64|amd64)   echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        *)
            err "Unsupported architecture: ${machine}"
            exit 1
            ;;
    esac
}

# ── Resolve version (latest or pinned) ─────────────────────────────

resolve_version() {
    local version="$1"
    if [ -n "${version}" ]; then
        echo "${version}"
        return
    fi

    step "Fetching latest release..."
    local url="https://api.github.com/repos/${REPO}/releases/latest"
    local tag

    if command -v curl >/dev/null 2>&1; then
        tag="$(curl -fsSL "${url}" | grep '"tag_name"' | head -1 | sed -E 's/.*"tag_name"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/')"
    elif command -v wget >/dev/null 2>&1; then
        tag="$(wget -qO- "${url}" | grep '"tag_name"' | head -1 | sed -E 's/.*"tag_name"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/')"
    else
        err "Neither curl nor wget found. Cannot fetch latest release."
        exit 1
    fi

    if [ -z "${tag}" ]; then
        err "Failed to determine latest version."
        exit 1
    fi

    echo "${tag}"
}

# ── Download helper ────────────────────────────────────────────────

download() {
    local url="$1" dest="$2"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL -o "${dest}" "${url}"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "${dest}" "${url}"
    else
        err "Neither curl nor wget found."
        exit 1
    fi
}

# ── Download and verify asset ──────────────────────────────────────

download_asset() {
    local version="$1" os="$2" arch="$3"
    local asset_name="${BINARY_NAME}-${version}-${os}-${arch}.tar.gz"
    local base_url="https://github.com/${REPO}/releases/download/${version}"
    local asset_url="${base_url}/${asset_name}"
    local checksum_url="${base_url}/checksums.txt"

    # TMP_DIR is set by the caller (main).

    local archive_path="${TMP_DIR}/${asset_name}"
    local checksum_path="${TMP_DIR}/checksums.txt"

    step "Downloading ${asset_name} (${version})..."
    download "${asset_url}" "${archive_path}"
    download "${checksum_url}" "${checksum_path}"

    # Verify checksum
    step "Verifying checksum..."
    local expected_line
    expected_line="$(grep "${asset_name}" "${checksum_path}" || true)"
    if [ -z "${expected_line}" ]; then
        # Try .zip variant (some releases may only have zip)
        asset_name="${BINARY_NAME}-${version}-${os}-${arch}.zip"
        asset_url="${base_url}/${asset_name}"
        archive_path="${TMP_DIR}/${asset_name}"

        step "Trying .zip variant..."
        download "${asset_url}" "${archive_path}"
        expected_line="$(grep "${asset_name}" "${checksum_path}" || true)"

        if [ -z "${expected_line}" ]; then
            err "Asset not found in checksums.txt"
            err "Tried: ${BINARY_NAME}-${version}-${os}-${arch}.tar.gz"
            err "Tried: ${asset_name}"
            exit 1
        fi
    fi

    local expected_hash
    expected_hash="$(echo "${expected_line}" | awk '{print $1}')"

    local actual_hash
    if command -v sha256sum >/dev/null 2>&1; then
        actual_hash="$(sha256sum "${archive_path}" | awk '{print $1}')"
    elif command -v shasum >/dev/null 2>&1; then
        actual_hash="$(shasum -a 256 "${archive_path}" | awk '{print $1}')"
    else
        err "No SHA256 tool found (sha256sum or shasum required)."
        exit 1
    fi

    if [ "${actual_hash}" != "${expected_hash}" ]; then
        err "Checksum mismatch!"
        err "  Expected: ${expected_hash}"
        err "  Got:      ${actual_hash}"
        exit 1
    fi

    ok "Checksum verified."
    echo "${archive_path}"
}

# ── Extract and install binary ─────────────────────────────────────

install_binary() {
    local archive_path="$1" install_dir="$2" os="$3" arch="$4" version="$5"

    step "Installing to ${install_dir}..."
    mkdir -p "${install_dir}"

    local extract_dir="${TMP_DIR}/extract"
    mkdir -p "${extract_dir}"

    # Extract based on file type
    case "${archive_path}" in
        *.tar.gz|*.tgz)
            tar -xzf "${archive_path}" -C "${extract_dir}"
            ;;
        *.zip)
            if command -v unzip >/dev/null 2>&1; then
                unzip -qo "${archive_path}" -d "${extract_dir}"
            else
                err "unzip not found. Cannot extract .zip archive."
                exit 1
            fi
            ;;
        *)
            err "Unknown archive format: ${archive_path}"
            exit 1
            ;;
    esac

    # Search for the binary — exact name or versioned pattern
    local binary_path=""
    local candidate

    # Priority 1: exact binary name
    candidate="$(find "${extract_dir}" -type f -name "${BINARY_NAME}" | head -1)"
    if [ -n "${candidate}" ]; then
        binary_path="${candidate}"
    fi

    # Priority 2: platform-specific name (e.g. gitmap-linux-amd64)
    if [ -z "${binary_path}" ]; then
        candidate="$(find "${extract_dir}" -type f -name "${BINARY_NAME}-${os}-${arch}" | head -1)"
        if [ -n "${candidate}" ]; then
            binary_path="${candidate}"
        fi
    fi

    # Priority 3: versioned pattern (e.g. gitmap-v2.55.0-linux-amd64)
    if [ -z "${binary_path}" ]; then
        candidate="$(find "${extract_dir}" -type f -regex ".*/${BINARY_NAME}-v[0-9][0-9.]*-${os}-${arch}" | head -1)"
        if [ -n "${candidate}" ]; then
            binary_path="${candidate}"
        fi
    fi

    # Priority 4: any executable file in the archive
    if [ -z "${binary_path}" ]; then
        candidate="$(find "${extract_dir}" -type f -executable | head -1)"
        if [ -n "${candidate}" ]; then
            binary_path="${candidate}"
        fi
    fi

    if [ -z "${binary_path}" ]; then
        err "Archive did not contain a recognizable binary."
        err "Archive contents:"
        find "${extract_dir}" -type f | while read -r f; do
            err "  ${f}"
        done
        exit 1
    fi

    local target_path="${install_dir}/${BINARY_NAME}"

    # Rename-first strategy for running binary
    if [ -f "${target_path}" ]; then
        mv -f "${target_path}" "${target_path}.old" 2>/dev/null || true
    fi

    mv -f "${binary_path}" "${target_path}"
    chmod +x "${target_path}"

    # Cleanup .old
    rm -f "${target_path}.old" 2>/dev/null || true

    if [ ! -f "${target_path}" ]; then
        err "Install failed: ${BINARY_NAME} was not written to ${install_dir}"
        exit 1
    fi

    ok "Installed ${BINARY_NAME} to ${install_dir}"
}

# ── Add to PATH ────────────────────────────────────────────────────

# add_path_to_profile writes an export line to a single profile file (idempotent).
# Returns 0 if written, 1 if already present.
add_path_to_profile() {
    local dir="$1" profile_file="$2" is_fish="$3"

    local path_line
    if [ "${is_fish}" = true ]; then
        path_line="fish_add_path ${dir}"
    else
        path_line="export PATH=\"\$PATH:${dir}\""
    fi

    if [ -f "${profile_file}" ] && grep -qF "${dir}" "${profile_file}"; then
        return 1
    fi

    mkdir -p "$(dirname "${profile_file}")"
    printf '\n# Added by gitmap installer\n%s\n' "${path_line}" >> "${profile_file}"
    return 0
}

add_to_path() {
    local dir="$1"
    local has_session_path=false

    case ":${PATH}:" in
        *":${dir}:"*)
            has_session_path=true
            ;;
    esac

    # Detect primary shell
    local shell_name
    shell_name="$(basename "${SHELL:-/bin/bash}")"
    PATH_SHELL="${shell_name}"

    local primary_profile=""
    local profiles_written=""
    local profiles_skipped=""

    # ── Write to all relevant POSIX/bash/zsh profiles ──────────────
    # This ensures gitmap is available regardless of which shell the user opens.

    # zsh profiles (both, to cover login + interactive shells)
    if [ "${shell_name}" = "zsh" ] || [ -f "${HOME}/.zshrc" ] || [ -f "${HOME}/.zprofile" ]; then
        # .zshrc — interactive shells (most terminal emulators)
        if add_path_to_profile "${dir}" "${HOME}/.zshrc" false; then
            profiles_written="${profiles_written} ~/.zshrc"
        else
            profiles_skipped="${profiles_skipped} ~/.zshrc"
        fi
        # .zprofile — login shells (macOS Terminal.app)
        if add_path_to_profile "${dir}" "${HOME}/.zprofile" false; then
            profiles_written="${profiles_written} ~/.zprofile"
        else
            profiles_skipped="${profiles_skipped} ~/.zprofile"
        fi
    fi

    # bash profiles
    if [ "${shell_name}" = "bash" ] || [ -f "${HOME}/.bashrc" ] || [ -f "${HOME}/.bash_profile" ]; then
        if add_path_to_profile "${dir}" "${HOME}/.bashrc" false; then
            profiles_written="${profiles_written} ~/.bashrc"
        else
            profiles_skipped="${profiles_skipped} ~/.bashrc"
        fi
        if [ -f "${HOME}/.bash_profile" ]; then
            if add_path_to_profile "${dir}" "${HOME}/.bash_profile" false; then
                profiles_written="${profiles_written} ~/.bash_profile"
            else
                profiles_skipped="${profiles_skipped} ~/.bash_profile"
            fi
        fi
    fi

    # POSIX ~/.profile — catch-all for sh and other POSIX shells
    if add_path_to_profile "${dir}" "${HOME}/.profile" false; then
        profiles_written="${profiles_written} ~/.profile"
    else
        profiles_skipped="${profiles_skipped} ~/.profile"
    fi

    # fish (only if fish is installed or is the default shell)
    if [ "${shell_name}" = "fish" ] || command -v fish >/dev/null 2>&1; then
        local fish_config="${HOME}/.config/fish/config.fish"
        if add_path_to_profile "${dir}" "${fish_config}" true; then
            profiles_written="${profiles_written} ~/.config/fish/config.fish"
        else
            profiles_skipped="${profiles_skipped} ~/.config/fish/config.fish"
        fi
    fi

    # Determine primary profile for reload instruction
    case "${shell_name}" in
        zsh)    primary_profile="${HOME}/.zshrc" ;;
        bash)   primary_profile="${HOME}/.bashrc" ;;
        fish)   primary_profile="${HOME}/.config/fish/config.fish" ;;
        *)      primary_profile="${HOME}/.profile" ;;
    esac

    PATH_TARGET="${primary_profile}"

    if [ "${shell_name}" = "fish" ]; then
        PATH_LINE="fish_add_path ${dir}"
        PATH_RELOAD="source ${primary_profile}"
    else
        PATH_LINE="export PATH=\"\$PATH:${dir}\""
        PATH_RELOAD=". ${primary_profile}"
    fi

    # Report what was written
    if [ -n "${profiles_written}" ]; then
        ok "Added to PATH in:${profiles_written}"
        PATH_STATUS="added"
    else
        step "PATH already configured in all profiles"
        PATH_STATUS="already present"
    fi

    if [ -n "${profiles_skipped}" ]; then
        step "Already present in:${profiles_skipped}"
    fi

    if [ "${has_session_path}" = true ]; then
        return
    fi

    # Update current session (only effective when script is sourced, not piped)
    export PATH="${PATH}:${dir}"
}

print_install_summary() {
    local installed_version="$1" bin_path="$2"

    echo ""
    step "Install summary"
    printf '    Version: %s\n' "${installed_version}" >&2
    printf '    Binary: %s\n' "${bin_path}" >&2
    printf '    Install dir: %s\n' "$(dirname "${bin_path}")" >&2
    if [ "${NO_PATH}" = true ]; then
        printf '    PATH target: skipped (--no-path)\n' >&2

        return
    fi
    printf '    Shell: %s\n' "${PATH_SHELL}" >&2
    printf '    PATH target: %s (%s)\n' "${PATH_TARGET}" "${PATH_STATUS}" >&2
    printf '    Reload: %s\n' "${PATH_RELOAD}" >&2
}

# ── Resolve install directory ──────────────────────────────────────

resolve_install_dir() {
    local dir="$1"
    if [ -n "${dir}" ]; then
        echo "${dir}"
        return
    fi

    # Use ~/.local/bin if it exists or is standard; fallback to /usr/local/bin
    if [ -d "${HOME}/.local/bin" ] || [ -w "${HOME}/.local" ]; then
        echo "${HOME}/.local/bin"
    elif [ -w "/usr/local/bin" ]; then
        echo "/usr/local/bin"
    else
        echo "${HOME}/.local/bin"
    fi
}

# ── Parse arguments ────────────────────────────────────────────────

parse_args() {
    VERSION=""
    INSTALL_DIR=""
    ARCH_FLAG=""
    NO_PATH=false

    while [ $# -gt 0 ]; do
        case "$1" in
            --version)
                VERSION="$2"
                shift 2
                ;;
            --dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            --arch)
                ARCH_FLAG="$2"
                shift 2
                ;;
            --no-path)
                NO_PATH=true
                shift
                ;;
            --help|-h)
                echo "Usage: install.sh [--version <tag>] [--dir <path>] [--arch <arch>] [--no-path]"
                echo ""
                echo "Options:"
                echo "  --version <tag>   Install a specific version (e.g. v2.55.0)"
                echo "  --dir <path>      Target directory (default: ~/.local/bin)"
                echo "  --arch <arch>     Force architecture: amd64, arm64 (default: auto)"
                echo "  --no-path         Skip adding install directory to PATH"
                exit 0
                ;;
            *)
                err "Unknown option: $1"
                err "Run with --help for usage."
                exit 1
                ;;
        esac
    done
}

# ── Main ───────────────────────────────────────────────────────────

main() {
    echo ""
    echo "  gitmap installer"
    printf '  \033[90mgithub.com/%s\033[0m\n' "${REPO}"
    echo ""

    parse_args "$@"

    local os arch version install_dir archive_path

    os="$(detect_os)"
    arch="$(detect_arch "${ARCH_FLAG}")"
    version="$(resolve_version "${VERSION}")"
    install_dir="$(resolve_install_dir "${INSTALL_DIR}")"

    # Create TMP_DIR in parent scope so install_binary and cleanup can access it.
    TMP_DIR="$(mktemp -d)"
    archive_path="$(download_asset "${version}" "${os}" "${arch}")"

    install_binary "${archive_path}" "${install_dir}" "${os}" "${arch}" "${version}"

    if [ "${NO_PATH}" = false ]; then
        add_to_path "${install_dir}"
    fi

    # Verify the binary works
    local bin_path="${install_dir}/${BINARY_NAME}"
    local installed_version="${version}"
    if [ -f "${bin_path}" ]; then
        echo ""
        local version_output
        if version_output="$("${bin_path}" version 2>&1)"; then
            installed_version="${version_output}"
            ok "gitmap ${version_output}"
        else
            err "Binary found but failed to run."
        fi
    else
        err "Binary not found at ${bin_path}"
    fi

    print_install_summary "${installed_version}" "${bin_path}"
    if [ "${NO_PATH}" = false ]; then
        echo ""
        printf '  \033[32m✓\033[0m  To start using gitmap \033[1mright now\033[0m, run:\n' >&2
        echo "" >&2
        printf '      \033[36m%s\033[0m\n' "${PATH_RELOAD}" >&2
        echo "" >&2
        printf '     Or open a new terminal window.\n' >&2
        echo "" >&2
        printf '  \033[90mInstalled to: %s\033[0m\n' "${install_dir}/${BINARY_NAME}" >&2
        printf '  \033[90mPATH added to: %s %s %s %s\033[0m\n' "~/.zshrc" "~/.zprofile" "~/.bashrc" "~/.profile" >&2
    fi

    echo ""
    ok "Done! Run 'gitmap --help' to get started."
    echo ""
}

main "$@"
