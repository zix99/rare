#!/bin/bash
set -e

REPO=zix99/rare
LATEST_TAG_URL="https://api.github.com/repos/${REPO}/releases/latest"
RELEASES_URL="https://github.com/${REPO}/releases"
FILE_BASENAME="rare"

# Get latest version tag
echo "Fetching latest version..."
LATEST=$( curl -sf "$LATEST_TAG_URL" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' )

if [[ -z $LATEST ]]; then
    echo "Unable to get version" >&2
    exit 1
fi
echo "  Version: ${LATEST}"

# Temp dir and cleanup
TMP_DIR="$(mktemp -d)"
trap "rm -rf \"$TMP_DIR\"" EXIT INT TERM

# Detect version
OS="$(uname -s)"
ARCH="$(uname -m)"

TAR_SUFFIX="tar.gz"
INSTALL_DIR="$HOME/.local/bin"
case $OS in
    aarch64)
        OS=arm64
        ;;
    MSYS_NT*)
        OS=Windows
        TAR_SUFFIX="zip"
        ;;
    Darwin)
        INSTALL_DIR="$HOME/bin"
        ;;
esac

# Download
TAR_FILENAME="${FILE_BASENAME}_${LATEST}_${OS}_${ARCH}.${TAR_SUFFIX}"
URL="$RELEASES_URL/download/$LATEST/$TAR_FILENAME"

cd $TMP_DIR
echo "Downloading $TAR_FILENAME to $TMP_DIR..." >&2
curl -sfLO $URL
tar xzf $TAR_FILENAME

if [[ ! -f rare ]]; then
    echo "Something went wrong, download did not include binary" >&2
    exit 1
fi

# Install
if [[ $USER == "root" ]]; then
    INSTALL_DIR="/usr/bin"
    echo "Installing as root to $INSTALL_DIR..." >&2
    chown root:root rare*
else
    echo "Installing to user home $INSTALL_DIR ..." >&2
    mkdir -p $INSTALL_DIR
fi

if [[ ! $PATH =~ $INSTALL_DIR ]]; then
    echo "WARNING: Installation path not in \$PATH environment" >&2
fi

chmod +x rare
mv rare $INSTALL_DIR
echo "Installed rare to $INSTALL_DIR" >&2

if [[ -f rare-pcre ]]; then
    chmod +x rare-pcre
    mv rare-pcre $INSTALL_DIR
    echo "Installed rare-pcre to $INSTALL_DIR" >&2
fi

# Install man page if able
if [[ $USER == "root" && -f rare.1.gz ]]; then
    cp rare.1.gz /usr/share/man/man1/
    echo "Installed man page" >&2
fi
