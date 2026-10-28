# ship-open-pr-verified

evidence · `01M6093X354JKP7XW41A8ACTG7` · actor [claude-code-5bb600f3](../actor/claude-code-5bb600f3.md)

## result

PASS

## statement

At bc5ec2b, ran go test ./internal/api -run TestShipOpenPR -count=1 -v: TestShipOpenPR and TestShipOpenPRRefusals pass. With a bare git origin and a fake gh, the first call pushes cs/ship-it (origin head == HEAD) and makes exactly one pr create titled 'Ship it' (the change set's title) with a body starting '# Review of run R1' holding '### alpha: passed'; after a new commit and journal entry, the second call pushes the new head, edits PR #7's body (now holding 'step code: pass') and makes no second create; the ship state file records PR 7, URL, branch and pushed head; ship.check with pr_open then reads PR 7 and passes. With no review journaled the call exits 2 naming 'no review journaled', and with RHIZA_RUN unset exits 2 naming RHIZA_RUN, with no gh call and no branch on origin. The built CLI's 'rhiza ship open-pr' outside a run also exits 2 naming RHIZA_RUN. make test passes across all packages; fmt-check and vet pass; rhiza check on a clean clone at the commit with a clean env reports ok (the in-session self-check failures name only actors from this verifier session's env, not the commit). Read the diff: ship_openpr.go, its tests, describe-golden and rendered actor files.

## Links

### Out

- verifies [ship-open-pr](../capability/ship-through-rhiza/ship-open-pr.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
