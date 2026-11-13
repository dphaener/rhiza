# fixer-lint-section-over-topic-cap

question · `01M7B83W52G6VQ2TQF6AQ3MQ29`

## statement

Leaf fixer-lint-hits-guide: verifier finding f1 requires restoring the fixer topic verbatim ("the rest of the fixer topic is unchanged"), but the restored topic already costs 994 of TopicCap 1000 tokens (guide.Tokens = chars/4), leaving 6 tokens: no lint-fix section (needs ~120 tokens with the Review:/session-done sentence f2 asks for) can fit without trimming other fixer rules, which f1 forbids. Which should the builder do? (a) raise TopicCap (e.g. to 1200) for all topics; (b) put the lint-fix instructions in a new topic (e.g. `rhiza guide lint-fix`) and have fixer only point to it, accepting that the pointer itself still needs a small trim; (c) allow tightening the existing fixer wording without dropping any rule (attempt 1 dropped rules); (d) something else.

## Links

### Out

- blocks [fixer-lint-hits-guide](../capability/fixer-lint-hits-guide.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
