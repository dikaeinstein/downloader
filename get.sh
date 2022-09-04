#!/bin/sh
ARCH=="amd64"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
VERSION=$1
# $VERSION with v prefix stripped
VERSION_STRIPPED="${VERSION:1}"

echo "Downloading downloadctl release ${VERSION} for ${OS}_${ARCH} ..."
echo ""

RELEASE_URL="https://github.com/dikaeinstein/downloader/releases/download/${VERSION}/downloadctl_${VERSION_STRIPPED}_${OS}_${ARCH}.tar.gz"
code=$(curl -w '%{http_code}' -L $RELEASE_URL -o /tmp/downloadctl_.tar.gz)

if [ $code != "200" ]; then
  echo ""
  echo "[error] Failed to download downloadctl release ${VERSION} for $OS $ARCH."
  echo "Received HTTP status code $code"
  echo ""
  echo "Supported versions of the downloadctl are:"
  echo " - darwin_amd64"
  echo " - linux_amd64"
  echo ""
  exit 1
fi

tar -xzf /tmp/downloadctl.tar.gz -C /tmp
sudo chmod +x /tmp/downloadctl
sudo mv /tmp/downloadctl /usr/local/bin/

echo ""
echo "downloadctl ${VERSION} for ${OS}_${ARCH} installed."
