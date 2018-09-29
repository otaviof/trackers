# `trackers`

A tool to keep track of Torrent tracker services.

## Motivation

Tracker services is something that can disappear and re-appear pretty much all over the sudden, it
can be subject of blockades, take down actions and other infrastructure related problems. So, at
the end of the day you want to know which trackers are working, and with this information you can
optimize your client, to only spend time on working trackers.

For this book-keeping work `trackers` helps you to understand which ones are in use, test if they
are working, overwrite, cache and update your running torrents.

## Installation

``` bash
go get -u github.com/otaviof/trackers/cmd/trackers
```

## Usage

The following sub-commands are available in `trackers`, via `trackers <sub-command>` in comment
line. All commands and sub-command support a `--help` option, so please consider.

In the top level, `trackers` supports only `--config` option, please consider
[Configuration](#configuration) section.

### Harvest

``` bash
trackers harvest --dry-run
```

Inspect Torrent-Client to read all trackers in use, and save new entries in Trackers local
database. Trackers execute a simple collection job, not trying to inspect trackers for DNS or if
service it working, at this phase only discovering new trackers is the goal of `harvest`
sub-command.

### Monitor

``` bash
trackers monitor --dry-run
```

Check if trackers in database are working as expected, changing the tracker status depending in
the outcome of monitoring probes. Consider [List](#list) sub-command to see the possible status.

### List

``` bash
trackers list
```

List sub-command is used to report about existing trackers. It also provide a `--etc-hosts` to
format trackers list as you would have in `/etc/hosts`, therefore you can use this to cache
trackers addresses and make sure your client only use tracker instances are only using working
instances.

Trackers may have the following status:
- **0**: service is reachable and responding;
- **1**: Can't resolve tracker's hostname;
- **2**: service does not respond;
- **3**: tracker was overwritten by `trackers overwrite`;

### Overwrite

``` bash
trackers overwrite --hostname hostname --addresses 127.0.0.1,127.0.0.2
```

Allows you to overwrite the hostname address for trackers, keep in mind a given hostname may mean
more than one tracker in the database, possibly several. On overwriting a hostname, Trackers will
inspect all services that matches the informed hostname, and probe the service using the new
hostname, only overwriting the entries that have a successful probe.

The overwritten entries will have status **3** in the database.

### Update

``` bash
trackers update --dry-run
```

Based on informed torrent-status and tracker-status, `update` sub-command will inspect running
torrents in configured torrent-client, and update it's trackers list with instances of local
database.

On [List](#list) sub-command you can see the possible options for `--status` parameter, and
regarding torrent status (`--torrent-status`), those the possible options:

- **0**: torrent is stopped;
- **1**: torrent check pending;
- **2**: torrent is checking;
- **3**: torrent download pending;
- **4**: torrent is downloading;
- **5**: torrent seed pending;
- **6**: torrent is seeding;

## Configuration

By default, Trackers expect to find configuration file at `/etc/trackers/trackers.yaml`, however, you
can change that location by using `--config` option.

In the YAML configuration file, it's expected to contain:

``` yaml
---
# transmission client settings
transmission:
  # transmission RPC URL
  url: http://127.0.0.1:9091/transmission/rpc
  # transmission RPC username
  username: user
  # transmission RPC password
  password: pass

# sqlite based persistence settings
persistence:
  # database file path
  dbPath: /var/lib/trackers/trackers.sqlite

# probe sub-command default behaviour
probe:
  # probe timeout in seconds
  timeout: 11

# nameservers used in Trackers, they must support TLS
nameservers:
  - 1.1.1.1:853 # cloudfare tls based dns server
  - 1.0.0.1:853 # cloudfare tls based dns server
```

## Workflow

A common workflow for Trackers, is the following:

``` bash
# harvest running trackers
trackers harvest && \
  # monitor functional status
  trackers monitor && \
  # filter working trackers and map them in /etc/hosts
  trackers list --etc-hosts |sort >> /etc/hosts && \
  # update running torrents with working trackers
  trackers update
```

You may also schedule sub-commands to run, so here goes a suggestion:

- `harvest`: every few minutes, for instance 10 minutes;
- `monitor`: once a day, maybe twice a day;
- `list`: and updating local `/etc/hosts`, should be execute right after you run `monitor` command;
- `update`: also should be executed right after you run `monitor` command;

## Contributing

Testing this project requires a available instance of Transmission running, with RPC enabled. To
configure tests, export the following environment variables:

- `TRANSMISSION_RPC_URL`: Transmission RPC URL;
- `TRANSMISSION_RPC_USERNAME`: Transmission RPC username;
- `TRANSMISSION_RPC_PASSWORD`: Transmission RPC password;

And then:

``` bash
make bootstrap
make test
```