# guide-install-concurrent-first-verified

evidence · `01M64T9NP0WAHW57YB0W8V4K4E` · actor [claude-code-cca964e4](../actor/claude-code-cca964e4.md)

## result

PASS

## statement

At 8a3af6f (detached checkout): go test ./internal/api -run 'TestGuideInstallConcurrent$' -count=300 -race -cpu 1,2,4,8 passed (ok, 12.3s, uncached). TestGuideInstallConcurrent is unchanged: the leaf's two commits (d892a75, 8a3af6f) touch only internal/api/guide_install.go and cmd/rhiza/testdata/script/guide-install.txtar. make check passed (exit 0: fmt-check, vet, go test ./..., self-check) with the run's RHIZA_* env unset; with them set, self-check reports only stale renders of this run's own session actors in the shared store, not of the commit. Diff read: pluginForeign now compares the other plugin files' bytes to want, so a racing install's partial write is not foreign, and the guide-install script's --project conflict case now uses a SKILL.md with different bytes, fixing attempt 1's finding f1.

## Links

### Out

- verifies [guide-install-concurrent-first](../capability/rhiza/guide/guide-install-concurrent-first.md)

## History

- `01M5ZNVHN2YNK6DDVYY8GK5CEQ` TestGuideInstallConcurrent passes reliably: the flake that fails it on main too (seen once in make check during the build workflow's review) is found and fixed at its cause, without skipping or loosening the test.
