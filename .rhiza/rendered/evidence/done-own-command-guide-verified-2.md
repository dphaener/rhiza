# done-own-command-guide-verified-2

evidence · `01M785J1T45V66FBPGBNTF6GBS` · actor [claude-code-1ae78044](../actor/claude-code-1ae78044.md)

## result

PASS

## statement

At d46c2728c: go test -trimpath -count=1 ./internal/guide -run 'Build|Review' -v passed (TestBuildText and TestReviewText require 'rhiza session done', 'as its own command' and '`;` or `&&`' in build.md and review.md; both texts hold the rule); go test -trimpath -count=1 ./internal/guide and ./internal/workflow passed; TestBuiltinPromptsDoneOwnCommand scans every builtin prompt naming rhiza session done (all prompts are single-line, so all are scanned) for '; ' / '&& ' chaining and requires 'very last action', and passes; grep finds no chained session done in internal/workflow/builtin/*.toml. Read the fix diff: review.md scopes the throwaway checkout back to 'Checks and critics ... at the branch head', resolving the major quality/model findings.

## Links

### Out

- verifies [done-own-command-guide](../capability/rhiza/guide/done-own-command-guide.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
