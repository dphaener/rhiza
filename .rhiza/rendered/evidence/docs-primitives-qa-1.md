# docs-primitives-qa-1

evidence · `01M7EAC1YY3X0TRXCC5EKPTAA1` · actor [claude-code-0cae456c](../actor/claude-code-0cae456c.md)

## result

PASS

## statement

QA: read docs/primitives.md; it has sections for agent, spawn, gate, push/pr, wait, merge, sign-off and hook, the handback section, all nine step kinds and links reference/workflow-file.md; README.md, docs/runs-and-daemon.md and docs/custom-workflows.md link it. The flags it names (run spawn --role/--leaf/--brief/--model/--check, ship check --settled/--no-rebase/--main, ship rebase --no-push) exist in the scratch build's help. TestDocsPrimitives and TestScripts/docs-commands-exist pass; make self-check passes in a clean detached worktree at HEAD.

## Links

### Out

- verifies [docs-primitives](../capability/rhiza/documentation/docs-primitives.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
