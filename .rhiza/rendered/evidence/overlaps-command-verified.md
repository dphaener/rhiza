# overlaps-command-verified

evidence · `01M6AE89REGNM0TVNFQ3WR28HY` · actor [claude-code-cb168241](../actor/claude-code-cb168241.md)

## result

PASS

## statement

At 9b04e081 I ran go test ./internal/model -run Overlap -v: TestOverlapsNamedDraft passes and checks that a draft counts only when named (no draft named gives nil, draft A named gives one A/B overlap on field:title, naming a proposed or absent set gives nil). go test ./cmd/rhiza -run 'TestScripts/(overlaps|describe-golden)$' passes. The new overlaps.txtar has a proposed change set B and a draft C both setting field statement of base. It checks 'No overlaps' and an empty JSON list without --change, then with --change C the text line 'change sets B and C both write field:statement of base' and the JSON a/b/target/aspects. describe-golden.txtar includes the overlaps operation and its cmp passes, so it is regenerated. I read the diff of internal/api/overlaps.go and model.Overlaps(sets, drafts...). make check exits 0; one earlier run failed on a timing test that later passed.

## Links

### Out

- verifies [overlaps-command](../capability/rhiza/views/overlaps-command.md)

## History

- `01M5SWTMR4G5EFRVRQHSEKVA15` Genesis reviews each plan against the model and the other runs in flight before work starts, as an agent check on the plan rather than the code.
