# evidence-names-actor

capability · `01M5WR40B08FCCSMTK2QFHHPPM` · parent [evidence-actor](../evidence-actor.md)

## acceptance

A testscript writes Evidence with RHIZA_ACTOR=verifier, then sets a field on it as another actor: show --json still gives actor verifier, and the rendered page under .rhiza/rendered/ names verifier.

## statement

An Evidence node's actor is the actor stamped on the log entry that created it: rhiza show --json gives it as actor on every Evidence node, and the rendered Evidence page names it. No field can set or change it.

## Links

### Out

- decided-by [verifier-writes-evidence](../../../../decision/verifier-writes-evidence.md)
- decided-by [evidence-actor-is-creator](../../../../decision/evidence-actor-is-creator.md)

### In

- depends-on [check-evidence-actor](check-evidence-actor.md)
- verifies [evidence-names-actor-verified](../../../../evidence/evidence-names-actor-verified.md)
- verifies [pr-36-leaf-dependencies-and](../../../../evidence/pr-36-leaf-dependencies-and.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
