#!/bin/sh

set -e

if [ "$1" = "remove" ]; then
	if [ -x "/usr/bin/deb-systemd-helper" ]; then
		deb-systemd-helper mask 'rscad.service' >/dev/null || true
	fi
fi

if [ "$1" = "purge" ]; then
	if [ -x "/usr/bin/deb-systemd-helper" ]; then
		deb-systemd-helper purge 'rscad.service' >/dev/null || true
		deb-systemd-helper unmask 'rscad.service' >/dev/null || true
	fi
fi

exit 0
