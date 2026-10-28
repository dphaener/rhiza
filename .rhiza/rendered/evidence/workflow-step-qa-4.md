# workflow-step-qa-4

evidence · `01M61CDE7W547CNBFQ0GAPF4AT` · actor [claude-code-bae43f5f](../actor/claude-code-bae43f5f.md)

## result

PASS

## statement

Built this branch's binary (go build ./cmd/rhiza), fresh rhiza init repo with a bare origin, only rex/claude/gh faked, built-in build/review/ship-change unmodified. ship-change's workflow steps ran the built-in build and review inside one run and journal: steps journaled as build/build, build/check, review/code, review/review in order; review's signoff was skipped (the run went from review/review to open-pr); both workflow steps passed and the outer on.pass steps followed.

## Links

### Out

- verifies [workflow-step](../capability/ship-through-rhiza/workflow-step.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
