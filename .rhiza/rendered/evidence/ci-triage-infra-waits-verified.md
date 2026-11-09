# ci-triage-infra-waits-verified

evidence · `01M6YRRZ17846ZM5DWH9PJYMXZ` · actor [claude-code-d729ffad](../actor/claude-code-d729ffad.md)

## result

PASS

## statement

Verifier at 92ba3601. Read the diff: a failing ship check (no --settled/--main) records red_head and red_checks (job ID or run ID or name, plus the attempt from the details URL) and leaves Head and Checks unset; --test-failures stays read-only; --settled --changed reports not settled while the head and check IDs/attempts match the recorded baseline, and its detail names the head, the failed checks and 'no failing test: an infra failure'; with no baseline it acts as --settled; --changed is refused without --settled. In the run engine a wait keeps the until command's last non-zero output and appends it to the timeout detail. ship-change.toml routes ci-triage on.fail to wait-ci-rerun (until 'rhiza ship check --settled --changed', 6h, max_visits 3, on.pass pr-green). RetryStep now picks the last-visited wait when two waits feed pr-green. Ran (RHIZA_* and GOTMPDIR unset, from a scratch worktree outside ~/code so the rootguard does not trip on the live daemon's usage.live.json): go test -trimpath ./internal/api/... ./internal/workflow/... gave ok (TestShipCheckFailed and TestShipCheckSettledChanged pass: baseline recorded, unchanged red head not settled and named, new job ID running not settled, then finished settled, newer attempt settled, new head settled, no baseline acts as --settled). go test -trimpath ./internal/run/... gave ok (TestWaitTimeoutDetailHoldsLastOutput: 'did not hold within 5s: <last output>'; Retry* tests pass). go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-rerun-fail|run-ship-infra-wait|run-ship-ci-fix)' gave PASS for all three. run-ship-infra-wait drives ci-triage fail to wait-ci-rerun, which passes on a rerun to pr-green (visit 2) to signoff to end.

## Links

### Out

- verifies [ci-triage-infra-waits](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-triage-infra-waits.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
