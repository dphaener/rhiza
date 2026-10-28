# render-before-open-pr

decision · `01M61EWAEN7Z8H3075Z91EPN0P`

## reasoning

QA's cycle-4 finding was confirmed by its critic. After request-changes, signoff goes back to review/review, and the fix cycle's verifier and checks then commit Evidence and actor records that nothing renders before ship open-pr pushes. So the re-pushed PR fails rhiza check in CI, and the run fails at pr-green. A render step just before open-pr covers both paths with one rule. The answered cycle had failed, so the question was escalated to the owner, who picked this option.

## statement

Ship-change renders again right before every push. A second render command step sits between the review workflow step and open-pr (review, then render-pr, then open-pr). It runs rhiza render and commits .rhiza/rendered and the uncommitted .rhiza/log entries, like the render step after build. signoff's on.changes stays review/review, so the push after request-changes is rendered too. run-ship-signoff clones the head re-pushed after request-changes and runs rhiza check on it. Leaf ship-change-file's statement names both render steps.

## wording

the agent's option, picked by the owner: "Render before open-pr (Recommended)" (add a second render step between review and open-pr, so the first push and every push after request-changes is rendered and committed; extend run-ship-signoff to clone the re-pushed head and run rhiza check)

## Links

### Out

- answers [leaf-failed-01m60atz-4](../question/leaf-failed-01m60atz-4.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
