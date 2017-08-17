#!/bin/bash
set -e
# Install on linux only
# Download latest awless binary from Github

DOWNLOAD_URL=$(curl https://api.github.com/repos/wallix/awless-scheduler/releases/latest | grep browser_download_url | sed 's/"browser_download_url": //g' | sed 's/[ ""]//g')

echo "Downloading scheduler from $DOWNLOAD_URL"

TARFILE="awless-scheduler.tar.gz"

if ! curl --fail -o $TARFILE -L $DOWNLOAD_URL; then
    exit
fi

mkdir -p /var/awless-scheduler/bin

tar -C /var/awless-scheduler/bin -xzf $TARFILE 


cat <<'EOF' >> /etc/systemd/system/awless-scheduler.service
[Unit]
Description=Awless Scheduler Daemon
After=network.target

[Service]
ExecStart=/var/awless-scheduler/bin/awless-scheduler -http-mode
Restart=always

[Install]
WantedBy=multi-user.target
EOF

service awless-scheduler restart