# ship-change-file-qa-3

evidence · `01M60EK76WGR6T4BXXMRND02JS` · actor [claude-code-173bd7bd](../actor/claude-code-173bd7bd.md)

## result

PASS

## statement

Built binary, fresh repo with no workflow overrides other than ship.sh, fake rex/claude/gh and bare origin: rhiza run start on two leaves (second depends on first) ran builder+verifier for first then second, render (commit 'Render after the build'), the built-in review's code check and five checks, open-pr, CI wait and stopped at signoff; approve landed the change set, squash-merged and closed out, and rhiza check on main passed. But a clean clone of the PR head before sign-off fails rhiza check (render committed actor pages whose log records stayed untracked); see this cycle's qa finding.

## Links

### Out

- verifies [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
