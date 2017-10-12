#!/bin/bash

source vars.sh

function installService {
    for f in ${UNITS}
    do
	    cp -f "${f}" "$SYSTEMD_DIR" || exit 1
	    systemctl enable "$(basename ${f})" || exit 1
    done
}

cp -f "${BUILD_EXEC}" "${INSTALL_FILE}" || exit 1
installService
