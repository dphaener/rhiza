# ship-needs-leaves-qa-5

evidence · `01M61FJ7AG0WNHFAEC5R8QRGRV` · actor [claude-code-0f55b0ee](../actor/claude-code-0f55b0ee.md)

## result

PASS

## statement

Built this branch's binary at 731d7f7 (go build ./cmd/rhiza) and drove it in a fresh rhiza init repo with a bare origin, only rex/claude/gh faked, built-in build/review/ship-change unmodified, on a two-leaf change set (second depends-on first). rhiza run start on an accepted change set with no leaves exited 3 with 'has no leaves: ship-change builds and reviews its leaves, so plan it first' and next 'rhiza run start --workflow plan', no new run journal and no worktree made; the same command on the two-leaf change set started the run.

## Links

### Out

- verifies [ship-needs-leaves](../capability/ship-through-rhiza/ship-needs-leaves.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
