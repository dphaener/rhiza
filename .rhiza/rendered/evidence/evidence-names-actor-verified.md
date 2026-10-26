# evidence-names-actor-verified

evidence · `01M5X0Z374E5EB0F8RWXYXTPF0` · actor [claude-code-b6232fc2-verifier](../actor/claude-code-b6232fc2-verifier.md)

## result

PASS under make check at commit a345b25 on branch cs/leaf-dependencies-and, run by an independent verifier (claude-code-b6232fc2-verifier)

## statement

Evidence written with RHIZA_ACTOR=verifier and then set by another actor (including a field named actor) still shows actor verifier in show --json and show, in the draft what-if view too, and its rendered page under .rhiza/rendered/ names verifier linked to the actor page; other kinds carry no actor.

## tests

testscript evidence-actor-show; go test internal/changeset: TestFoldHistoryActor; internal/model: TestCreatedBy; internal/render: TestEvidenceActor, TestGolden

## Links

### Out

- verifies [evidence-names-actor](../capability/rhiza/integrity/evidence-actor/evidence-names-actor.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
