# ship-open-pr-qa-4

evidence · `01M61CDDG78V91RJ2NGE2T0ZAR` · actor [claude-code-bae43f5f](../actor/claude-code-bae43f5f.md)

## result

PASS

## statement

Built this branch's binary (go build ./cmd/rhiza), fresh rhiza init repo with a bare origin, only rex/claude/gh faked, built-in build/review/ship-change unmodified. Inside the ship-change run, open-pr pushed cs/ship-it and ran gh pr create --title 'Ship it' with the run summary as body ('# Review of run ...'); after request-changes the second open-pr pushed the fixer's commits and ran gh pr edit, with exactly one pr create; run status named the PR URL. Outside a run, rhiza ship open-pr exited 2 saying RHIZA_RUN is not set.

## Links

### Out

- verifies [ship-open-pr](../capability/ship-through-rhiza/ship-open-pr.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
