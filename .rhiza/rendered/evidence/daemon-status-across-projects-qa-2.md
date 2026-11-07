# daemon-status-across-projects-qa-2

evidence · `01M6SYP70RHZJW2D685DZTZDS7` · actor [claude-code-fea95113](../actor/claude-code-fea95113.md)

## result

PASS

## statement

QA cycle 2 at 23ada88d with a scratch build (go build -trimpath ./cmd/rhiza), RHIZA_* env cleared. On the real root ~/code, rhiza daemon status --json: root daemon once (running), projects parts nebari (0 runs, last_started) and rhiza (2 live runs of this milestone, 12 skips with reasons: umbrella, usage reserve), skipped_projects empty. Text shows Project nebari:/Project rhiza: sections with live runs, last started/resumed and Skipped lines. daemon-status-projects testscript (now asserting each project live run and skip under its name) passes.

## Links

### Out

- verifies [daemon-status-across-projects](../capability/rhiza/orchestration/daemon-status-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
