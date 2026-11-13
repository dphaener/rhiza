# uncovered-proposals-get-plan-runs

decision · `01M7A22M37FB7599ZNJ6J69TWA`

## reasoning

The change set's intent names both remedies ('the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals') and its title has two defects: the proposal is never accepted, and nothing says so. Only the daemon half removes the hand step (the orchestrator ran rhiza run start --workflow plan for nine change sets), and only the inbox half covers accept = owner, where an auto-started run cannot accept. Starting a plan run rather than accepting directly in the daemon keeps acceptance in one place (runRunGate, internal/api/run_review.go) and gives proposals written outside a run the same plan review that run-made proposals get. The inbox today lists only st.Gated proposals (internal/api/inbox.go:344), so widening it to uncovered proposals is a small change in the same code. For the intent: @m4-orchestrated-by-rhiza says change sets ship with the daemon orchestrating and no orchestrator session, which this stall broke. The closest earlier fix, 01M797TR (hand-queued plan runs get the daemon's permission mode), serves the same intent. @m4-runs-hold-up is about runs that finish without repair, and here no run had started.

## statement

The leaves build both halves of the fix. (1) The daemon starts an autonomous plan run for every proposed change set that already has leaves and that no plan run covers (no live run on it, and no finished plan run whose gated entry names it). It counts against maxRuns and the load holds like any other start, and it takes the daemon's permission mode. The plan run reviews the existing leaves and reaches its gate, so runRunGate stays the only place that accepts, and rhiza.autonomy.accept decides there as it does today. (2) The inbox lists such a proposal as waiting for acceptance whenever no plan run covers it yet or the gate left it for the owner, so a proposed change set is never silent. rhiza ready keeps its 'not accepted: proposed' reason. Change set 01M7A1P serves the landed intent @m4-orchestrated-by-rhiza.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-d64a7225-1](../question/question-d64a7225-1.md)

### In

- decided-by [plan-review-workflow](../capability/plan-review-workflow.md)
- decided-by [daemon-plans-uncovered-proposals](../capability/daemon-plans-uncovered-proposals.md)
- decided-by [inbox-lists-uncovered-proposals](../capability/inbox-lists-uncovered-proposals.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
