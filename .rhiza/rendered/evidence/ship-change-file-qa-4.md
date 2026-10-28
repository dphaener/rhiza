# ship-change-file-qa-4

evidence · `01M61CDCQTE51GSFPVKVHST9YG` · actor [claude-code-bae43f5f](../actor/claude-code-bae43f5f.md)

## result

FAIL

## statement

Built this branch's binary (go build ./cmd/rhiza), fresh rhiza init repo with a bare origin, only rex/claude/gh faked, built-in build/review/ship-change unmodified. Two-leaf change set (second depends on first): rhiza run start --permission-mode manual ran build/build, build/check, render, review/code, review/review (5 of 5 checks), open-pr, wait-pr, pr-green and stopped at signoff; a fresh clone of the pushed cs/ship-it passed rhiza check; rhiza run approve ran rebase through land, merge, main-green, close-out and ship, change set landed, main passed rhiza check, no shipit anywhere. But after rhiza run request-changes (one finding on a leaf) the run goes back to review/review without rendering again, and a fresh clone of the re-pushed branch fails rhiza check: missing rendered/actor page for the cycle-2 verifier, index.md differs.

## Links

### Out

- verifies [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
