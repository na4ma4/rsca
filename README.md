# rsca

![ci](https://github.com/na4ma4/rsca/workflows/ci/badge.svg)
[![GoDoc](https://godoc.org/github.com/na4ma4/rsca/?status.svg)](https://godoc.org/github.com/na4ma4/rsca)
[![GitHub issues](https://img.shields.io/github/issues/na4ma4/rsca)](https://github.com/na4ma4/rsca/issues)
[![GitHub forks](https://img.shields.io/github/forks/na4ma4/rsca)](https://github.com/na4ma4/rsca/network)
[![GitHub stars](https://img.shields.io/github/stars/na4ma4/rsca)](https://github.com/na4ma4/rsca/stargazers)
[![GitHub license](https://img.shields.io/github/license/na4ma4/rsca)](https://github.com/na4ma4/rsca/blob/main/LICENSE)

Remote Service Check Acceptor (alternative to NSCA)

## Design

Three components:

- `rsc`: admin client.
- `rsca`: check client, executes checks on servers.
- `rscad`: check daemon, receives check results and writes them to the nagios command file.

## Deployment

### rsc tool

This can be run on any host with network access to the `rscad` server.

### rsca service

This should be run on the server to check, it runs the checks in the config file and sends them to `rscad` on schedule.

Copy [rsca.service](systemd/client/rsca.service) to `/etc/systemd/system/rsca.service`.
Copy [rsca.toml](testdata/rsca.toml) to `/etc/nagios/rsca.toml`.

### rscad service

This should be run on the nagios server, it handles the connections from the `rsca` clients.

Copy [rscad.service](systemd/server/rscad.service) to `/etc/systemd/system/rscad.service`.

## Support

Reach out to the maintainer at one of the following places:

[GitHub discussions](https://github.com/na4ma4/rsca/discussions)
[GitHub issues](https://github.com/na4ma4/rsca/issues)
The email which is located in [GitHub profile](https://github.com/na4ma4).
