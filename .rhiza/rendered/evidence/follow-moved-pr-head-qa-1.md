# follow-moved-pr-head-qa-1

evidence · `01M68E9AB9S8WZ0FK45HCWTFZN` · actor [claude-code-8d3acb4b](../actor/claude-code-8d3acb4b.md)

## result

PASS

## statement

At 14d9f17, as QA with a scratch build (go build -o $TMPDIR) in a scratch repo with a local bare origin and a fake gh that reports origin's branch head: after the run's own ship rebase pushed eb82547 (recorded as pushed), a second clone pushed 0b34323 on top; ship check --settled exited 0 at once, a green ship check wrote head=pushed=0b34323, and ship merge --json returned head 0b34323. A force-pushed rewrite 592a57a that dropped those commits gave 'not settled: the pull request's head 592a57a drops the pushed 0b34323', and --settled exited 1. go test ./internal/api -run TestShip and the ship-check-follows-moved-head testscript pass.

## Links

### Out

- verifies [follow-moved-pr-head](../capability/rhiza/orchestration/shipping/wait-green-pr/follow-moved-pr-head.md)

## History

- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing
