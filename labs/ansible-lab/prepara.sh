#!/usr/bin/env bash
# Crea la coppia di chiavi SSH del laboratorio (solo la prima volta).
set -euo pipefail
cd "$(dirname "$0")"
mkdir -p .chiavi
chmod 700 .chiavi
if [[ ! -f .chiavi/id_ed25519 ]]; then
  ssh-keygen -q -t ed25519 -N '' -C ansible-lab -f .chiavi/id_ed25519
  echo "Chiavi create in .chiavi/"
else
  echo "Chiavi già presenti in .chiavi/"
fi
