# follow-moved-pr-head-verified

evidence · `01M68CT6EWFJVJXNN8RWFEFG0W` · actor [claude-code-f4ee0ef5](../actor/claude-code-f4ee0ef5.md)

## result

PASS

## statement

At 1324262, as verifier: go test ./internal/api -run TestShipCheck -count=1 passes, including the new TestShipCheckFollowsMovedHead beside TestShipCheckWaitsForPushedHead. It uses a fake git and gh. With no merge-base --is-ancestor answer, --settled fails with 'drops the pushed 0ld0000' and names abc1234. Once the ancestor check holds, --settled passes without writing, and a green check records Head=Pushed=abc1234def. ship.merge then returns Head abc1234def. The older-head case still reports 'not the pushed ... yet'. The testscript ship-check-follows-moved-head.txtar passes: a second clone pushes to a local bare origin and ship check passes, then a force-pushed rewrite fails --settled with 'drops the pushed'. Its merge step does not itself compare the merged head with the moved head, so I checked that by hand. With a scratch build in a scratch repo, the same fake gh and a second-clone push on top of the landed head: ship check --settled (the wait-pr until command) exits 0 at once. The green check writes head=pushed=the moved sha. ship merge --json pins head to the moved sha. A force-pushed rewrite gives 'not settled: the pull request's head b46ecbd drops the pushed af90358', exit 1. make check passes. Reading the diff: ship.go fetches the branch and adopts v.Head as Pushed only when it contains the pushed head, and keeps the 'not yet' wait when the PR head is behind.

## Links

### Out

- verifies [follow-moved-pr-head](../capability/rhiza/orchestration/shipping/wait-green-pr/follow-moved-pr-head.md)

## History

- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing
