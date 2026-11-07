# roadmap-across-projects-qa-1

evidence · `01M6STV5FD2A8MKAT1K5E6DT4F` · actor [claude-code-8ba077ae](../actor/claude-code-8ba077ae.md)

## result

PASS

## statement

QA at 538fc842 with a scratch build. Scratch root of two hand-built projects (alpha with milestone ma and an accepted change set serving it, beta with a backlog change set): rhiza roadmap at the root prints ## alpha: ma and ## beta: Backlog, and --json tags every group and change set with its project; from inside beta only beta's backlog untagged, --all from inside beta both; --project alpha only alpha's milestone untagged; --project nope exits 2 not registered. Real root ~/code: only rhiza's groups (milestone-29, Backlog, areas) tagged rhiza, as nebari has no in-flight change sets (main's rhiza roadmap in nebari also shows Nothing planned yet).

## Links

### Out

- verifies [roadmap-across-projects](../capability/rhiza/orchestration/roadmap-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
