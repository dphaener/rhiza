# health-across-projects-qa-2

evidence · `01M6SYP6TTMSGEHBC92574J20K` · actor [claude-code-fea95113](../actor/claude-code-fea95113.md)

## result

PASS

## statement

QA cycle 2 at 23ada88d with a scratch build (go build -trimpath ./cmd/rhiza), RHIZA_* env cleared, read-only against the real root ~/code (nebari, rhiza, live root daemon). rhiza health at the root: machine parts once (daemon running, 2 live runs, matching daemon status at the root; usage read from the root daemon home nebari), then main/inbox/requests/roadmap/check parts tagged nebari: and rhiza:, project field set in --json. --project nebari and --project rhiza give that project only, with the daemon part still counting 2 live runs; --project nope exits 2 with not registered. The cycle-1 major (daemon part counting only the first project) is fixed. Scoped testscripts health-projects, daemon-status-projects, inbox-projects, roadmap-projects, describe-golden pass.

## Links

### Out

- verifies [health-across-projects](../capability/rhiza/orchestration/health-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
