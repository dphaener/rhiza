# update-target-documented-verified

evidence · `01M69924G7JSM94BDG0ZD0MX07` · actor [claude-code-d4a39548](../actor/claude-code-d4a39548.md)

## result

PASS

## statement

At bf71b8b: make help prints 'update  Rebuild into $(PREFIX), then rewrite the Rex script and reload Rex'. TestMakeUpdateDryRun (cmd/rhiza/makeupdate_test.go), the make-update test, now also runs make help and asserts a line whose first field is update (fixes f2). TestClaudeMDCommandsNameMakeUpdate reads CLAUDE.md, extracts the fenced block under '## Commands' and requires a line starting with 'make update' (fixes f1); CLAUDE.md has 'make update     # make install, then rewrite the Rex script and reload Rex' there. go test ./cmd/rhiza -run 'TestMakeUpdateDryRun|TestClaudeMDCommandsNameMakeUpdate' -count=1 -v: both PASS. Mutation in a scratch git-archive copy (moved make update out of the Commands block into trailing prose): TestClaudeMDCommandsNameMakeUpdate FAILs. go vet and go test ./... pass (modelapp timing test flaked once at 62ms under load, passed on rerun). make check's self-check fails only on stale .rhiza/rendered actor/index files from this run's in-flight actors, not this leaf's code.

## Links

### Out

- verifies [update-target-documented](../capability/rhiza/rex-install-command/update-target-documented.md)

## History

- `01M686RWSEKXXBCNC7RZMT07HB` make update rebuilds rhiza in place, rewrites the Rex script and reloads Rex in one command
