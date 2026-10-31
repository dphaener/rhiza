# tui-launch-usage-qa-1

evidence · `01M686VFKGRRFNGMMK40G6PWR4` · actor [claude-code-17358f03](../actor/claude-code-17358f03.md)

## result

PASS

## statement

At dccf0f0 built go build -o $TMPDIR/rhiza-qa ./cmd/rhiza and ran it by hand: 'rhiza tui launch' printed 'usage: rhiza tui app <name> [args] [--once]' on stderr and exited 2; bare 'rhiza tui' did the same; nothing opened and stdout was empty. rhiza --help lists only 'tui app', no launch.

## Links

### Out

- verifies [tui-launch-usage](../capability/rhiza/rex-app-actions/tui-launch-usage.md)

## History

- `01M6489HCR688XKH99YV9E2VGY` Remove the all-apps picker, as decision rex-picker-removed says: rhiza tui launch and the rhiza_launch action written by rhiza rex install go, and the capabilities that describe the picker are retired or revised, so each app opens only from its own key in the rhiza key mode or its own action in Rex's command palette.
