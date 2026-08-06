#!/bin/bash

set -euo pipefail

# Set the userpass credentials for loading env vars
export VAULT_ADDR_1="http://127.0.0.1:8200"
export VAULT_ADDR_2="http://vault:8200"
VAULT_USERNAME="service"
VAULT_PASSWORD="gruent"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Source the bash libraries
source "$SCRIPT_DIR/../../../libs/bash/source.sh"

# Vault KV paths and environment files
API_VAULT_KV_PATH="kv/config/http-api"
API_ENV_FILE="$SCRIPT_DIR/.env.api"
CONTROLLER_VAULT_KV_PATH="kv/config/controller"
CONTROLLER_ENV_FILE="$SCRIPT_DIR/.env.controller"

function main() {
    local logged_in=false

    # Try each Vault address until one works
    for addr in "$VAULT_ADDR_1" "$VAULT_ADDR_2"; do
        echo "Trying to connect to Vault at: $addr"
        if lib_login_to_vault "$addr" "$VAULT_USERNAME" "$VAULT_PASSWORD"; then
            logged_in=true
            break
        fi
        echo "Failed to connect to $addr"
    done

    if [ "$logged_in" != true ]; then
        echo "Could not log in to Vault" >&2
        exit 1
    fi

    lib_create_kv_entries "$API_VAULT_KV_PATH" "$API_ENV_FILE"
    lib_create_kv_entries "$CONTROLLER_VAULT_KV_PATH" "$CONTROLLER_ENV_FILE"
}

main