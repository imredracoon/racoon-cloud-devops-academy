# Base ufficiale Microsoft per i Dev Container (Ubuntu 24.04 LTS)
FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

# Strumenti di uso quotidiano usati nelle lezioni
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get install -y --no-install-recommends \
       bash-completion dnsutils iproute2 iputils-ping traceroute netcat-openbsd tcpdump \
       jq tree less nano vim shellcheck sqlite3 postgresql-client apt-transport-https \
       ca-certificates gnupg curl \
    && rm -rf /var/lib/apt/lists/*

# Google Cloud CLI dal repository ufficiale Google
# https://cloud.google.com/sdk/docs/install#deb
RUN curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg \
       | gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg \
    && echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" \
       > /etc/apt/sources.list.d/google-cloud-sdk.list \
    && apt-get update && apt-get install -y --no-install-recommends google-cloud-cli \
    && rm -rf /var/lib/apt/lists/*

# kind (Kubernetes in Docker) dal rilascio ufficiale
# https://kind.sigs.k8s.io/docs/user/quick-start/#installing-from-release-binaries
ARG KIND_VERSION=v0.30.0
RUN ARCH=$(dpkg --print-architecture) \
    && curl -fsSLo /usr/local/bin/kind "https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-linux-${ARCH}" \
    && chmod +x /usr/local/bin/kind
