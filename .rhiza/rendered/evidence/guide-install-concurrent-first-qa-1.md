# guide-install-concurrent-first-qa-1

evidence · `01M64TX6GRMZKBXBT68BQR87VY` · actor [claude-code-0678fac1](../actor/claude-code-0678fac1.md)

## result

PASS

## statement

At a6b1a5d, built rhiza with go build to a scratch path. Real processes: 200 rounds of 8 concurrent `rhiza guide install --agent claude` into a fresh HOME each: 0/1600 failures. With the binary: removing plugin.json from an installed guide and reinstalling writes it back (1 written, 6 unchanged, exit 0); no plugin.json plus a SKILL.md holding other bytes gives the 'exists and was not written by rhiza guide install' conflict (exit 3) with next --force, and --force installs; a foreign plugin.json is still a conflict; an empty folder installs all 7 files. go test ./internal/api -run 'TestGuideInstallConcurrent$' -count=300 -race -cpu 1,2,4,8 passes at HEAD (ok, 30.5s) and fails 4 times at the merge-base with the foreign-folder error, so the fix removes the flake at its cause. internal/api/guide_install_test.go is unchanged from the merge-base. make check exits 0 with the run's RHIZA_* env unset.

## Links

### Out

- verifies [guide-install-concurrent-first](../capability/rhiza/guide/guide-install-concurrent-first.md)

## History

- `01M5ZNVHN2YNK6DDVYY8GK5CEQ` TestGuideInstallConcurrent passes reliably: the flake that fails it on main too (seen once in make check during the build workflow's review) is found and fixed at its cause, without skipping or loosening the test.
