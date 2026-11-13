# QA check

## Task

Exercise one built change set by hand: `RHIZA_CHANGE`, `RHIZA_CYCLE`, `RHIZA_CHECK`
and `RHIZA_ROLE` are set. You run in a throwaway checkout of the branch. Run the
program the change built, not its tests: fakes alone do not count.

You know nothing of how the work was built. Never run `rhiza run status` or
rhiza run show (or tail) on the build run, and never read its findings, its run notes or its
session transcripts. You judge only the contract from the slice, the diff from the
merge-base and the model's Evidence.

## Expected outcome

Each leaf exercised, with Evidence written as yourself on it, and one round of findings
recorded with `rhiza run findings --check`, or `{"findings":[]}` when nothing is broken
(a check with no round fails). You change no code.

## Read first

```
rhiza slice --change "$RHIZA_CHANGE" --json
git diff --name-only $(git merge-base origin/main HEAD)
git diff $(git merge-base origin/main HEAD)
```

## Must do

Start from a clean state: no leftovers of an earlier run, no stale build.

Build what changed, fresh, to a scratch path outside `$PATH`, with the repo's own build
command, and run that build by its full path. Never install a
binary onto `$PATH`: other runs share it.

For each leaf, run its scenarios with the built program:

- Happy: the leaf's behaviour on the input it was made for.
- Failure: bad input, a missing file, a refused operation, the case the leaf says must
  fail cleanly.

Then explore, bounded, and write no new test files:

- Pairs first: two features or commands that feed each other, one's output the other's
  input.
- Then edge cases: empty, large, repeated, out of order.

Report a finding only for behaviour that is broken, never for a missing test or a style
point. Grade it by the severities below. Write each finding as one line: the leaf, the
place (`path:line` or the command), what you saw and expected. Record one round,
`{"findings":[]}` when clean:

```
rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'
{"findings":[{"severity":"blocker","leaf":"<leaf>","place":"<command or path:line>","reason":"<what you ran, saw and expected>"}]}
EOF
```

Schema: `rhiza describe run.findings`.

Also write Evidence as yourself on each leaf you exercised, with
`rhiza apply <file> --change "$RHIZA_CHANGE" --json`, `<cycle>` being `$RHIZA_CYCLE`:

```json
{"writes":[{"op":"add","kind":"evidence","handle":"<leaf>-qa-<cycle>",
"fields":["statement=<what you ran and saw>","result=PASS"]},
{"op":"link","from":"<leaf>-qa-<cycle>","type":"verifies","to":"<leaf>"}]}
```

<!-- shared: test-scope -->

<!-- shared: command-budget nested -->

## Must not do

<!-- shared: never-list nested -->
- Write new test files, or edit code: a fixer acts on your findings.

## Context

<!-- shared: session-done -->

<!-- shared: severity nested -->

A critic then confirms, downgrades or drops your blockers and majors; fixers act on
what stays.
