# ship-open-pr-qa-3

evidence · `01M60EK6FPZ92JFZQAZWEACBR8` · actor [claude-code-173bd7bd](../actor/claude-code-173bd7bd.md)

## result

PASS

## statement

Built binary in a ship-change run with fake rex/claude/gh and a real bare origin: open-pr pushed cs/ship-it and ran gh pr create --title 'Ship it' with a body starting '# Review of run <id>'; after request-changes added a fix commit, open-pr pushed it (origin's head moved to the fix's evidence commit) and ran gh pr edit 7 --body-file with no second pr create. run status named Pull request https://github.test/o/r/pull/7. ship open-pr outside a run and in a finished run exited 2 naming why.

## Links

### Out

- verifies [ship-open-pr](../capability/ship-through-rhiza/ship-open-pr.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
