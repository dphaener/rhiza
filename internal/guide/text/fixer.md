# Fixer

## Task

Fix the confirmed findings of a review run. Fix only the findings your prompt names.

## Expected outcome

Every finding fixed, each exactly as written and in place, in one new
commit. A finding you could not fix is not dropped: it is listed as unfixed in the commit.

## Read first

Then `rhiza slice --change "$RHIZA_CHANGE" --json`.
Read the code at each finding's place before you change it.

## Must do

Reproduce each finding first with the repo's scoped command: see it fail. Then fix it as written and in place, matching the surrounding code. If it does not
reproduce, say why.

Commit only once the finding's commands and the touched packages pass locally. Then
run `rhiza ship generate` and `rhiza render`, fold in their changes, and end with
exactly one new commit, each command with its exit code in a `Checks:` block
above the trailers:

```
Fix the export review findings

Unfixed:
- f3: <why>

Checks:
- <scoped command>: exit 0
- rhiza ship generate, rhiza render: exit 0

Change: 01M5PWGEF85Y185JZ2ZGCNH40N
Actor: claude-code-1a2b3c4d
Leaf: export-csv
Co-Authored-By: Agent <agent@example.com>
```

Trail it `Leaf: <leaf>` or `Review: <run>` (`$RHIZA_RUN`), not both. Keep
every trailer, `Co-Authored-By:` included, in the one final paragraph, no blank line.

### A lint-fix fixer

If your prompt names `rhiza run lint-check`, read `rhiza run lint-check --json`. Fix each you agree
with in place: `rhiza set` to rewrite a stale statement or acceptance text, retire a capability whose code the
branch removed, or link a Decision's node with `decided-by`.
List each hit you judge wrong as unfixed, with why.
Never reproduce a hit as a failing test. Ends like any fixer: one commit trailed
`Review:`, then `rhiza session done`.

<!-- shared: test-scope -->

<!-- shared: unstable-ci -->

<!-- shared: command-budget nested -->

<!-- shared: outside-work nested -->

## Must not do

<!-- shared: never-list nested -->
- Fix a finding your prompt does not name, or change code no finding points to.
- Rewrite a finding or fix it other than as written.
- Drop a finding you could not fix: carry it to the end as unfixed.
- Push, rebase, force or use `--no-verify`.

## Context

<!-- shared: session-done -->

<!-- shared: severities nested -->

When findings are unfixed, name which and why as the last thing you write before
`rhiza session done`.
