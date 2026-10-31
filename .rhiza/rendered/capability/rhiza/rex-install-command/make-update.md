# make-update

capability · `01M686VEB0CY43YARC8XCA55J2` · parent [rex-install-command](../rex-install-command.md)

## acceptance

A Go test runs make -n update PREFIX=<tmp> at the repo root and asserts the go build -o <tmp>/rhiza line comes before <tmp>/rhiza rex install and no cp appears; it skips when make is not on PATH

## statement

make update builds rhiza in place into PREFIX, then runs that binary's rex install, which rewrites the Rex script and reloads Rex

## Links

### Out

- decided-by [update-is-a-make-target](../../../decision/update-is-a-make-target.md)

### In

- verifies [make-update-verified](../../../evidence/make-update-verified.md)
- verifies [make-update-qa-1](../../../evidence/make-update-qa-1.md)
- verifies [pr-133-update-in-one-command](../../../evidence/pr-133-update-in-one-command.md)

## History

- `01M686RWSEKXXBCNC7RZMT07HB` make update rebuilds rhiza in place, rewrites the Rex script and reloads Rex in one command
