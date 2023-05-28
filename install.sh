#!/usr/bin/env bash

echo "Installing Subtitlr..."
echo "------------------------"

# Determining the Linux distribution and architecture
distro=$(lsb_release -i -s)
arch=$(uname -m)

echo "Distribution: $distro"
echo "Architecture: $arch"

# Subtitlr version
version="v0.1.0"

echo "Version: $version"
echo "------------------------"

# URL for downloading the archive based on the distribution and architecture
url=""

case "$distro" in
  "Darwin")
    case "$arch" in
      "x86_64")
        url="https://github.com/yoanbernabeu/Subtitlr/releases/download/${version}/Subtitlr-${version}-darwin-amd64.tar.gz"
        ;;
      "arm64")
        url="https://github.com/yoanbernabeu/Subtitlr/releases/download/${version}/Subtitlr-${version}-darwin-arm64.tar.gz"
        ;;
      *)
        echo "Unsupported architecture"
        exit 1
        ;;
    esac
    ;;
  "Ubuntu"|"Debian"|"Raspbian")
  echo "Downloading Subtitlr..."
    case "$arch" in
      "i686")
        url="https://github.com/yoanbernabeu/Subtitlr/releases/download/${version}/Subtitlr-${version}-linux-386.tar.gz"
        ;;
      "x86_64")
        url="https://github.com/yoanbernabeu/Subtitlr/releases/download/${version}/Subtitlr-${version}-linux-amd64.tar.gz"
        echo $url
        ;;
      "arm64")
        url="https://github.com/yoanbernabeu/Subtitlr/releases/download/${version}/Subtitlr-${version}-linux-arm64.tar.gz"
        ;;
      *)
        echo "Unsupported architecture"
        exit 1
        ;;
    esac
    ;;
  *)
    echo "Unsupported distribution"
    exit 1
    ;;
esac

# Downloading the archive to home directory (and check if url is not 404)
echo "Downloading Subtitlr..."
wget -q --spider $url
if [ $? -eq 0 ]; then
  wget -O ~/Subtitlr.tar.gz $url -q --show-progress
else
  echo "------------------------"
  echo "Subtitlr archive not found"
  echo "------------------------"
  exit 1
fi

# Extracting the archive (if it exists)
echo "Extracting Subtitlr..."
if [ -f ~/Subtitlr.tar.gz ]; then
  tar -xzf ~/Subtitlr.tar.gz -C ~/
else
  echo "Subtitlr archive not found"
  exit 1
fi

# Removing the archive
echo "Removing archive..."
rm ~/Subtitlr.tar.gz

# Moving the binary to /usr/local/bin
echo "Moving Subtitlr to /usr/local/bin..."
sudo mv ~/Subtitlr /usr/local/bin/

# Making the binary executable
echo "Making Subtitlr executable..."
sudo chmod +x /usr/local/bin/Subtitlr

# Sending a message to the user
echo "-----------------------------------------"
echo "Subtitlr successfully installed"
echo "-----------------------------------------"