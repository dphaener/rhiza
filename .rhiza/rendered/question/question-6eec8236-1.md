# question-6eec8236-1

question · `01M77NDXVWW54D78AQT3HRKZKZ`

## session

6eec8236-ad63-49dd-824c-01a1cf546583

## statement

Run 01M77DW's land step was refused because change set 01M763Z4's update to ship-check-ci-test-findings' statement rests on a stale base: 01M76B50 (landed via #455) has since rewritten that field. `rhiza check` only warns about an in-flight change set's conflict, so every check passed. land-triage only handles a failed pre-push check, so the run ended with green PR #480 unmerged. How should ship-change handle a land refused on a model conflict?

## Links

### In

- answers [ship-land-model-conflict-repair](../decision/ship-land-model-conflict-repair.md)

## History

- `01M77N83F027XBGGZPK9SV6YN2` Fix watchdog unmerged PR at step land-triage: run 01M77DWNRDSFY0PEP0S1P58P5W; log: PR green at 52ace47 and unmerged after step land-triage failed: step "land-triage" ended with outcome "fail" | exit 1 | No head the local gate failed on after a rebase is on the branch.
