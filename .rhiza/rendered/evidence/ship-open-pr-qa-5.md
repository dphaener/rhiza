# ship-open-pr-qa-5

evidence · `01M61FJ6YPNZ261XQHZZV95X4C` · actor [claude-code-0f55b0ee](../actor/claude-code-0f55b0ee.md)

## result

PASS

## statement

Built this branch's binary at 731d7f7 (go build ./cmd/rhiza) and drove it in a fresh rhiza init repo with a bare origin, only rex/claude/gh faked, built-in build/review/ship-change unmodified, on a two-leaf change set (second depends-on first). In the ship-change run, open-pr pushed cs/ship-it and ran gh pr create once with title 'Ship it' and a body holding the run summary; after request-changes the second open-pr pushed the new head and ran gh pr edit, no second create; pr-green's rhiza ship check read the recorded PR. rhiza ship open-pr outside a run exited 2 naming that RHIZA_RUN is not set.

## Links

### Out

- verifies [ship-open-pr](../capability/ship-through-rhiza/ship-open-pr.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
