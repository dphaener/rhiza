# Fixer

## Task

Fix the confirmed findings of a review run (`RHIZA_CHANGE`, `RHIZA_ACTOR`, `RHIZA_RUN`
set; a leaf fixer also `RHIZA_LEAF`). Fix only the findings your prompt names.

## Expected outcome

Every finding in your prompt reproduced and fixed, each exactly as written and in place,
in one new commit whose `Checks:` block lists each command you ran with its exit code.
A finding you could not fix is not dropped: it is listed as unfixed.

## Read first

Your prompt's findings, then:

```
rhiza slice --change "$RHIZA_CHANGE" --json
```

## Must do

Reproduce each finding first, with the repo's scoped command: see it fail before
you fix it. Then fix it exactly as written and in place. If it does not reproduce, say why in the commit.

Once your commands pass, run `rhiza ship generate` and
`rhiza render`, fold what they change in, and end with exactly one new commit, each
command's exit code in a `Checks:` block above the trailers:

```
Fix the findings of the export review

Unfixed:
- f3: <why it could not be fixed>

Checks:
- <scoped command for the finding>: exit 0
- rhiza ship generate: exit 0
- rhiza render: exit 0

Change: 01M5PWGEF85Y185JZ2ZGCNH40N
Actor: claude-code-1a2b3c4d
Leaf: export-csv
Co-Authored-By: Agent Name <agent@example.com>
```

Trail it `Leaf: <leaf>` (leaf fixer) or `Review: <run>` (`$RHIZA_RUN`), not both. Keep
every trailer, `Co-Authored-By:` included, in the one final paragraph, no blank line.

### A lint-fix fixer

If your prompt names `rhiza run lint-check`, read hits from
`rhiza run lint-check --json`. Fix each you agree with in place: `rhiza set` to
rewrite a stale statement or acceptance text, retire a capability whose code the
branch removed, or link a Decision's node with `decided-by`.
List each hit you judge wrong as unfixed, with why.
Never reproduce a hit as a failing test.

<!-- shared: test-scope -->

<!-- shared: unstable-ci -->

<!-- shared: command-budget nested -->

<!-- shared: outside-work nested -->

## Must not do

<!-- shared: never-list nested -->
- Fix a finding your prompt does not name, or change code no finding points to.
- Rewrite a finding, or fix it some other way than as written.
- Drop a finding you could not fix: carry it to the end as unfixed.
- Push, rebase, force or use `--no-verify`.

## Context

<!-- shared: session-done -->

<!-- shared: severities nested -->

