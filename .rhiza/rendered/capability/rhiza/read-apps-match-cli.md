# read-apps-match-cli

capability · `01M62JP9RTA2QHK5B6A53852M8` · parent [rhiza](../rhiza.md)

## acceptance

For each of the four commands, a testscript on a fixture model checks the app's result equals the command's --json output, and the test starts no daemon and binds no socket

## statement

ready, show, why and slice run as apps that call the internal/api registry in-process and show the same result as the matching rhiza command --json

## Links

### Out

- serves [tui-read-apps](../../intent/rhiza-tui/tui-read-apps.md)
- decided-by [tui-many-mini-apps](../../decision/tui-many-mini-apps.md)

### In

- verifies [read-apps-match-cli-verified](../../evidence/read-apps-match-cli-verified.md)
- verifies [read-apps-match-cli-qa-1](../../evidence/read-apps-match-cli-qa-1.md)
- verifies [read-apps-match-cli-verified-2](../../evidence/read-apps-match-cli-verified-2.md)
- verifies [read-apps-match-cli-qa-2](../../evidence/read-apps-match-cli-qa-2.md)
- verifies [read-apps-match-cli-qa-3](../../evidence/read-apps-match-cli-qa-3.md)
- verifies [pr-65-tui-shell-and-read-apps](../../evidence/pr-65-tui-shell-and-read-apps.md)
- verifies [pr-122-remove-the-all-apps-picker](../../evidence/pr-122-remove-the-all-apps-picker.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
- `01M6489HCR688XKH99YV9E2VGY` Remove the all-apps picker, as decision rex-picker-removed says: rhiza tui launch and the rhiza_launch action written by rhiza rex install go, and the capabilities that describe the picker are retired or revised, so each app opens only from its own key in the rhiza key mode or its own action in Rex's command palette.
