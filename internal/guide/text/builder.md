# Builder

## Task

Build one leaf of a change set: `RHIZA_CHANGE`, `RHIZA_LEAF` and `RHIZA_ATTEMPT` are set.
Build only that leaf.

## Expected outcome

The leaf's `acceptance` is the contract: when it holds, you are done. One new commit
carries the code, the tests and your staged actor record, and its `Checks:` block lists
each command you ran with its exit code.

## Read first

```
rhiza slice --change "$RHIZA_CHANGE" --json
rhiza show "$RHIZA_LEAF" --json
```

Read the slice, the leaf and the existing code and tests that follow the same pattern
before you write anything. Match that pattern.

If the owner must answer something, raise a question before you start work, not midway.

## Must do

Write code and tests to the leaf's acceptance. Then run its acceptance
commands, `rhiza ship generate` and `rhiza render`; commit what they change.
End with exactly one new commit, each command's exit code in a
`Checks:` block above the trailers:

```
Stream CSV exports

Checks:
- <acceptance command>: exit 0
- rhiza ship generate: exit 0
- rhiza render: exit 0

Change: 01M5PWGEF85Y185JZ2ZGCNH40N
Actor: claude-code-1a2b3c4d
Leaf: export-csv
Co-Authored-By: Agent Name <agent@example.com>
```

Keep every trailer, `Co-Authored-By:` included, in the one final paragraph, no blank line.

Stage your untracked actor record in `.rhiza/log`, never another actor's:

```
git add $(grep -l "Record agent actor $RHIZA_ACTOR " .rhiza/log/*/*.json | xargs -n1 dirname | sort -u)
```

On a retry (`RHIZA_ATTEMPT` > 1), load `/rhiza-build-retry` (`rhiza guide build-retry`) first.

<!-- shared: test-scope -->

<!-- shared: unstable-ci -->

<!-- shared: command-budget nested -->

<!-- shared: outside-work nested -->

## Must not do

<!-- shared: never-list nested -->

## Context

<!-- shared: session-done -->

A verifier then re-runs your `Checks:` and the acceptance; findings come back as a retry.
