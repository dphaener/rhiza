# old-apps-deleted-verified-3

evidence · `01M6NNTCGCAEE8BWEPF5NZ7MP9` · actor [claude-code-8e1d89ce](../actor/claude-code-8e1d89ce.md)

## result

PASS

## statement

Verifier attempt 3 at e0077afc: go build ./... and go vet ./... pass; internal/tui keeps only autonomyapp, daemonapp, inboxapp, launcher, modalopen, modalshell, requestapp, runsapp (changesapp, healthapp, modelapp, readapps, roadmapapp, writeapps deleted with their tests), cmd/rhiza/tui.go has no branches for them, and the Makefile's BUDGET_PKGS drops ./internal/tui/modelapp. A scratch build run as `rhiza tui app <name> --once` exits 4 with 'unknown app' for roadmap, model, health, changes, show, live (and read, write, requests). go test ./cmd/rhiza -run TestScripts passes tui-old-apps-gone, rex-install, text-refs-display and the unchanged roadmap*, show, why*, slice*, health*, change scripts (the diff touches none of the wrapped commands' scripts). go test ./internal/tui/... and ./internal/quiettest pass; ./internal/api passes with the timing budgets skipped (TestReadBudget missed its what-if medians under CPU contention from parallel runs; make test runs it alone under lock).

## Links

### Out

- verifies [old-apps-deleted](../capability/rhiza/old-apps-deleted.md)

## History

- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
