#!/bin/sh

set -e

case "${1}" in
	configure|abort-upgrade|abort-deconfigure|abort-remove)
		if [ -x "/usr/bin/deb-systemd-helper" ]; then
			deb-systemd-helper unmask 'rsca.service' >/dev/null || true
			if deb-systemd-helper --quiet was-enabled 'rsca.service'; then
				deb-systemd-helper enable 'rsca.service' >/dev/null || true
			else
				deb-systemd-helper update-state 'rsca.service' >/dev/null || true
			fi
		fi
	;;
esac

exit 0
