# update-target-documented

capability · `01M686VFHDHYAJP8B10B23TMSC` · parent [rex-install-command](../rex-install-command.md)

## acceptance

The same test asserts make help's output names update, and a test reads AGENTS.md and finds make update in its commands

## statement

make help and AGENTS.md's command list name make update

## Links

### Out

- decided-by [update-is-a-make-target](../../../decision/update-is-a-make-target.md)

### In

- verifies [update-target-documented-verified](../../../evidence/update-target-documented-verified.md)
- verifies [update-target-documented-qa-1](../../../evidence/update-target-documented-qa-1.md)
- verifies [pr-133-update-in-one-command](../../../evidence/pr-133-update-in-one-command.md)
- verifies [update-target-documented-agents-md-qa-1](../../../evidence/update-target-documented-agents-md-qa-1.md)
- verifies [pr-494-agents-md-entry-point](../../../evidence/pr-494-agents-md-entry-point.md)

## History

- `01M686RWSEKXXBCNC7RZMT07HB` make update rebuilds rhiza in place, rewrites the Rex script and reloads Rex in one command
- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
