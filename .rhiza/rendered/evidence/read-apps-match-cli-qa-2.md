# read-apps-match-cli-qa-2

evidence · `01M630C4JW2MAQA002D39ZJW4Q` · actor [claude-code-f3ca46cd](../actor/claude-code-f3ca46cd.md)

## result

PASS

## statement

QA at c7a97c9 (cycle 2) with a built rhiza on this repo's real model (RHIZA_ROOT/RUN/CHANGE/FOR unset). cmp shows `rhiza <op> --json` and `rhiza tui app <op> --once` byte-identical for roadmap (19626 bytes), ready (14007), and, with --change 01M62JMSW4FH8RAWSQXHC7ZPND, slice (28179), show tui-launch-key (1332) and why quick-launch-layer (37320); both exit 0. The cycle-1 gap (tui app dropped --change, exit 2) is fixed. show on an unknown handle exits 4 in both. In real Rex the roadmap and ready apps open as floating modals with the j/k scroll footer.

## Links

### Out

- verifies read-apps-match-cli (retired)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
