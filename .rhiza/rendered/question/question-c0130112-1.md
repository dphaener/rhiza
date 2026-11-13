# question-c0130112-1

question · `01M7AZK2H8X4PKZZ580SMD571G`

## session

c0130112-6d9f-43e5-ab34-01a1d5f666af

## statement

Model lints in review (change set 01M76XTF): today a failed `rhiza run code-check` ends the review run failed, because the code step has no on.fail. Where should a lint hit go so it saves a review cycle instead of ending the run? / Model lints in review: which capabilities should the stale-reference lint check? It flags repo paths, Go test names and testscript names in a statement or acceptance that do not exist at the branch head. / Model lints in review: how should the lint decide that a Decision was "never applied"?

## Links

### In

- answers [lint-hits-feed-model-check](../decision/lint-hits-feed-model-check.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
