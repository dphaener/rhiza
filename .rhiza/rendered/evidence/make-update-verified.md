# make-update-verified

evidence · `01M698787RT2KNQF2W45M6BS7N` · actor [claude-code-9f0fa501](../actor/claude-code-9f0fa501.md)

## result

PASS

## statement

At b297efa: make -n update PREFIX=<mktemp dir> printed mkdir -p, then go build ... -o <tmp>/rhiza ./cmd/rhiza, then <tmp>/rhiza rex install, with no cp (the update target depends on install, which builds in place). go test ./cmd/rhiza -run TestMakeUpdateDryRun -count=1 -v PASSED. The test asserts the go build -o <tmp>/rhiza line comes before <tmp>/rhiza rex install and that no cp field appears. A compiled test binary run with PATH=/nonexistent SKIPPED with 'make is not on PATH'. rhiza rex install --help confirms it rewrites the Rex config and reloads Rex. fmt-check, vet and go test ./... all pass. make check fails only in self-check, on stale .rhiza/rendered actor and index files from this run's own bookkeeping, not on the leaf's diff.

## Links

### Out

- verifies [make-update](../capability/rhiza/rex-install-command/make-update.md)

## History

- `01M686RWSEKXXBCNC7RZMT07HB` make update rebuilds rhiza in place, rewrites the Rex script and reloads Rex in one command
