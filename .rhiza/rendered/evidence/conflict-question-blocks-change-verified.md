# conflict-question-blocks-change-verified

evidence · `01M6D298FZWP1DZJHDM8TSFBTT` · actor [claude-code-a0909b6f](../actor/claude-code-a0909b6f.md)

## result

PASS

## statement

At 9394ffd1 ran go test ./cmd/rhiza -run TestScripts/run-ship-change-rebase -count=1 -testwork: PASS. Case 2 (code conflict on feature.txt) asserts rhiza show rebase-conflict-ship-two --json has a blocks link, rhiza inbox --json lists it with kind blocker and blocks 'this change set', a second rhiza ship rebase leaves the .rhiza/log byte count unchanged (cmp) while still reporting the question, and rhiza question answer rebase-conflict-ship-two 'Keep main' --change <cs> succeeds and folds. In the kept work dir, rhiza show on the question listed 'blocks change set 01M5PWGEF85Y18: Ship two' and 'answered by answer-rebase-conflict-ship-two', so the answer was written. Read the diff of internal/api/ship_rebase.go: new questions get add+blocks link; an existing identical question gains the link only when LinkBetween reports not-found, else nothing is written. make check passes (fmt, vet, all tests, self-check run with the RHIZA_* run env unset; with the env set, self-check only flagged this verifier's own unrendered actor record).

## Links

### Out

- verifies conflict-question-blocks-change (retired)

## History

- `01M6D12NNWQZHEWFBX60JS04AG` A run's rebase-conflict question blocks its change set, so the inbox and question answer see it
