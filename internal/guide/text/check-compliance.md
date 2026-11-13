# Compliance check

## Task

Check one built change set against its contract: `RHIZA_CHANGE`, `RHIZA_CYCLE`,
`RHIZA_CHECK` and `RHIZA_ROLE` are set. You run in a throwaway checkout of the branch.
Judge two things beside each other: each leaf's `acceptance` met as written, and the
scope, meaning nothing missing and nothing extra.

You know nothing of how the work was built. Never run `rhiza run status` or
rhiza run show (or tail) on the build run, and never read its findings, its run notes or its
session transcripts. You judge only the contract from the slice, the diff from the
merge-base and the model's Evidence.

## Expected outcome

One round of findings recorded with `rhiza run findings --check`, each finding one line
with its evidence, or `{"findings":[]}` when every leaf is met and the scope is clean (a
check with no round fails). You change no code.

## Read first

```
rhiza slice --change "$RHIZA_CHANGE" --json
git diff $(git merge-base origin/main HEAD)
```

The slice gives the intent and the leaves, each with its `acceptance`; the diff is
everything the branch changed. The model's Evidence on a leaf is a claim to check, not a
verdict.

## Must do

Work leaf by leaf, and within a leaf point by point through its `acceptance`: find the code
or test in the diff that meets each point exactly as written, and run the test the
acceptance names when it names one. Then read the diff once more against the intent for
what no leaf asked for. Grade as follows.

- Under-built: a leaf with nothing in the diff for it is a `blocker`; an acceptance point
  that is unmet or only partly met is a `major`.
- Scope creep: behaviour outside the intent and the leaves is a `major`.
- Unaccounted files: a changed file that no leaf and not the intent accounts for is a
  `major` if it changes behaviour and a `minor` otherwise.

Write each finding as one line with its evidence, so a fixer can act without
re-diagnosing: the leaf, the place (`path:line`), the point it breaks, what is missing or extra. Record one round, `{"findings":[]}` when
clean:

```
rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'
{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":"<acceptance point unmet and what the code does instead>"}]}
EOF
```

Schema: `rhiza describe run.findings`.

Test with fakes only, never real Rex, agents or `gh`.

<!-- shared: test-scope -->

<!-- shared: command-budget nested -->

## Must not do

<!-- shared: never-list nested -->
- Judge style, naming or test quality: the quality check does that.
- Edit code: a fixer acts on your findings.

## Context

<!-- shared: session-done -->

<!-- shared: severities nested -->

A critic then confirms, downgrades or drops your blockers and majors; fixers act on
what stays.
