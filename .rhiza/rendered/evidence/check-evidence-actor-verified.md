# check-evidence-actor-verified

evidence · `01M5X0Z3Q26XEPV9ZZ73ASK7ZJ`

## result

PASS under make check at commit 92051c5 on branch cs/leaf-dependencies-and, run by an independent verifier (claude-code-b6232fc2-verifier)

## statement

In a git repo with a commit carrying Change: and Actor: builder that changes a file outside .rhiza/, Evidence written by builder verifying the change set's leaf gives check --json ok with one warn finding on rule evidence-actor naming the evidence, capability, actor and commit; written by verifier it gives none; a commit with no Actor: trailer, or a Leaf: trailer naming another leaf only, gives none; a commit with the builder's trailers that changes only .rhiza/ files, or an empty commit, gives none, and a later code commit with the same trailers gives one; outside a git work tree the rule is skipped. The Go tests run the rule on given commit trailers and skip it with no git; TestGitCommits checks that git log is asked only for commits that change a file outside .rhiza/ (pathspec :/ :(exclude).rhiza).

## tests

testscript check-evidence-actor (builder, verifier, trailers, modelonly, nogit); go test internal/check: TestRunEvidenceActor; internal/api: TestParseTrailerLog, TestGitCommits

## Links

### Out

- verifies [check-evidence-actor](../capability/rhiza/integrity/evidence-actor/check-evidence-actor.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
