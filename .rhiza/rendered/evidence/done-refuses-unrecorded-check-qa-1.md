# done-refuses-unrecorded-check-qa-1

evidence · `01M7859XJ15MZKY2A0F65ATGA4` · actor [claude-code-c8877637](../actor/claude-code-c8877637.md)

## result

PASS

## statement

QA cycle 1 with the binary built at HEAD c345f85e6. Live: this QA check session's own `rhiza session done --dry-run --json` (check qa, cycle 1, no round yet) exits 3, code conflict, message names `rhiza run findings --check qa -`, next `rhiza run findings --check qa`. Scratch repo with a hand-written review run journal (step checks, check bugs, cycle 1): check session refused (exit 3) naming rhiza run findings and written no done mark; after a findings round (blocker f1, major f2, nit f3) it says done. Critic with f1,f2 unjudged refused naming `rhiza run verdict --check bugs -` and listing (f1, f2); after f1 confirmed still refused listing (f2); after f2 dropped it says done with nit f3 unjudged. Fixer session says done with no round. go test -trimpath ./internal/api -run SessionDone passes.

## Links

### Out

- verifies [done-refuses-unrecorded-check](../capability/rhiza/workflows/workflow-runs/agent-step-done-signal/done-refuses-unrecorded-check.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
