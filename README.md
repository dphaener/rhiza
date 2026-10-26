# Rhiza

A living, connected model of a product (what it is, what it should become, and why), built for working with coding agents.

Rhiza was called Genesis until its rename in milestone 2. Its logs keep the names they were written with.

Early days: nothing to install yet.

## Development

`make check` runs everything CI runs. `make install` builds the `rhiza` binary into `~/bin` (set `PREFIX` to change it). The Go module is `github.com/dphaener/rhiza`. Rhiza keeps its own log in `.rhiza/`, and the environment variables it reads start with `RHIZA_`.
