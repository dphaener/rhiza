# ship-change-file-verified-3

evidence · `01M60C4T5PDY6QNCQ56XDZ1RSE` · actor [claude-code-8ed323e1](../actor/claude-code-8ed323e1.md)

## result

PASS

## statement

Verifier attempt 3 at 179a96f: go test -count=1 ./cmd/rhiza -run 'TestScripts/(run-ship-change|describe-golden)' passes. run-ship-change.txtar (fake rex, claude, gh) runs rhiza run start on a two-leaf change set (second depends-on first) and checks: builder then verifier for first, then for second; render, review/code and two review checks passing; the PR opened once with the review summary as its body; the run stopped at signoff until rhiza run approve; then rebase, land, squash merge pinned to the landed head, main-green, close-out and the ship command; the change set landed. TestShipChangeBuiltin passes: no shipit or --pr-open in ship-change.toml, review/signoff skipped, step order as stated plus a render step between build and review. describe-golden passes, so the golden is current. Read the diffs of 603521f, 7c7f59d and 179a96f. make fmt-check, vet and test pass. make self-check fails only on this verifier session's own unrendered actor entry, which is not in the leaf's tree.

## Links

### Out

- verifies [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
