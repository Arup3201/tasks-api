#!/usr/bin/env bash
# import_env.sh
# Usage: source ./import_env.sh [file]
# Default file: .env

ENV_FILE="${1:-.env}"

if [ ! -f "$ENV_FILE" ]; then
  echo "ERROR: $ENV_FILE not found in $(pwd)"
  return 1 2>/dev/null || exit 1
fi

# Read .env line by line
while IFS= read -r line || [ -n "$line" ]; do
  # Trim leading/trailing whitespace
  line="$(echo "$line" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')"

  # Skip empty lines and comments
  [ -z "$line" ] && continue
  case "$line" in
    \#*) continue ;;
  esac

  # Remove optional "export " prefix
  line="${line#export }"

  # Split KEY=VALUE
  if [[ "$line" == *"="* ]]; then
    key="${line%%=*}"
    value="${line#*=}"

    # Remove surrounding quotes if any
    value="${value%\"}"
    value="${value#\"}"
    value="${value%\'}"
    value="${value#\'}"

    export "$key=$value"
    echo "Exported $key"
  fi
done < "$ENV_FILE"

echo "Done. Variables loaded from $ENV_FILE"