# answer-refuses-closed-or-unblocking-verified-2

evidence · `01M6B229Z5HYP3RPQKBMM1ADYS` · actor [claude-code-6062b3b2](../actor/claude-code-6062b3b2.md)

## result

PASS

## statement

Verifier, attempt 2 at 7214f794. Ran go test ./cmd/rhiza -run TestScripts/question-answer -count=1 -v: PASS. In the script, store-q and cache-q block store and cache, which change set Accepted touches; both are answered (accepted). Re-answering store-q exits 2 with code usage ('not an open question blocking this change set') and answer-store-q-2 does not exist. apart-q blocks apart; the new fix has change set Other touch apart (set apart statement=Apart, shown in change show). Answering apart-q with --change Accepted then exits 2 with code usage, and answer-apart-q exists in neither the current view nor what-if:Accepted. fold-ok passes. This closes compliance finding f1. make fmt-check vet test exits 0 (one earlier run exited 2 and passed on rerun with no FAIL lines). make self-check fails only on stale .rhiza/rendered actor/index files, which also fail at the parent commit e655f3a8, so the leaf did not cause them. The diff adds only test lines.

## Links

### Out

- verifies [answer-refuses-closed-or-unblocking](../capability/rhiza/change-sets/question-answer-op/answer-refuses-closed-or-unblocking.md)

## History

- `01M69VQG7GJARDSWZWKDJ4CWQ1` question.answer refuses a question that is not open or does not block the target change set, so an answer always lands in the change set the question blocks
