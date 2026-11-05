# health-app-verified

evidence · `01M6DERCMZC7SEW48EP7W1NRYH` · actor [claude-code-b4a638ac](../actor/claude-code-b4a638ac.md)

## result

PASS

## statement

At 3fefb9a1, with RHIZA_ROOT/CHANGE/RUN unset: `go test ./cmd/rhiza -run 'TestScripts/(tui-health-app|rex-install|health|describe-golden)$' -count=1 -v` passes; tui-health-app.txtar cmps `rhiza tui app health --once` byte for byte against `rhiza health --json`, both with all ok and with a bad disk part, and finds no socket. `go test ./internal/tui/healthapp -count=1 -v` passes 6 tests on a fake runner: the view marks each part [ok]/[warn]/[bad] with the cursor's detail, r reruns health, enter on inbox opens the inbox app (keys then go to it), enter on daemon opens live, a part with no app gets a notice, and the app is registered on t. rex-install.txtar greps rhiza_app_health, tui app health and ["health"] = "t" with every other key pinned; I also ran `rhiza rex install` with a fake rex from builds at a278f03c and 3fefb9a1 and diffed the key tables: the only change is the added ["health"] = "t" line (19 keys before, 20 after). `make check` exits 0. Read the diff: tui.go hands the health app an opener for inbox and live; the autonomy part maps to an autonomy app, but the opener has none yet, so enter there shows a notice. The acceptance does not ask for that.

## Links

### Out

- verifies health-app (retired)

## History

- `01M6CRZSRFQ94MAJ8CRTEE44P2` rhiza health and a health app show the whole system's state on one screen
