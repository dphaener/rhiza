# workflow-step-qa-1

evidence · `01M60B4SSPGM5J4GV0SFS0RC7G` · actor [claude-code-f5af7e3e](../actor/claude-code-f5af7e3e.md)

## result

PASS

## statement

QA cycle 1, binary built from branch head 20f3cee, run by hand in scratch repos. With repo override workflows, `rhiza run start --workflow X` exited 2, with no run journaled, for: a workflow step with no workflow; one naming nosuch; self -> self; selfa -> selfb -> selfa; skip naming a step inner lacks (listing one, three, two); and skip or workflow on a command step. An outer workflow with step w over a three-step inner (skip=[two]), inner/three failing on a flag file: the run failed at inner/three, journaled inner/one then inner/three, and two never ran. With the flag removed, `rhiza run resume <run> --retry` re-ran only inner/three, then after, and the run succeeded (order log: one, three, after).

## Links

### Out

- verifies [workflow-step](../capability/ship-through-rhiza/workflow-step.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
