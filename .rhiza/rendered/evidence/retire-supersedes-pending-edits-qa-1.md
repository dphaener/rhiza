# retire-supersedes-pending-edits-qa-1

evidence · `01M6XT122X93K62G6NG7NYPYFG` · actor [claude-code-2d293fb3](../actor/claude-code-2d293fb3.md)

## result

PASS

## statement

QA at 27e17a09 with the built binary in scratch git repos (RHIZA_* unset). New binary: set --handle store2 + statement, mv --parent, then retire store2 left one node.retire on the landed base; propose/accept/land exited 0. Old binary (merge-base) wrote set + mv + retire all on the landed base; land --json exited 3 with {"error":{"code":"conflict",...stale...}}; the new binary's retire <node-id> exited 0, dropped the update and move, left one retire on the landed base, and land exited 0. A second retire on a lone retire on the landed base is refused (exit 3, 'is retired'), as is retire of a node already retired on main; retire of a node created in the change set still drops it entirely. Retiring by handle after retire fails (handle no longer resolves), so the repair needs the node ID, as the PR 383 note gives. go test -trimpath ./internal/api -run 'TestRetire|Ship' ok; TestScripts/retire-after-update ok in a clean env.

## Links

### Out

- verifies [retire-supersedes-pending-edits](../capability/rhiza/change-sets/conflicts/retire-supersedes-pending-edits.md)

## History

- `01M6XHCX4FRER8CR0MG5D4D76D` Fix watchdog unmerged PR at step land: run 01M6X2EMMHKMQ9W4TDJH3ERXSD; log: PR green at 6d6c4d1 and unmerged after step land failed: step "land" ended with outcome "fail" | exit 1 | rhiza: ship: land with the PR head's build: /home/rhiza/code/rhiza/.git/rhiza/ship/01M6WMCCRC0KSCWEV17VY8KVQP/rhiza land --change 01M6WMCCRC0KSCWEV17VY8KVQP --json: exit status 3
