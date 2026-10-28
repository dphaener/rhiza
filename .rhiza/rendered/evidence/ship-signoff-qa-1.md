# ship-signoff-qa-1

evidence · `01M60B4TVT0DNRXK3H15EPVDPB` · actor [claude-code-f5af7e3e](../actor/claude-code-f5af7e3e.md)

## result

PASS

## statement

QA cycle 1, binary built from branch head 20f3cee. Built-in ship-change, with review.toml overridden without its code step and fake rex/claude/gh. The run stopped at signoff after pr-green, with run status showing Current step: signoff and the PR URL, and no pr merge. `rhiza run approve` ran rebase through ship, and the change set landed with one `pr merge 7 --squash --match-head-commit`. `rhiza run request-changes` with one major on leaf fourth ran a fixer, then that leaf's verifier (cycle 2), and re-ran only check alpha (cycle 2), journaled 'signoff (visit 1): changes, then review/review'. review/review, open-pr, wait-pr and pr-green ran as visit 2, with one pr edit and no second create, and the run waited at signoff again, then landed on approve. With --autonomous, signoff passed 'approved by rule autonomous' and the run landed. approve and request-changes exit 2 on a finished run.

## Links

### Out

- verifies [ship-signoff](../capability/ship-through-rhiza/ship-signoff.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
