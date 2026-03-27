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
REMOTE_TARGET_DIR="/var/www/vhosts/api-note-thing.truvis.co"
REMOTE_CONFIG_DIR="/etc/api-note-thing"
SERVICE_NAME="api-note-thing"

BUILD_DIR="${ROOT_DIR}/dist"
mkdir -p "${BUILD_DIR}"

echo "Building linux binary..."
cd "${ROOT_DIR}"
GOOS=linux GOARCH=amd64 go build -o "${BUILD_DIR}/note-thing" .

echo "Generating config.ini from environment..."
cat > "${BUILD_DIR}/config.ini" <<EOF
DATABASE_URL=${DATABASE_URL:-}
GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID:-}
GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET:-}
GOOGLE_REDIRECT_URL=${GOOGLE_REDIRECT_URL:-}
JWT_SECRET=${JWT_SECRET:-}
PORT=${PORT:-18611}
FRONTEND_URL=${FRONTEND_URL:-}
STRIPE_SECRET_KEY=${STRIPE_SECRET_KEY:-}
STRIPE_WEBHOOK_SECRET=${STRIPE_WEBHOOK_SECRET:-}
ADMIN_EMAILS=${ADMIN_EMAILS:-}
STRIPE_PRICE_MIGRATION_GRACE_DAYS=${STRIPE_PRICE_MIGRATION_GRACE_DAYS:-0}
EOF

echo "Ensuring target directories exist on server..."
ssh -i "${SSH_KEY_PATH}" "${REMOTE_USER}@${REMOTE_HOST}" \
	"sudo mkdir -p ${REMOTE_CONFIG_DIR} && mkdir -p ${REMOTE_TARGET_DIR}"

echo "Uploading binary..."
rsync -avz -e "ssh -i ${SSH_KEY_PATH}" \
	"${BUILD_DIR}/note-thing" \
	"${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_TARGET_DIR}/note-thing"

echo "Uploading config.ini..."
scp -i "${SSH_KEY_PATH}" \
	"${BUILD_DIR}/config.ini" \
	"${REMOTE_USER}@${REMOTE_HOST}:/tmp/config.ini"

ssh -i "${SSH_KEY_PATH}" "${REMOTE_USER}@${REMOTE_HOST}" \
	"sudo mv /tmp/config.ini ${REMOTE_CONFIG_DIR}/config.ini && \
	 sudo chown root:root ${REMOTE_CONFIG_DIR}/config.ini && \
	 sudo chmod 600 ${REMOTE_CONFIG_DIR}/config.ini"

echo "Uploading config sample..."
rsync -avz -e "ssh -i ${SSH_KEY_PATH}" \
	"${ROOT_DIR}/config.ini.sample" \
	"${REMOTE_USER}@${REMOTE_HOST}:/tmp/config.ini.sample"

ssh -i "${SSH_KEY_PATH}" "${REMOTE_USER}@${REMOTE_HOST}" \
	"sudo mv /tmp/config.ini.sample ${REMOTE_CONFIG_DIR}/config.ini.sample && sudo chown root:root ${REMOTE_CONFIG_DIR}/config.ini.sample"

echo "Installing systemd unit..."
scp -i "${SSH_KEY_PATH}" \
	"${ROOT_DIR}/deploy/${SERVICE_NAME}.service" \
	"${REMOTE_USER}@${REMOTE_HOST}:/tmp/${SERVICE_NAME}.service"

ssh -i "${SSH_KEY_PATH}" "${REMOTE_USER}@${REMOTE_HOST}" \
	"sudo mv /tmp/${SERVICE_NAME}.service /etc/systemd/system/${SERVICE_NAME}.service && \
	 sudo systemctl daemon-reload && \
	 sudo systemctl enable ${SERVICE_NAME}.service && \
	 sudo systemctl restart ${SERVICE_NAME}.service"

# Clean up local build artifacts
rm -f "${BUILD_DIR}/config.ini"

echo "Backend deployed."
