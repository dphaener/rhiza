# done-refuses-unrecorded-verifier-qa-1

evidence · `01M7859XKJHQ52Y988FPY9FKE0` · actor [claude-code-c8877637](../actor/claude-code-c8877637.md)

## result

PASS

## statement

QA cycle 1 with the binary built at HEAD c345f85e6, in a scratch repo (rhiza init, change set with capability leaf-a, actors ada and bob) with a hand-written build run journal (leaf-session builder and verifier on leaf-a attempt 1). Builder session done succeeds. Verifier with nothing recorded: exit 3 conflict naming `rhiza run findings --leaf leaf-a -`, session status stays unset (no done mark). After bob writes Evidence verifying leaf-a via rhiza apply: still refused, message says leaf-a-bob was written by another actor. After a findings entry on attempt 1: a verifier of attempt 1 says done. A verifier of attempt 2 is refused until ada writes her own Evidence, then says done. Human text form prints 'Session ... is done.'. go test -trimpath ./internal/api -run SessionDone passes.

## Links

### Out

- verifies [done-refuses-unrecorded-verifier](../capability/rhiza/workflows/workflow-runs/agent-step-done-signal/done-refuses-unrecorded-verifier.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
