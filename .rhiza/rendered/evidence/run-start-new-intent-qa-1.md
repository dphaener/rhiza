# run-start-new-intent-qa-1

evidence · `01M6A8S0C4WXK59R9CYZ34QENB` · actor [claude-code-b0fb6bcc](../actor/claude-code-b0fb6bcc.md)

## result

PASS

## statement

QA at cafec91c with a scratch build of ./cmd/rhiza in a scratch git repo (real git, fake rex and claude on PATH, RHIZA_* run env unset). run start --workflow plan --new X <cs> exits 2 'takes no change set'; --workflow build --new X exits 2 'workflow "build" is not plan'; --milestone without --new exits 2; none opened a change set or a run. With RHIZA_TEST_NOW/SEED fixed, --dry-run --json predicted change 01M6A6FYW86T3BKSK9V4MCGPCP and wrote nothing; the real run (--permission-mode manual) opened that same ID as a draft intention 'Plan reports', made it current, created worktree repo.plan-reports on cs/plan-reports, started the agent with 'Plan 01M6A6FYW86T3BKSK9V4MCGPCP', and the run journal's change field named it. --new ... --milestone opened a draft record. With RHIZA_SESSION set the start is refused (exit 2) and no draft is left behind.

## Links

### Out

- verifies [run-start-new-intent](../capability/rhiza/workflows/workflow-runs/run-start/run-start-new-intent.md)

## History

- `01M686RVAD1EHGERBBDJN3S45G` The roadmap app starts a plan run for a new milestone or a new change set from a typed one-line intent
