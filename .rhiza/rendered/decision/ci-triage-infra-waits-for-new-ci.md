# ci-triage-infra-waits-for-new-ci

decision · `01M6YGVRV397Z6Q7B41Z7PXKN1`

## reasoning

The change set's intent says an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. In internal/workflow/builtin/ship-change.toml, ci-triage has only on.pass. The question notes that pointing on.fail straight at wait-pr would loop: `rhiza ship check --settled` passes at once on a settled red head, so the run would burn max_visits in seconds and end anyway. A wait that passes only when there is new CI (a rerun by the owner or by GitHub, or a new push) fixes that, and it covers exactly the PR #393 case: billing fixed, one rerun, green. It needs no human at the terminal. The timeout and max_visits still bound a run whose CI never recovers, so it fails there with a clear reason instead of hanging. Workflow edges are static, so the wait cannot return to the check that routed into ci-triage (pr-green or rebased-green). Routing to pr-green is the simplest correct target and matches the graph's existing rule that any re-entry after a red head goes back through signoff (ci-fix -> render-pr -> wait-pr -> pr-green -> signoff). The retry decision 01M6YGPS3 (resume --retry reopens at wait-pr or wait-rebased from the journal) is separate and still holds; with this on.fail in place, fewer runs reach a failed ci-triage at all.

## statement

In ship-change, ci-triage gets on.fail = a new wait step (wait-ci-rerun) instead of ending the run. Its until is `rhiza ship check --settled` with a new flag (e.g. --changed) that reports settled only once CI has moved past the state ci-triage judged: a new head SHA, or a check run/attempt newer than the one in the head's latest ship-check record, and that new CI has finished. So an infra-only failure never loops at once on the already-settled red head. The wait has a long timeout (6h, longer than the 2h CI waits, since an infra cause like billing needs a person) and max_visits = 3; its timeout outcome ends the run with a line naming the infra failure. Its on.pass goes to pr-green, which re-checks the head: green continues to signoff, red with failing tests goes to ci-triage/ci-fix, red without them waits again. A post-rebase run that comes back this way passes signoff again, the same as a ci-fix commit already sends it back through render-pr, open-pr and wait-pr.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-bc50c2e7-2](../question/question-bc50c2e7-2.md)
- derived-from [resume-retry-ci-triage-plan](../source/resume-retry-ci-triage-plan.md)

### In

- decided-by [ci-triage-infra-waits](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-triage-infra-waits.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
