# done-own-command-guide-verified

evidence · `01M784GTHRT3TCKZ4THQKJYJGD` · actor [claude-code-b092a740](../actor/claude-code-b092a740.md)

## result

PASS

## statement

At c44cd76a7: go test -trimpath -count=1 ./internal/guide -run 'Build|Review' -v passes; TestBuildText and TestReviewText now require 'rhiza session done', 'as its own command' and '`;` or `&&`', and build.md/review.md both say each role ends with rhiza session done as its own command after its recording command succeeded, never chained with ; or &&. go test -trimpath -count=1 ./internal/workflow passes, including new TestBuiltinPromptsDoneOwnCommand, which scans every builtin (BuiltinNames) prompt naming rhiza session done for '; rhiza session done', ';rhiza session done', '&& rhiza session done' and requires 'very last action'. Read the diff: all prompts in build, plan, plan-milestone, review and ship-change .toml are single-line, and grep shows no session done mention outside prompt lines, so the scan covers them all.

## Links

### Out

- verifies [done-own-command-guide](../capability/rhiza/guide/done-own-command-guide.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
