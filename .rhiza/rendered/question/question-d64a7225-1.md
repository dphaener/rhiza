# question-d64a7225-1

question · `01M7A1YT6CST3413RJMH65VY77`

## session

d64a7225-e14f-4ece-a633-01a1d41bff39

## statement

Change set 01M7A1P (proposed change sets no plan run covers stall): today only a plan run's gate (runRunGate) accepts, the daemon only starts plan runs from queued requests (daemonPlanRequests), and the inbox lists a proposal only when a plan run's gate proposed it (st.Gated). Which fix should the leaves build? / Change set 01M7A1P serves no intent yet. Which landed intent should it serve?

## Links

### In

- answers [uncovered-proposals-get-plan-runs](../decision/uncovered-proposals-get-plan-runs.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
