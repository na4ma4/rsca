# rsca

![ci](https://github.com/na4ma4/rsca/workflows/ci/badge.svg)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fna4ma4%2Frsca.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fna4ma4%2Frsca?ref=badge_shield)

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
Copy [rsca.toml](test/rsca.toml) to `/etc/nagios/rsca.toml`.

### rscad service

This should be run on the nagios server, it handles the connections from the `rsca` clients.

Copy [rscad.service](systemd/server/rscad.service) to `/etc/systemd/system/rscad.service`.


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fna4ma4%2Frsca.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fna4ma4%2Frsca?ref=badge_large)