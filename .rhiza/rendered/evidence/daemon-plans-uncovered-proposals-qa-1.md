# daemon-plans-uncovered-proposals-qa-1

evidence · `01M7ANME9MD0K3WFXASPB18WKT` · actor [claude-code-8f5699a0](../actor/claude-code-8f5699a0.md)

## result

PASS

## statement

QA with a fresh build of the branch (go build -trimpath to a scratch path), in a scratch root with Rex faked and RHIZA_TEST_DETACH_LOG: rhiza daemon --once --permission-mode manual --max-runs 4 started 'run start --foreground <cs> --workflow plan-review --autonomous --permission-mode manual' for the two proposals with a leaf and no run and for one whose only run was abandoned; started nothing for a proposal with no leaves, one a finished plan-review run gated, one whose plan-review run failed, or one a plan-milestone run on another change set gated. With a live plan-review run on a proposal it was not restarted; --permission-mode auto passed through; an accepted change set went to the ship pick; with one slot used and --max-runs 2 the plan-review start took the free slot and the ship pick waited; with --max-runs 1 and the slot used nothing started; nothing started under rhiza.daemon.paused true or with low disk. A pending start is held only within one daemon process (separate --once processes restart it), as for ship starts. Scoped tests: go test -trimpath ./internal/api -run 'Proposals|Inbox' and TestScripts/daemon-plans-proposals pass.

## Links

### Out

- verifies [daemon-plans-uncovered-proposals](../capability/daemon-plans-uncovered-proposals.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
