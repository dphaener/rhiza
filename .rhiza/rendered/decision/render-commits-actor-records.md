# render-commits-actor-records

decision · `01M61BG9ECFWH96DYYPCKWV6KB`

## reasoning

QA's cycle-3 finding was confirmed by its critic. ship-change.toml's render step stages only .rhiza/rendered. Builder actor records (recordAgentActor) stay uncommitted until land, so the PR head has rendered actor pages with no records, and this repo's CI make check fails before sign-off. Staging .rhiza/log in the same step is the smallest fix, and it matches what land already does (git add .rhiza). Changing the leaf instead would leave ship-change unable to ship in Rhiza's own repo, which is the change set's intent.

## statement

Ship-change's render step commits the worktree's uncommitted .rhiza/log entries together with .rhiza/rendered, so the branch it pushes holds every actor record its rendered pages are built from and rhiza check passes on a clone of it. A test proves it: a builder that commits only its own files, then rhiza check on a fresh clone of the pushed branch passes. Leaf ship-change-file keeps its scope and its statement names this.

## wording

decided by the orchestrator acting as arbiter, on the owner's word that the orchestrator signs off on this change set

## Links

### Out

- answers [leaf-failed-01m60atz-1](../question/leaf-failed-01m60atz-1.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
