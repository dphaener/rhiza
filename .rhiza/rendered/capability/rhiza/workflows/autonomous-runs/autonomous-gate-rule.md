# autonomous-gate-rule

capability · `01M5WRHE82S6TFSBBQ3VN42HGQ` · parent [autonomous-runs](../autonomous-runs.md)

## acceptance

A testscript proposes a change set and runs rhiza accept --rule autonomous three ways: with no RHIZA_RUN, and with RHIZA_RUN naming a run that is not autonomous, it exits 2 with nothing written; with RHIZA_RUN naming an unfinished autonomous run it accepts, and change show --json gives rule autonomous and a reason naming the run.

## statement

An owner's gate skipped in autonomous mode is recorded with rule autonomous: rhiza accept --rule autonomous is allowed only when RHIZA_RUN names an unfinished autonomous run, and with no --reason the entry's reason names that run. Anywhere else the rule is refused, so it always means an autonomous run made the move.

## Links

### Out

- depends-on [run-start-autonomous](run-start-autonomous.md)
- decided-by [autonomous-mode](../../../../decision/autonomous-mode.md)
- decided-by [two-owner-gates](../../../../decision/two-owner-gates.md)
- decided-by [autonomous-accepts-journalled](../../../../decision/autonomous-accepts-journalled.md)
- decided-by [reject-guarded-like-accept](../../../../decision/reject-guarded-like-accept.md)
- decided-by [autonomous-rule-any-change-set](../../../../decision/autonomous-rule-any-change-set.md)
- decided-by [journal-failure-reported](../../../../decision/journal-failure-reported.md)

### In

- depends-on [autonomous-review](autonomous-review.md)
- verifies [autonomous-gate-rule-verified](../../../../evidence/autonomous-gate-rule-verified.md)
- verifies [pr-39-autonomous-mode-and-the](../../../../evidence/pr-39-autonomous-mode-and-the.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
