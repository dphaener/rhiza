# inbox-lists-uncovered-proposals-qa-1

evidence · `01M7ANMEK32WF805VJCE83D7AB` · actor [claude-code-8f5699a0](../actor/claude-code-8f5699a0.md)

## result

PASS

## statement

QA with the scratch build: rhiza inbox --json listed each proposal with a leaf and no covering run as one accept item, detail 'waiting for acceptance: no plan run covers it', next 'rhiza accept --rule owner --change <id>', under rhiza.autonomy.accept auto; a proposal gated by a plan-review or plan-milestone run kept its single 'waiting for acceptance' item and was not listed twice; one whose plan-review run failed showed only its failed item; after rhiza accept (the item's next) and after a live plan-review run journaled, the item was gone; rhiza ready still says 'not accepted: proposed'. Rough edge: plain-text inbox and inbox open print 'run ,' for these run-less items (minor finding). TestScripts/inbox-uncovered-proposal and inbox-owner-run-accept pass.

## Links

### Out

- verifies [inbox-lists-uncovered-proposals](../capability/inbox-lists-uncovered-proposals.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
