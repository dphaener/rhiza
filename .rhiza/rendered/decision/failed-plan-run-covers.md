# failed-plan-run-covers

decision · `01M7A26BZ2KESAHX4212BR36EK`

## reasoning

A plan-review run that fails before its gate (milestone-check, a stuck review) journals no gated entry; under the arbiter's rule alone the daemon would start a new run on every poll. The failed run is already an inbox item with rhiza run resume --retry next, so the proposal is not silent

## statement

A run on the change set of any workflow, plan-review included, that ended failed and was not abandoned also covers it: the daemon does not start another, and the inbox lists the failure with its retry instead of an uncovered proposal; an abandoned run covers nothing

## wording

(derived by the planner from the arbiter's Decision uncovered-proposals-get-plan-runs, which names live runs and gated finished runs)

## Links

### Out

- derived-from [uncovered-proposal-plan](../source/uncovered-proposal-plan.md)

### In

- decided-by [daemon-plans-uncovered-proposals](../capability/daemon-plans-uncovered-proposals.md)
- decided-by [inbox-lists-uncovered-proposals](../capability/inbox-lists-uncovered-proposals.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
