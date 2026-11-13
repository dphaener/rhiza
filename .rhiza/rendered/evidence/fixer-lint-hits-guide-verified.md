# fixer-lint-hits-guide-verified

evidence · `01M7B8K0698874FFBNJBGEA424` · actor [claude-code-0554662e](../actor/claude-code-0554662e.md)

## result

PASS

## statement

At c86be0c27 (on 453848a9d): go test -count=1 -trimpath ./internal/guide/... exit 0; rhiza ship generate exit 0; rhiza render exit 0, tree clean afterwards. The diff touches only internal/guide/text/fixer.md, internal/guide/fixer_test.go and the actor's .rhiza records. fixer.md has '### A lint-fix fixer' inside Must do, naming `rhiza run lint-check --json`, the three fixes (rhiza set for stale statement/acceptance, retire a capability whose code was removed, decided-by link), unfixed-with-why, never reproduce a hit as a failing test, and ends with one Review:-trailed commit then rhiza session done. fixer_test.go checks these phrases, the existing rules (fix only named findings, Checks: block, never push) and the five rules restored per decision fixer-lint-section-tighten-keep-rules. The rendered topic is about 3987 chars (~997 of TopicCap 1000) and the cap test passes. No stubs or TODOs. Minor wording left by the trim: Read first starts with a dangling 'Then', and the '(leaf fixer)' hint on the Leaf: trailer was dropped.

## Links

### Out

- verifies [fixer-lint-hits-guide](../capability/fixer-lint-hits-guide.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
