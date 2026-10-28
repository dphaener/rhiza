# ship-needs-leaves-qa-4

evidence · `01M61CDD4CAPZEC086VA7P4HZY` · actor [claude-code-bae43f5f](../actor/claude-code-bae43f5f.md)

## result

PASS

## statement

Built this branch's binary (go build ./cmd/rhiza), fresh rhiza init repo with a bare origin, only rex/claude/gh faked, built-in build/review/ship-change unmodified. rhiza run start <leafless accepted change set> --permission-mode manual exited 3 with 'has no leaves: ship-change builds and reviews its leaves, so plan it first' and next 'rhiza run start --workflow plan'; no run directory was added and no worktree was made.

## Links

### Out

- verifies [ship-needs-leaves](../capability/ship-through-rhiza/ship-needs-leaves.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
