#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

SSH_KEY_PATH="${1:-}"
if [[ -z "$SSH_KEY_PATH" ]]; then
	echo "usage: deploy_backend.sh /path/to/ssh_key"
	exit 1
fi

REMOTE_USER="grimlock"
REMOTE_HOST="web1"
REMOTE_TARGET_DIR="/var/www/vhosts/apinote-thing.truvis.co"
REMOTE_CONFIG_DIR="/etc/note-thing"

BUILD_DIR="${ROOT_DIR}/dist"
mkdir -p "${BUILD_DIR}"

echo "Building linux binary..."
cd "${ROOT_DIR}"
GOOS=linux GOARCH=amd64 go build -o "${BUILD_DIR}/note-thing" ./cmd/server

echo "Uploading binary..."
rsync -avz -e "ssh -i ${SSH_KEY_PATH}" \
	"${BUILD_DIR}/note-thing" \
	"${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_TARGET_DIR}/note-thing"

echo "Uploading config sample..."
rsync -avz -e "ssh -i ${SSH_KEY_PATH}" \
	"${ROOT_DIR}/config.ini.sample" \
	"${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_CONFIG_DIR}/config.ini.sample"

echo "Ensuring config.ini exists on server..."
ssh -i "${SSH_KEY_PATH}" "${REMOTE_USER}@${REMOTE_HOST}" \
	"sudo mkdir -p ${REMOTE_CONFIG_DIR} && \
	 if [ ! -f ${REMOTE_CONFIG_DIR}/config.ini ]; then \
	  sudo cp ${REMOTE_CONFIG_DIR}/config.ini.sample ${REMOTE_CONFIG_DIR}/config.ini; \
	 fi"

echo "Installing systemd unit..."
scp -i "${SSH_KEY_PATH}" \
	"${ROOT_DIR}/deploy/note-thing.service" \
	"${REMOTE_USER}@${REMOTE_HOST}:/tmp/note-thing.service"

ssh -i "${SSH_KEY_PATH}" "${REMOTE_USER}@${REMOTE_HOST}" \
	"sudo mv /tmp/note-thing.service /etc/systemd/system/note-thing.service && \
	 sudo systemctl daemon-reload && \
	 sudo systemctl enable note-thing.service && \
	 sudo systemctl restart note-thing.service"

echo "Backend deployed."
