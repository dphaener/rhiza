# ship-gate-fails-own-model-conflict-qa-1

evidence · `01M77SJEJGVN714VQXBX2H3T7V` · actor [claude-code-e1432792](../actor/claude-code-e1432792.md)

## result

PASS

## statement

QA by hand with the binary built at 668f5cf08 in a scratch repo where the accepted change set's own update of first rests on a base the landed change set Other rewrote: rhiza check only warned (exit 0); rhiza ship gate --model-conflict exited 1 with 'Gate failed: the change set's own operations rest on a stale base: <node> (field:statement) rewritten by op ... Restate them on the landed version (rhiza set <ref> --rebase).' and --json model_conflicts naming op, node, fields, base, landed op and change set; rhiza ship check (fake gh: PR #7 green) exited 1 with detail 'failed: 1 operation(s) of the change set rest on a stale base: ...' and recorded the conflict in .git/rhiza/ship/<cs>.json. After rhiza set first --rebase --field 'statement=Theirs and ours', ship gate --model-conflict passed (exit 0, no gate set) and ship check --model-conflict exited 1. Scoped go test ./internal/api -run 'TestShipCheck|TestShipGate|TestShipLand' ok.

## History

- `01M77N83F027XBGGZPK9SV6YN2` Fix watchdog unmerged PR at step land-triage: run 01M77DWNRDSFY0PEP0S1P58P5W; log: PR green at 52ace47 and unmerged after step land-triage failed: step "land-triage" ended with outcome "fail" | exit 1 | No head the local gate failed on after a rebase is on the branch.
