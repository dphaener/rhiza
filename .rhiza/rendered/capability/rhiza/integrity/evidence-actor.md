# evidence-actor

capability · `01M5WR40418Q8WECTQVWKZZ4QC` · parent [integrity](../integrity.md)

## acceptance

A testscript in a git repo has a builder commit a leaf's code and then write Evidence that verifies it: rhiza check --json exits 0 with one evidence-actor warning; written by a different actor, there is none.

## statement

Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the code of the leaf the Evidence verifies.

## Links

### Out

- decided-by [verifier-writes-evidence](../../../decision/verifier-writes-evidence.md)

### In

- verifies [pr-36-leaf-dependencies-and](../../../evidence/pr-36-leaf-dependencies-and.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
