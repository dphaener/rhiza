# daemon-status-across-projects-qa-1

evidence · `01M6STV5BXZ65VV4WRBPE9TFFG` · actor [claude-code-8ba077ae](../actor/claude-code-8ba077ae.md)

## result

PASS

## statement

QA at 538fc842 with a scratch build (go build -trimpath ./cmd/rhiza), read-only against the machine's real root ~/code (projects nebari and rhiza, live root daemon), RHIZA_* env cleared. rhiza daemon status --json at the root: the root daemon once (running true, top-level runs and skipped empty), a projects part for nebari (0 runs, its last_started) and for rhiza (the 2 live runs of this milestone, 13 skipped change sets with reasons such as usage reserve and umbrella), skipped_projects empty (both initialised). Human text shows Project nebari:/Project rhiza: sections. --project nebari gives only nebari's part (no projects entries); run inside ~/code/rhiza it gives rhiza's 2 runs and 13 skips at the top level; --all from inside rhiza covers both; --project nope exits 2 with not registered. Nit seen: at the root the top-level human line reads Live runs: none above the per-project sections that list 2 live runs.

## Links

### Out

- verifies [daemon-status-across-projects](../capability/rhiza/orchestration/daemon-status-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
