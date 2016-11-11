#!/bin/bash

APP_NAME="micro"
APP_DIR="/usr/local/bin"
SYSTEMD_DIR="/etc/systemd/system"

function usage {
	printf "Micro api Installer\nUsage: $1 [user]\n"
}

function install {
    cp -f "${APP_NAME}" "$APP_DIR/$APP_NAME" || exit 1
    usr=$1
    if [ -z "$usr" ];then
        return
    fi
    chown -R "$usr:$usr" "$APP_DIR/$APP_NAME"
}

function installService {
	mkdir -p "$SYSTEMD_DIR" || exit 1
	cp -f "${APP_NAME}.service" "$SYSTEMD_DIR" || exit 1
	systemctl enable "${APP_NAME}.service" || exit 1
}

## Begin processing script
./systemdUninstaller.sh
echo "Installing..."
install $1
installService
echo "Done!"
exit 0
