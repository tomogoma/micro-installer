#!/bin/bash

APP_NAME="micro"
APP_FILE="/usr/local/bin/$APP_NAME"
SYSTEMD_FILE="/etc/systemd/system/${APP_NAME}.service"

echo "Uninstalling..."
if [ -f "$APP_FILE" ]; then
	systemctl stop ${APP_NAME}.service >/dev/null
	rm -f "$APP_FILE"
fi
if [ -f "$SYSTEMD_FILE" ]; then
	rm -f "$SYSTEMD_FILE"
fi
systemctl daemon-reload
echo "Uninstall complete!"
exit 0
