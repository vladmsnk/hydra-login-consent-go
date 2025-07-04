#!/usr/bin/env bash
set -euo pipefail

SERVICE_NAME=hydra-login-consent
INSTALL_DIR=/opt/$SERVICE_NAME
BIN_PATH=$INSTALL_DIR/login_consent
SERVICE_FILE=/etc/systemd/system/$SERVICE_NAME.service
ENV_FILE=/etc/$SERVICE_NAME.env

# Build the binary
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
echo "Building binary for $GOOS/$GOARCH"
go build -o login_consent ./app

sudo mkdir -p "$INSTALL_DIR"

# Copy files
sudo mv login_consent "$BIN_PATH"
sudo cp -r ui "$INSTALL_DIR/"
sudo cp -r config "$INSTALL_DIR/"

# Install service file
sudo install -m 644 hydra-login-consent.service "$SERVICE_FILE"

# Create environment file if missing
if [ ! -f "$ENV_FILE" ]; then
    sudo tee "$ENV_FILE" >/dev/null <<EOT
HYDRA_ADMIN_URL=https://oauthidm.ru
HYDRA_USERNAME=adminuser
HYDRA_PASSWORD=1234
HOST=0.0.0.0
PORT=3000
EOT
fi

sudo systemctl daemon-reload
sudo systemctl enable "$SERVICE_NAME"
sudo systemctl restart "$SERVICE_NAME"

echo "Service $SERVICE_NAME installed and started"
