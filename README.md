# Tcp Server written in Go

This is a relatively simple Tcp server but can be extended by writing plugins

## Build

### Dependencies

- [golang](https://go.dev)
- GNU Make

### Build server

    make build

### Build plugins

Build all plugins from `plugin_*` dirs and and save it in the `plugins` directory

    make build-plugins
