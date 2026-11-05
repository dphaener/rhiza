# question-b947d57e-1

question · `01M6JWYSPB5EGQSMT7B69C5NW8`

## session

b947d57e-7bee-4636-b577-01a1a5cce3e1

## statement

PR #164 already switched rhiza run abandon to `git branch -d`, and run-abandon-unmerged.txtar shows it keeping a never-pushed branch with no upstream. The model has no leaf for this, and nothing tests a branch that has an upstream plus local commits never pushed to it. What should change set 01M6JWSC build?

## Links

### In

- answers [abandon-unpushed-branch-scope](../decision/abandon-unpushed-branch-scope.md)

## History

- `01M6JWSCZYYQVFR6ZTVWS6JWJ0` Fix rhiza run abandon keeps a local branch holding commits that were never pushed, instead of force-deleting them (replaces draft change set 01M6B9FCNTTH65, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
