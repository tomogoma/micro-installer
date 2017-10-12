#!/bin/bash

source vars.sh

for f in ${UNITS}
do
    sysDF="${SYSTEMD_DIR}/$(basename ${f})"
    if [ -f "$sysDF" ]; then
        systemctl stop ${f} >/dev/null
        rm -f ${sysDF} || exit 1
    fi
done

systemctl daemon-reload

if [ -f "$INSTALL_FILE" ]; then
	rm -f "$INSTALL_FILE" || exit 1
fi

