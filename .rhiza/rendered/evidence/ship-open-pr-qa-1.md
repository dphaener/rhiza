# ship-open-pr-qa-1

evidence · `01M60B4TGH977Z1Y9QZKZ8XZSJ` · actor [claude-code-f5af7e3e](../actor/claude-code-f5af7e3e.md)

## result

PASS

## statement

QA cycle 1, binary built from branch head 20f3cee. `rhiza ship open-pr` with RHIZA_RUN unset exits 2 naming RHIZA_RUN. Inside built-in ship-change runs (real bare origin, fake gh), open-pr pushed the cs/ branch and made one `gh pr create --base main --head cs/ship-it --title Ship it --body-file`. The body began '# Review of run' and held '## Checks of step review/review', '### alpha: passed' and '### beta: passed'. After request-changes, open-pr visit 2 made one pr edit and no second pr create. run status named the PR URL https://github.test/o/r/pull/7.

## Links

### Out

- verifies [ship-open-pr](../capability/ship-through-rhiza/ship-open-pr.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
