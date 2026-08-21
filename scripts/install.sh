#!/bin/sh
set -e

repo="jeonjw85/Ksecret"
bin="ksecret"
dest="${1:-/usr/local/bin}"

os=$(uname -s | tr "[:upper:]" "[:lower:]")
arch=$(uname -m)
case "$arch" in
x86_64) arch=amd64 ;;
aarch64 | arm64) arch=arm64 ;;
esac

url="https://github.com/${repo}/releases/latest/download/${bin}_${os}_${arch}.tar.gz"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

curl -sSfL "$url" | tar -xz -C "$tmp"
if ! mkdir -p "$dest" 2>/dev/null || [ ! -w "$dest" ]; then
	echo "쓰기 권한이 없습니다: $dest" >&2
	echo "예: curl -sSfL https://raw.githubusercontent.com/${repo}/main/scripts/install.sh | sh -s -- \"\$HOME/.local/bin\"" >&2
	exit 1
fi
mv "$tmp/$bin" "$dest/$bin"
chmod +x "$dest/$bin"
echo "설치완료 : $dest/$bin"
