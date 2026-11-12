# ship-gate-fails-own-model-conflict-verified

evidence · `01M77S6KBBFS8AZGK08KZH5EDK` · actor [claude-code-f69cd2cd](../actor/claude-code-f69cd2cd.md)

## result

PASS

## statement

At 38638bb20, go test -trimpath ./internal/api -run 'TestShipCheck|TestShipGate' passes, including the new TestShipCheckAndGateFailOnOwnStaleBase (the shipping change set's own stale update fails ship check and ship gate --model-conflict with exit 1, the detail and human output naming the node ID and field:statement, recorded in the ship state so ship check --model-conflict exits 0; plain ship gate and --settled unaffected) and TestShipCheckAndGateFailOnOwnModelConflict (only another in-flight change set stale: check and gate pass). go test ./internal/workflow -run TestShipChange passes: pr-green fails to pr-model-conflict-triage, which passes to pr-model-conflict-fix; gate-land runs ship gate --model-conflict and fails to gate-model-conflict-triage, which passes to model-conflict-fix. Testscripts run-ship-model-conflict-fix (a stale own update fails pr-green, triage routes it to pr-model-conflict-fix, the fixer repairs it, and the run signs off and lands; a fixer that does not repair it fails on the until check), run-ship-change-fail, run-ship-ci-fix, run-ship-ci-lint-fix, run-ship-conflict-fix and describe-golden pass. The diff does not touch rhiza check. The .rhiza and builtin ship-change.toml changes match. make self-check fails only on stale rendered files, which is expected before the run's render step.

## History

- `01M77N83F027XBGGZPK9SV6YN2` Fix watchdog unmerged PR at step land-triage: run 01M77DWNRDSFY0PEP0S1P58P5W; log: PR green at 52ace47 and unmerged after step land-triage failed: step "land-triage" ended with outcome "fail" | exit 1 | No head the local gate failed on after a rebase is on the branch.
