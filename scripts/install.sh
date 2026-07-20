#!/bin/sh
# scripts/install.sh — installs the latest mcli release for this machine.
# Usage: curl -fsSL https://raw.githubusercontent.com/skygrime35/mcli/main/scripts/install.sh | sh
set -eu

REPO="skygrime35/mcli"

detect_os() {
	os=$(uname -s)
	case "$os" in
		Linux) echo "linux" ;;
		*)
			echo "mcli only supports Linux and Android/Termux (detected: $os)" >&2
			exit 1
			;;
	esac
}

detect_arch() {
	arch=$(uname -m)
	case "$arch" in
		x86_64|amd64) echo "amd64" ;;
		aarch64|arm64) echo "arm64" ;;
		*)
			echo "unsupported architecture: $arch (mcli releases only cover amd64/arm64)" >&2
			exit 1
			;;
	esac
}

is_termux() {
	# Termux sets $PREFIX to a path containing "com.termux".
	case "${PREFIX:-}" in
		*com.termux*) return 0 ;;
		*) return 1 ;;
	esac
}

install_dir() {
	if is_termux; then
		echo "${PREFIX}/bin"
	else
		echo "${HOME}/.local/bin"
	fi
}

latest_tag() {
	curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
		| grep '"tag_name":' \
		| sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
}

main() {
	os=$(detect_os)
	arch=$(detect_arch)
	tag=$(latest_tag)
	if [ -z "$tag" ]; then
		echo "could not determine the latest mcli release version" >&2
		exit 1
	fi
	version=${tag#v}

	asset="mcli_${version}_${os}_${arch}.tar.gz"
	url="https://github.com/${REPO}/releases/download/${tag}/${asset}"

	dest_dir=$(install_dir)
	mkdir -p "$dest_dir"

	tmp_dir=$(mktemp -d)
	trap 'rm -rf "$tmp_dir"' EXIT

	echo "Downloading mcli ${tag} for ${os}/${arch}..."
	curl -fsSL "$url" -o "${tmp_dir}/mcli.tar.gz"
	tar -xzf "${tmp_dir}/mcli.tar.gz" -C "$tmp_dir" mcli
	mv "${tmp_dir}/mcli" "${dest_dir}/mcli"
	chmod +x "${dest_dir}/mcli"

	echo "Installed mcli ${tag} to ${dest_dir}/mcli"
	case ":$PATH:" in
		*":${dest_dir}:"*) ;;
		*)
			echo ""
			echo "NOTE: ${dest_dir} is not in your PATH. Add this to your shell profile:"
			echo "  export PATH=\"${dest_dir}:\$PATH\""
			;;
	esac
}

# Guarded rather than an unconditional call: this lets the functions above be
# safely sourced (e.g. `. ./scripts/install.sh`) for testing the detection
# logic (detect_os, detect_arch, is_termux, install_dir) without triggering a
# real download. Set MCLI_INSTALL_SKIP_MAIN=1 before sourcing to opt out of
# running main. Normal usage (./install.sh, `sh install.sh`, or
# `curl -fsSL ... | sh`) is unaffected since the variable is unset by default.
if [ -z "${MCLI_INSTALL_SKIP_MAIN:-}" ]; then
	main "$@"
fi
