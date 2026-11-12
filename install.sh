#!/bin/sh
# Install rhiza: the newest release (or --channel latest) for this platform.
#
#   curl -fsSL https://raw.githubusercontent.com/dphaener/rhiza/main/install.sh | sh
#   curl -fsSL .../install.sh | sh -s -- --channel latest --dir /usr/local/bin
#
# Options:
#   --channel release|latest   release (default) is the newest v* tag; latest
#                              is the newest build of main.
#   --dir DIR                  install into DIR (default ~/.local/bin).
#   -h, --help                 show this help.
#
# The archive's sha256 is checked against the channel manifest; on a mismatch
# nothing is installed. RHIZA_INSTALL_BASE_URL replaces https://github.com/dphaener/rhiza
# (for mirrors and tests).
set -eu

base=${RHIZA_INSTALL_BASE_URL:-https://github.com/dphaener/rhiza}
channel=release
dir=${HOME:-}/.local/bin

die() { echo "install.sh: $*" >&2; exit 1; }

usage() { sed -n '2,/^set -eu/p' "$0" | sed '$d' | sed 's/^# \{0,1\}//'; }

while [ $# -gt 0 ]; do
	case $1 in
	--channel) [ $# -ge 2 ] || die "--channel needs a value"; channel=$2; shift 2 ;;
	--channel=*) channel=${1#*=}; shift ;;
	--dir) [ $# -ge 2 ] || die "--dir needs a value"; dir=$2; shift 2 ;;
	--dir=*) dir=${1#*=}; shift ;;
	-h|--help) usage; exit 0 ;;
	*) die "unknown option $1 (try --help)" ;;
	esac
done

case $channel in
release) manifest_url=$base/releases/latest/download/manifest.json ;;
latest) manifest_url=$base/releases/download/latest/manifest.json ;;
*) die "unknown channel $channel (want release or latest)" ;;
esac
[ -n "$dir" ] && [ "$dir" != /.local/bin ] || die "no install directory: set HOME or pass --dir"

case $(uname -s) in
Linux) os=linux ;;
Darwin) os=darwin ;;
*) die "unsupported OS $(uname -s): rhiza is built for linux and darwin" ;;
esac
case $(uname -m) in
x86_64|amd64) arch=amd64 ;;
aarch64|arm64) arch=arm64 ;;
*) die "unsupported architecture $(uname -m): rhiza is built for amd64 and arm64" ;;
esac

command -v curl >/dev/null 2>&1 || die "curl is required"
command -v tar >/dev/null 2>&1 || die "tar is required"
if command -v sha256sum >/dev/null 2>&1; then
	sha256() { sha256sum "$1" | cut -d' ' -f1; }
elif command -v shasum >/dev/null 2>&1; then
	sha256() { shasum -a 256 "$1" | cut -d' ' -f1; }
else
	die "sha256sum or shasum is required"
fi

tmp=$(mktemp -d "${TMPDIR:-/tmp}/rhiza-install.XXXXXX") || die "cannot make a temporary directory"
trap 'rm -rf "$tmp"' EXIT INT TERM

echo "Fetching the $channel manifest"
curl -fsSL "$manifest_url" -o "$tmp/manifest.json" || die "cannot fetch $manifest_url"

# The manifest is JSON with one field per line (tools/release writes it so);
# this reads the os/arch platform entry's url and sha256 without jq.
entry=$(sed 's/[{},]/\n/g' "$tmp/manifest.json" | awk -v os="$os" -v arch="$arch" '
	function val(l) { sub(/^[^:]*:[ \t]*"/, "", l); sub(/"[ \t]*$/, "", l); return l }
	/"os"[ \t]*:/ { o = val($0); a = ""; u = ""; next }
	/"arch"[ \t]*:/ { a = val($0); next }
	/"url"[ \t]*:/ { u = val($0); next }
	/"sha256"[ \t]*:/ { if (o == os && a == arch && u != "") { print u " " val($0); exit } }
')
[ -n "$entry" ] || die "the $channel manifest has no build for $os/$arch"
url=${entry% *}
want=${entry#* }
version=$(sed 's/[{},]/\n/g' "$tmp/manifest.json" | awk '/"version"[ \t]*:/ { sub(/^[^:]*:[ \t]*"/, ""); sub(/"[ \t]*$/, ""); print; exit }')

echo "Downloading rhiza ${version:-} for $os/$arch"
curl -fsSL "$url" -o "$tmp/rhiza.tar.gz" || die "cannot fetch $url"
got=$(sha256 "$tmp/rhiza.tar.gz")
[ "$got" = "$want" ] || die "sha256 mismatch for $url: manifest says $want, got $got; nothing installed"

mkdir "$tmp/x"
tar -xzf "$tmp/rhiza.tar.gz" -C "$tmp/x" rhiza || die "the archive has no rhiza binary"
[ -f "$tmp/x/rhiza" ] || die "the archive has no rhiza binary"

mkdir -p "$dir" || die "cannot create $dir"
chmod 755 "$tmp/x/rhiza"
# Move within the destination so the final rename is atomic.
mv "$tmp/x/rhiza" "$dir/.rhiza.new.$$" || die "cannot write to $dir"
mv -f "$dir/.rhiza.new.$$" "$dir/rhiza" || { rm -f "$dir/.rhiza.new.$$"; die "cannot write $dir/rhiza"; }

echo "Installed $dir/rhiza"
case ":${PATH:-}:" in
*":$dir:"*) ;;
*) echo "Add $dir to your PATH to run rhiza." ;;
esac
