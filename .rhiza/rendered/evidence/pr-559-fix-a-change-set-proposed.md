# pr-559-fix-a-change-set-proposed

evidence · `01M7B0NNBX8TXY00W1ANV5HVTV` · actor [dphaener](../actor/dphaener.md)

## head

651f9a96c0e79f7d1af2290fb74d954e1e8c1e05

## pr

https://github.com/dphaener/rhiza/pull/559

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #559 (https://github.com/dphaener/rhiza/pull/559) passed CI at 651f9a9

## status

verified

## Links

### Out

- verifies [plan-review-workflow](../capability/plan-review-workflow.md)
- verifies [daemon-plans-uncovered-proposals](../capability/daemon-plans-uncovered-proposals.md)
- verifies [inbox-lists-uncovered-proposals](../capability/inbox-lists-uncovered-proposals.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
