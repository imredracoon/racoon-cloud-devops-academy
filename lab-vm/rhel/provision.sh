#!/usr/bin/env bash
# Preparazione della VM rhel-lab (eseguito da Vagrant come root al primo avvio).
set -euo pipefail

dnf install -y --setopt=strict=0 bash-completion man-pages tree jq bind-utils traceroute tcpdump \
  lvm2 xfsprogs acl git curl nano vim-enhanced policycoreutils-python-utils setroubleshoot-server

# Tre dischi "finti" da 1 GB per gli esercizi su partizioni e LVM (come nella VM Ubuntu).
mkdir -p /srv/lab-disks
for n in 1 2 3; do
  [ -f /srv/lab-disks/disk$n.img ] || truncate -s 1G /srv/lab-disks/disk$n.img
done
cat > /etc/systemd/system/lab-disks.service <<'UNIT'
[Unit]
Description=Collega i dischi di laboratorio come loop device
DefaultDependencies=no
Before=local-fs-pre.target

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/bin/sh -c 'for n in 1 2 3; do losetup -fP /srv/lab-disks/disk$n.img; done'

[Install]
WantedBy=local-fs-pre.target
UNIT
systemctl daemon-reload
systemctl enable --now lab-disks.service
echo "VM rhel-lab pronta: SELinux $(getenforce), firewalld $(systemctl is-active firewalld)"
