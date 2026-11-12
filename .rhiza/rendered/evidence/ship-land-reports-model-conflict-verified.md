# ship-land-reports-model-conflict-verified

evidence · `01M77R3RTSAGJMDTSAG13SV0PR` · actor [claude-code-3ffdb309](../actor/claude-code-3ffdb309.md)

## result

PASS

## statement

Verifier, detached checkout at cc2a6d8de: go test -trimpath -count=1 ./internal/api -run 'TestShipLand' -v passes all 8 cases, including the new TestShipLandModelConflict. It sets up a fixture change set whose set on node 'checking' rests on a base that a landed change set has since rewritten. ship.land returns exit code 1 with Landed=false and one model_conflicts entry naming the node ID, fields [field:statement], base, landed op and the other change set; the human text says 'stale base'; the state file records the same conflict. ship check --model-conflict exits 0 while the conflict holds, exits 1 after a rebase repair on the landed version, and exits 1 for a state with no recorded conflict (another land failure). The diff refuses before Evidence or any write, and clears the record on a successful push. Also passed: go test ./internal/model (typed ConflictError still matches ErrConflict), go test ./internal/api -run 'TestShip|Describe|Schema' (all tests pass; the package's root-state guard flagged ~/code/.rhiza-root/usage.live.json changing, which this usage-free diff does not touch), go test ./cmd/rhiza -run 'TestScripts/(describe-golden|ship-.*)$' ok, gofmt clean, go vet clean.

## History

- `01M77N83F027XBGGZPK9SV6YN2` Fix watchdog unmerged PR at step land-triage: run 01M77DWNRDSFY0PEP0S1P58P5W; log: PR green at 52ace47 and unmerged after step land-triage failed: step "land-triage" ended with outcome "fail" | exit 1 | No head the local gate failed on after a rebase is on the branch.
