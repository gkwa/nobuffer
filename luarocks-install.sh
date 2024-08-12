#!/bin/bash

# Exit on error
set -e

# Update package list and upgrade existing packages
apk update
apk upgrade

# Install prerequisites
apk add --no-cache \
    wget \
    tar \
    gcc \
    libc-dev \
    make \
    openssl-dev \
    readline-dev \
    lua5.4 \
    lua5.4-dev

# Download and install LuaRocks
LUAROCKS_VERSION="3.11.1"
wget "https://luarocks.org/releases/luarocks-$LUAROCKS_VERSION.tar.gz"
tar zxpf "luarocks-$LUAROCKS_VERSION.tar.gz"
cd "luarocks-$LUAROCKS_VERSION"

./configure \
    --prefix=/usr \
    --with-lua-include=/usr/include/lua5.4 \
    --with-lua=/usr \
    --with-lua-interpreter=lua5.4

make
make install

# Clean up
cd ..
rm -rf "luarocks-$LUAROCKS_VERSION" "luarocks-$LUAROCKS_VERSION.tar.gz"

# Install LuaSocket as an example
luarocks install luasocket

echo "LuaRocks installation completed successfully!"
echo "You can now use 'luarocks' to install Lua packages."
echo "Example usage: luarocks install <package_name>"
