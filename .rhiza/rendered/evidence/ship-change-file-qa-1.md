# ship-change-file-qa-1

evidence · `01M60B4V75A706S7VR8YBQP5KE` · actor [claude-code-f5af7e3e](../actor/claude-code-f5af7e3e.md)

## result

FAIL

## statement

QA cycle 1, binary built from branch head 20f3cee. ship-change.toml has no shipit and no --pr-open. Unmodified built-in ship-change on a fresh `rhiza init` repo, two leaves with second depends-on first, fake rex/claude/gh: build/build passed 2 of 2 leaves in order and build/check passed, then review/code FAILED. rhiza check found the render stale for the actors of build's sessions and for index.md (qa finding f1). With review.toml overridden without the code step, as run-ship-change.txtar does, the run went build, review/review, open-pr, wait-pr, pr-green and signoff, then after approve rebase, land, merge, main check, close-out (worktree and branches gone) and ship, and the change set landed. go test for run-ship*, review-signoff and describe-golden passes.

## Links

### Out

- verifies [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
