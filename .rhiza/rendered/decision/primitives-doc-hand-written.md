# primitives-doc-hand-written

decision · `01M7E5H477ZR8V8PKZ06CS0RWN`

## question

Where does the hand-written primitives reference go? (docs/reference/ holds only files `rhiza docs` generates.)

## reasoning

docs/reference/ holds only files rhiza docs generates and rhiza check compares; the primitives page explains how step kinds and commands combine, which is prose rather than a key table, so it sits beside custom-workflows.md and links workflow-file.md for the keys instead of duplicating them.

## statement

The primitives reference is a hand-written docs/primitives.md: each primitive (agent, spawn, gate, push/pr, wait, merge, signoff, hook) with the step kind or command that is it, its outcomes and how a red hands back; it links the generated workflow-file reference for the keys and is linked from the README, runs-and-daemon.md and custom-workflows.md.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [byo-docs-planning](../source/byo-docs-planning.md)

### In

- decided-by [docs-primitives](../capability/rhiza/documentation/docs-primitives.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
