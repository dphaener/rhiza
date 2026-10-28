# ship-needs-leaves-qa-3

evidence · `01M60EK5S34CH9QS6GMTF3CADQ` · actor [claude-code-173bd7bd](../actor/claude-code-173bd7bd.md)

## result

PASS

## statement

Built binary (go build at 1913524), fresh rhiza init repo with a bare origin: rhiza run start <cs> --permission-mode manual on an accepted change set with no leaves exited 3 with 'has no leaves: ship-change builds and reviews its leaves, so plan it first' and next 'rhiza run start --workflow plan'; rhiza run status then said no runs yet and git worktree list showed only the main checkout. After adding two leaves the same command started the run.

## Links

### Out

- verifies [ship-needs-leaves](../capability/ship-through-rhiza/ship-needs-leaves.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
