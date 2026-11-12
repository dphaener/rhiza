# done-own-command-guide-qa-1

evidence · `01M7859XNPXWQJCHN2Z5XHYAPR` · actor [claude-code-c8877637](../actor/claude-code-c8877637.md)

## result

PASS

## statement

QA cycle 1 with the binary built at HEAD c345f85e6: `rhiza guide build` and `rhiza guide review` both say each role ends with `rhiza session done` as its own command after its recording command succeeded. `rhiza workflow describe <w> --json` for build, review, plan, plan-milestone and ship-change: no prompt contains '; rhiza session done' or '&& rhiza session done', and every mention of rhiza session done is matched by 'very last action' (2/2, 8/8, 5/5, 5/5, 16/16). grep over .rhiza/workflows and internal/ finds no chained done. go test -trimpath ./internal/guide -run 'Build|Review' and ./internal/workflow pass.

## Links

### Out

- verifies [done-own-command-guide](../capability/rhiza/guide/done-own-command-guide.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
