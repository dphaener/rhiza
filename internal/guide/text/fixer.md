# Fixer

## Task

Fix the confirmed findings of a review run: `RHIZA_CHANGE`, `RHIZA_ACTOR` and
`RHIZA_RUN` are set, a leaf fixer also `RHIZA_LEAF`. Fix only the findings your prompt names.

## Expected outcome

Every finding in your prompt reproduced and fixed, each exactly as written and in place,
in one new commit whose `Checks:` block lists each command you ran with its exit code.
A finding you could not fix is not dropped: it is listed as unfixed in the commit.

## Read first

Your prompt: the confirmed findings, each with severity, place and reason. Then:

```
rhiza slice --change "$RHIZA_CHANGE" --json
```

Read the code at each finding's place before you change anything.

## Must do

Reproduce each finding first, with the repo's own scoped command for it: see it fail
before you fix it. Then fix it exactly as written and in place, matching the
surrounding code. A finding that does not reproduce: say so and why in the commit.

Commit only once the commands for your findings and the packages you touched pass
locally. Then run `rhiza ship generate` and `rhiza render`, fold what they change into
the commit, and end with exactly one new commit, each command's exit code in a
`Checks:` block above the trailers:

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

When you end with findings unfixed, say which and why as the last thing you write
before `rhiza session done`.
