#!/usr/bin/env bash
# Eseguito una sola volta, alla creazione del Dev Container.
set -euo pipefail

python3 -m pip install --user --upgrade pip
python3 -m pip install --user -r requirements.txt   # per vedere il sito con "mkdocs serve"
python3 -m pip install --user ansible-core flask pytest requests pyyaml

echo
echo "Ambiente pronto. Controllo delle versioni:"
for cmd in git python3 docker aws gcloud terraform kubectl helm kind ansible; do
  printf '  %-10s ' "$cmd"
  case "$cmd" in
    kubectl) kubectl version --client 2>/dev/null | head -1 ;;
    gcloud)  gcloud version 2>/dev/null | head -1 ;;
    *)       "$cmd" --version 2>/dev/null | head -1 ;;
  esac
done
