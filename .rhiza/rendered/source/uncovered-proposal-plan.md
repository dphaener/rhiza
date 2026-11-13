# uncovered-proposal-plan

source · `01M7A255MHCF19H738M4774SJF`

## quotes

```json
[]
```

## raw_sha256

ff1d235d6949ef452df06bed2f558eea01ab7913326d94ce73f312ae799d65af

## session

d64a7225-e14f-4ece-a633-01a1d41bff39

## summary

Autonomous plan run of bug change set 01M7A1P: change sets proposed outside a plan run (by a PR) are never accepted, since only a plan run's gate (runRunGate) accepts, the daemon starts plan runs only from queued requests (daemonPlanRequests), and the inbox lists a proposal only when a plan run's gate named it. No owner was present. The planner asked two questions (the fix's shape, which intent it serves); the arbiter answered with Decision uncovered-proposals-get-plan-runs: the daemon starts an autonomous plan run for every proposed change set with leaves that no plan run covers (no live run on it, no finished plan run whose gated entry names it), within maxRuns and the holds and with the daemon's permission mode; the run reviews the existing leaves and reaches its gate, which stays the only place that accepts under rhiza.autonomy.accept; the inbox lists such a proposal as waiting for acceptance while no plan run covers it or when the gate left it for the owner; rhiza ready keeps 'not accepted: proposed'; the change set serves m4-orchestrated-by-rhiza. The planner derived two rulings from that: the run is a builtin plan-review workflow that includes plan and skips its interview, so the leaves are reviewed, not re-planned; and a plan run on the change set that ended failed (not abandoned) also covers it, so a failing run is not restarted every poll. Nothing is left open.

## title

Plan of the uncovered-proposal fix, autonomous, with the arbiter

## Links

### In

- derived-from [plan-review-skips-interview](../decision/plan-review-skips-interview.md)
- derived-from [failed-plan-run-covers](../decision/failed-plan-run-covers.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
