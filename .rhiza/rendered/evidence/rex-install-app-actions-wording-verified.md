# rex-install-app-actions-wording-verified

evidence · `01M6864BHEPJS6V3GV0DWXFB6P` · actor [claude-code-97a3a8dd](../actor/claude-code-97a3a8dd.md)

## result

PASS

## statement

At 36f3811: go test ./cmd/rhiza -run 'TestScripts/(rex-install|describe-golden)$' -count=1 passed; rex-install.txtar asserts 'Installed Rhiza.s app actions and key mode' and 'Removed Rhiza.s app actions and key mode from Rex' with ! stdout 'quick-launch' after both install and uninstall. Re-running describe-golden with -update left the tree clean (golden current), and grep finds 0 'quick-launch' in describe-golden.txtar. The diff changes the rex.install/rex.uninstall summaries, request doc comments and Human output in internal/api/rex_install.go; no quick-launch text remains in Go or txtar files.

## Links

### Out

- verifies [rex-install-app-actions-wording](../capability/rhiza/rex-install-command/rex-install-app-actions-wording.md)

## History

- `01M6489HCR688XKH99YV9E2VGY` Remove the all-apps picker, as decision rex-picker-removed says: rhiza tui launch and the rhiza_launch action written by rhiza rex install go, and the capabilities that describe the picker are retired or revised, so each app opens only from its own key in the rhiza key mode or its own action in Rex's command palette.
