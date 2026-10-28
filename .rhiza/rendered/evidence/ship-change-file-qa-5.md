# ship-change-file-qa-5

evidence · `01M61FJ66ZW6ZRSDNQVNX2SN31` · actor [claude-code-0f55b0ee](../actor/claude-code-0f55b0ee.md)

## result

PASS

## statement

Built this branch's binary at 731d7f7 (go build ./cmd/rhiza) and drove it in a fresh rhiza init repo with a bare origin, only rex/claude/gh faked, built-in build/review/ship-change unmodified, on a two-leaf change set (second depends-on first). rhiza run start ran build/build, build/check, render, review/code, review/review (5 of 5 checks), render-pr, open-pr, wait-pr, pr-green and stopped at signoff; a fresh clone of the pushed cs/ship-it passed rhiza check. rhiza run request-changes (one qa finding on leaf second) ran the fixer and verifier on second, re-ran only qa as cycle 2, then render-pr (visit 2) committed 'Render after the review', open-pr pushed and edited the PR; a fresh clone of the re-pushed head passed rhiza check, its rendered/actor pages including the cycle-2 fixer, verifier and qa sessions (the cycle-4 failure is gone). rhiza run approve ran rebase through land, merge, main-green, close-out and ship; the run succeeded, the change set landed, main passed rhiza check, worktree and origin branch gone, no shipit in any output.

## Links

### Out

- verifies [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
