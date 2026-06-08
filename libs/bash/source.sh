#!/bin/bash

lib_wait_for_vault() {
    local addr="$1"
    local retries="${2:-60}"

    while [ "$retries" -gt 0 ]; do
        if curl -sf "$addr/v1/sys/health" >/dev/null 2>&1; then
            return 0
        fi
        sleep 1
        retries=$((retries - 1))
    done

    return 1
}

lib_login_to_vault() {
    local addr="$1"
    local username="$2"
    local password="$3"

    if ! lib_wait_for_vault "$addr"; then
        return 1
    fi

    local response
    response=$(curl -sf -X POST \
        -H "Content-Type: application/json" \
        -d "{\"password\": \"$password\"}" \
        "$addr/v1/auth/userpass/login/$username") || return 1

    VAULT_TOKEN=$(printf '%s' "$response" | python3 -c 'import json,sys; print(json.load(sys.stdin)["auth"]["client_token"])')
    export VAULT_ADDR="$addr"
    export VAULT_TOKEN

    [ -n "$VAULT_TOKEN" ]
}

lib_create_kv_entries() {
    local kv_path="$1"
    local env_file="$2"

    if [ ! -f "$env_file" ]; then
        echo "Env file not found: $env_file" >&2
        return 1
    fi

    local payload
    payload=$(python3 - "$env_file" <<'PY'
import json
import re
import sys

data = {}
with open(sys.argv[1], encoding="utf-8") as env_file:
    for line in env_file:
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        match = re.match(r"([^=]+)=(.*)$", line)
        if not match:
            continue
        key = match.group(1).strip()
        value = match.group(2).strip()
        if len(value) >= 2 and value[0] == value[-1] and value[0] in "\"'":
            value = value[1:-1]
        data[key] = value

print(json.dumps(data))
PY
) || return 1

    curl -sf -X POST \
        -H "X-Vault-Token: $VAULT_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$payload" \
        "$VAULT_ADDR/v1/$kv_path" >/dev/null

    echo "Wrote secrets to $kv_path from $env_file"
}
