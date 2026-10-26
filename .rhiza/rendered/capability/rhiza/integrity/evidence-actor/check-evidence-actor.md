# check-evidence-actor

capability · `01M5WR40W9BDXA9Z2Q0T4GEM4S` · parent [evidence-actor](../evidence-actor.md)

## acceptance

A testscript in a git repo commits with Change: <id> and Actor: builder, then writes Evidence verifying the change set's leaf: as builder, check --json exits 0 with one warn finding on rule evidence-actor naming the four; as verifier, no such finding; a commit with no Actor: trailer yields none. Go tests run the rule on given commit trailers, so check itself still runs no git.

## statement

rhiza check warns, on a rule named evidence-actor, for each live Evidence that verifies a capability when the Evidence's actor is named in an Actor: trailer of a commit that changes a file outside .rhiza/ and carries the Change: trailer of the change set that added that capability, unless the commit's Leaf: trailers name other leaves only. The warning names the Evidence, the capability, the actor and the commit. Outside a git work tree the rule is skipped, silently.

## Links

### Out

- depends-on [evidence-names-actor](evidence-names-actor.md)
- depends-on [commit-actor-trailer](commit-actor-trailer.md)
- decided-by [verifier-writes-evidence](../../../../decision/verifier-writes-evidence.md)
- decided-by [evidence-actor-leaf-scope](../../../../decision/evidence-actor-leaf-scope.md)
- decided-by [evidence-actor-silent-skip](../../../../decision/evidence-actor-silent-skip.md)
- decided-by [evidence-actor-one-per-pair](../../../../decision/evidence-actor-one-per-pair.md)
- decided-by [evidence-actor-code-commits](../../../../decision/evidence-actor-code-commits.md)

### In

- relates-to [leaf-commit-scope](../../../../question/leaf-commit-scope.md) — Until it is settled the warning works per change set: an actor that committed any of the change set's code is warned about as the verifier of any of its leaves.
- verifies [check-evidence-actor-verified](../../../../evidence/check-evidence-actor-verified.md)
- verifies [pr-36-leaf-dependencies-and](../../../../evidence/pr-36-leaf-dependencies-and.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
