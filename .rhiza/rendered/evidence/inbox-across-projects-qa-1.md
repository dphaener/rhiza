# inbox-across-projects-qa-1

evidence · `01M6STV5E9K5NK7AFNXDXXWMBN` · actor [claude-code-8ba077ae](../actor/claude-code-8ba077ae.md)

## result

PASS

## statement

QA at 538fc842 with a scratch build. Real root ~/code, read-only: rhiza inbox at the root lists rhiza's failed run, nebari's failed run (each with project: line) and the machine-wide usage-reserve item once, untagged; inside nebari only nebari's items; --project rhiza from inside nebari only rhiza's; --all with --project exits 2. Scratch root of two hand-built projects alpha and beta (hermetic HOME, fake rex on PATH, the testscript's journal fixtures): inbox at the root tags the alpha question and beta sign-off; inbox open of alpha's question from the root focused block:q1 (rex.log shows block inspect and focus); from inside beta it is no inbox item, with --all it focuses; inbox approve --foreground --dry-run of beta's sign-off from the root reached beta's run SIGNOFFRUN01 (its fixture workflow has no steps); approve of the question is refused naming inbox open.

## Links

### Out

- verifies [inbox-across-projects](../capability/rhiza/orchestration/inbox-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
