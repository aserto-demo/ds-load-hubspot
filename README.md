# ds-load-hubspot
A sample `ds-load`` plugin for getting contact and company data out of hubspot

## Arguments

```
Usage: ds-load-hubspot <command>

Hubspot data loader

Commands:
  version              version information
  fetch                fetch hubspot data
  transform            transform hubspot data
  export-transform     export default transform template
  exec                 fetch and transform hubspot data
  get-refresh-token    obtain a refresh token from hubspot

Flags:
  -h, --help                  Show context-sensitive help.
  -c, --config=CONFIG-FLAG    Configuration file path
  -v, --verbosity=INT         Use to increase output verbosity.
```

The default command for ds-load-hubspot is `exec`, which will run both `fetch` and `transform`.

To view command-specific help: `ds-load-hubspot fetch --help`.

### Env variables
Parameters can also be passed by environment, as seen in the help message of each command, but the ones from config files and command line take precedence.

## Config files

config files are in yaml format:
```yaml
---
ds-load-arg: value
ds-load-another-arg: value

hubspot:
  arg: value
```

When passing custom config files to both the cli and the plugin, use `ds-load -c <config-path> <plugin-name> <command>` or `ds-load-hubspot -c <config-path> <command>`.

### CLI config
default location: `~/.config/ds-load/cfg/config.yaml` can be overridden using `-c/--config`

#### ds-load config example with hubspot plugin
```yaml
---
host: directory.eng.aserto.com:8443
api-key: secretapikey
tenant-id: your-tenant-id
hubspot:
  private-access-token: your-access-token
  contacts: false
  companies: true
```

### Plugin config
default location: `~/.config/ds-load/cfg/hubspot.yaml` can be overridden using `-c/--config`

#### example for hubspot
```yaml
---
hubspot:
  private-access-token: your-access-token
  contacts: false
  companies: true
```

## Transform
The data received from the fetcher is being transformed using a transformation template, which is written as a go template and it outputs objects and relations.

The default transformation template can be exported using `ds-load-hubspot export-transform`.

A custom transformation file can be provided when running the plugin in `exec` or `transform` mode via the `--template-file` parameter.

## Logs

Logs are printed to `stdout`. You can increase detail using the verbosity flag (e.g. `-vvv`).

## Usage examples

### Import from hubspot into the directory
```
ds-load --host=<directory-host> --api-key=<directory-api-key> --tenant-id=<tenant-id> hubspot --private-access-token=<PAT> --contacts --companies
```

### Import data with custom transformation file
```
ds-load --host=<directory-host> --api-key=<directory-api-key> --tenant-id=<tenant-id> hubspot --private-access-token=<PAT> --contacts --companies --template-file=<template-path>
```

### View contact data from hubspot
```
ds-load-hubspot fetch --private-access-token=<PAT> --contacts 
```

### Transform data from a previously saved hubspot fetch
```
ds-load-hubspot fetch --private-access-token=<PAT> --contacts > hubspot.json
cat hubspot.json | ds-load-hubspot transform
```

### Transform and import data from a previously saved hubspot fetch
```
ds-load hubspot fetch --private-access-token=<PAT> --contacts --companies > hubspot.json

cat hubspot.json | ds-load --host=<directory-host> --api-key=<directory-api-key> --tenant-id=<tenant-id> hubspot transform
```

### Pipe data from fetch to transform
```
ds-load-hubspot fetch --private-access-token=<PAT> --contacts | ds-load-hubspot transform
```

### Use config file to import data from hubspot into the directory

config.yaml
```yaml
---
host: "directory.eng.aserto.com:8443"
api-key: "secretapikey"
tenant-id: "your-tenant-id"
hubspot:
  private-access-token: "your-PAT"
  contacts: true
```

```
ds-load -c ./config.yaml hubspot
```
