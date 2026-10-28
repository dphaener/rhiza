# ship-signoff-qa-3

evidence · `01M60EK6TYVS3HWDQ2H98GB6Y9` · actor [claude-code-173bd7bd](../actor/claude-code-173bd7bd.md)

## result

PASS

## statement

Built binary, fake rex/claude/gh: ship-change stopped at signoff after pr-green with run status naming the PR URL and next run summary/approve/request-changes. request-changes with one qa finding on leaf first ran the fixer and verifier of first, then only qa as cycle 2 of review/review, then open-pr (pr edit), wait-pr, pr-green and signoff again. run approve then went on through rebase, land, merge (pr merge 7 --squash --match-head-commit), close-out and ship; the change set landed. approve and request-changes on the finished run both exited 2.

## Links

### Out

- verifies [ship-signoff](../capability/ship-through-rhiza/ship-signoff.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
