# ship-signoff-qa-4

evidence · `01M61CDDW1SARHN0V250H79GWN` · actor [claude-code-bae43f5f](../actor/claude-code-bae43f5f.md)

## result

PASS

## statement

Built this branch's binary (go build ./cmd/rhiza), fresh rhiza init repo with a bare origin, only rex/claude/gh faked, built-in build/review/ship-change unmodified. ship-change stopped at signoff after green CI with run status naming the PR URL; rhiza run approve carried on to rebase, land and merge. rhiza run request-changes with one finding (check qa, leaf thing2) ran a fixer then the leaf's verifier, re-ran only check qa as cycle 2 under review/review, then open-pr (pr edit), wait-pr, pr-green and stopped at signoff again.

## Links

### Out

- verifies [ship-signoff](../capability/ship-through-rhiza/ship-signoff.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
