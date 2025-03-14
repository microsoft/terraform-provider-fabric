#!/usr/bin/env bash

set -e
export DEBIAN_FRONTEND="noninteractive"

if [ -n "$SUDO_USER" ]; then
  HOME=$(getent passwd "$SUDO_USER" | cut -d: -f6)
fi

# Variables
# shellcheck disable=SC2153
PROVIDER_NAME=${PROVIDERNAME}
if [ -z "$PROVIDER_NAME" ]; then
  echo "The 'providerName' is not set"
  exit 1
fi
echo "PROVIDER_NAME is: $PROVIDER_NAME"

WORKSPACE_PATH=${WORKSPACE:-$PWD}
echo "WORKSPACE_PATH is: $WORKSPACE_PATH"

USER_HOME=${_REMOTE_USER_HOME:-$HOME}
echo "USER_HOME is: $USER_HOME"

SCRIPT_PATH="$(dirname "$(realpath "$0")")"
echo "SCRIPT_PATH is: $SCRIPT_PATH"

USER="$_REMOTE_USER"
if [ -n "$WSL_INTEROP" ] || [ -n "$WSL_DISTRO_NAME" ]; then
  USER=${SUDO_USER:-$(whoami)}
fi
echo "USER is: $USER"

set_ownership_to_current_user() {
  target_path=$1
  if [ -d "$target_path" ]; then
    chown -R "$USER" "$target_path"
  fi
}

set_terraformrc() {
  target_path=$1

  os_type=$(uname)
  os_type=${os_type,,}

  os_arch=$(uname -m)
  if echo "$os_arch" | grep -E -q 'x86_64|amd64'; then
    os_arch="amd64"
  elif echo "$os_arch" | grep -E -q 'i386|i686'; then
    os_arch="386"
  elif echo "$os_arch" | grep -E -q 'aarch64|arm64'; then
    os_arch="arm64"
  elif echo "$os_arch" | grep -E -q 'armv'; then
    os_arch="arm"
  else
    echo "Unsupported architecture"
    exit 1
  fi

  sed <"$SCRIPT_PATH/terraformrc.tmpl" -e "s|{{PROVIDER_NAME}}|${PROVIDER_NAME}|g" -e "s|{{PROVIDER_BIN_PATH}}|${WORKSPACE_PATH}/bin/${os_type}-${os_arch}|g" >"$target_path"
}

install() {
  TFRC_PATH="$USER_HOME/.terraformrc"
  set_terraformrc "$TFRC_PATH"
  set_ownership_to_current_user "$TFRC_PATH"
  set_ownership_to_current_user "$WORKSPACE_PATH"
  set_ownership_to_current_user "$(go env GOPATH)"
}

install
