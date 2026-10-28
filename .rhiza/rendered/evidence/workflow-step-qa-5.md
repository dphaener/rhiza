# workflow-step-qa-5

evidence · `01M61FJ7P7T0H8FY3B6858W6BX` · actor [claude-code-0f55b0ee](../actor/claude-code-0f55b0ee.md)

## result

PASS

## statement

Built this branch's binary at 731d7f7 (go build ./cmd/rhiza) and drove it in a fresh rhiza init repo with a bare origin, only rex/claude/gh faked, built-in build/review/ship-change unmodified, on a two-leaf change set (second depends-on first). ship-change's build and review workflow steps ran their inner steps in the same run and journal (build/build, build/check, review/code, review/review), review's signoff skipped, the outer on.pass taken to render and render-pr; request-changes re-entered at review/review.

## Links

### Out

- verifies [workflow-step](../capability/ship-through-rhiza/workflow-step.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
