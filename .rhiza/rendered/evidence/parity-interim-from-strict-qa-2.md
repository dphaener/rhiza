# parity-interim-from-strict-qa-2

evidence · `01M7GWB8AZ3SKWRQPTNKMSVNSH` · actor [claude-code-e084a010](../actor/claude-code-e084a010.md)

## result

PASS

## statement

QA cycle 2 with a scratch build (go build -o scratch/rhiza-qa ./cmd/rhiza) in a scratch repo with two landed heads on x (same field) and y (disjoint fields), 5 proposals (leaves under landed par, under landed z, at a root, set-only): inbox --json lists the 4 leaf-bearing proposals, set-only none, identical to a 556194e52 build; health, roadmap, ready, slice, change show/diff, overlaps match base; on this repo's own log ready/roadmap/inbox --json are byte-identical to base; scripts daemon-picks, daemon-plans-proposals, inbox-uncovered-proposal, health pass.

## Links

### Out

- verifies [parity-interim-from-strict](../capability/rhiza/views/fast-reads/parity-interim-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
