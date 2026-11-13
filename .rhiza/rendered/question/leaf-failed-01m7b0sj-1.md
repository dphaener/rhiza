# leaf-failed-01m7b0sj-1

question · `01M7B86XM4758S3QE6K0MN1XT9`

## statement

Leaf fixer-lint-hits-guide failed all 3 attempts of round 1: attempt 1 verifier failed: verifier found 2 problem(s) (f1 internal/guide/text/fixer.md: the leaf statement says "The rest of the fixer topic is unchanged", but commit 453848a9d rewrote and dropped existing fixer rules…; f2 internal/guide/text/fixer.md "### A lint-fix fixer": the statement requires the section to say the lint-fix fixer ends with one commit trailed `Review:` as any…); attempt 2 builder failed: no new commit since 453848a9d655; attempt 3 builder failed: no new commit since 453848a9d655. How should it be built, or should it change?

## Links

### Out

- blocks [fixer-lint-hits-guide](../capability/fixer-lint-hits-guide.md)

### In

- answers [fixer-lint-section-tighten-keep-rules](../decision/fixer-lint-section-tighten-keep-rules.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
