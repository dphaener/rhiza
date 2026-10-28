# ship-signoff-qa-5

evidence · `01M61FJ6JSPXGW78XBC7DDDT0K` · actor [claude-code-0f55b0ee](../actor/claude-code-0f55b0ee.md)

## result

PASS

## statement

Built this branch's binary at 731d7f7 (go build ./cmd/rhiza) and drove it in a fresh rhiza init repo with a bare origin, only rex/claude/gh faked, built-in build/review/ship-change unmodified, on a two-leaf change set (second depends-on first). The run stopped at signoff after green CI, run status naming https://github.test/o/r/pull/7; request-changes with one finding on leaf second started that leaf's fixer (its prompt held the finding) and verifier, re-ran only the qa check as cycle 2 under review/review, called open-pr again (1 pr create, 1 pr edit), waited for CI and stopped at signoff visit 2; rhiza run approve then moved on to rebase and land and the run succeeded; a second rhiza run approve exited 2.

## Links

### Out

- verifies [ship-signoff](../capability/ship-through-rhiza/ship-signoff.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
