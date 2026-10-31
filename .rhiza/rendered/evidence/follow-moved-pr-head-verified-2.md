# follow-moved-pr-head-verified-2

evidence · `01M68F24362DBMNEY03MQ4CC37` · actor [claude-code-508a7cb3](../actor/claude-code-508a7cb3.md)

## result

PASS

## statement

At f77d332 (attempt 2, fixing the model check's confirmed f1): the fix amends Decision follow-head-only-when-it-contains-pushed to say a head that drops the run's commits leaves ship check unsettled with a detail naming the dropped head and raises no question, which now matches internal/api/ship.go and the leaf statement; no code changed. Re-ran the acceptance: go test ./internal/api -run 'TestShipCheck(FollowsMovedHead|WaitsForPushedHead)' -count=1 passes (TestShipCheckFollowsMovedHead: dropped head stays unsettled naming 'drops the pushed 0ld0000' and abc1234; with merge-base --is-ancestor holding, --settled passes, a green check writes Head and Pushed = abc1234def, ship merge pins Head abc1234def). go test ./cmd/rhiza -run TestScripts/ship-check-follows-moved-head passes: a second clone pushes to a local bare origin, the run's ship check passes on the moved head, and a force-push dropping the run's commits leaves ship check --settled failing with 'drops the pushed'. make check: fmt-check, vet and go test ./... pass; self-check fails only on stale .rhiza/rendered actor files and index.md for this run's session actors, which the run writes and renders when it lands, not this leaf.

## Links

### Out

- verifies [follow-moved-pr-head](../capability/rhiza/orchestration/shipping/wait-green-pr/follow-moved-pr-head.md)

## History

- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing
