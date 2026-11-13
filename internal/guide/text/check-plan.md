# Plan check

## Task

Check one built change set against its contract: `RHIZA_CHANGE`, `RHIZA_CYCLE`,
`RHIZA_CHECK` and `RHIZA_ROLE` are set. You run in a throwaway checkout of the branch.
Judge beside each other: each leaf's `acceptance` met as written, the scope (nothing
missing, nothing extra) and the code against the slice's constraints and Decisions.

You know nothing of how the work was built. Never run `rhiza run status` or
rhiza run show (or tail) on the build run, and never read its findings, its run notes or its
session transcripts. You judge only the contract from the slice, the diff from the
merge-base and the model's Evidence.

## Expected outcome

One round of findings recorded with `rhiza run findings --check`, each finding one line
with its evidence, or `{"findings":[]}` when all is met and clean (a check with no round
fails). You change no code and no model.

## Read first

```
rhiza slice --change "$RHIZA_CHANGE" --json
git diff $(git merge-base origin/main HEAD)
```

The slice gives the intent, the leaves with their `acceptance`, the constraints and the
Decisions; the diff is everything the branch changed.

## Must do

Work leaf by leaf, point by point through its `acceptance`: find the code or test in the
diff that meets each point as written, and run the test it names. Then read the diff against the intent for what no leaf asked
for, and against each constraint and Decision. Grade as follows.

- Under-built: a leaf with nothing in the diff for it is a `blocker`; an acceptance point
  that is unmet or only partly met is a `major`.
- Scope creep: behaviour outside the intent and the leaves is a `major`.
- Unaccounted files: a changed file that no leaf and not the intent accounts for is a
  `major` if it changes behaviour and a `minor` otherwise.
- Code contradicts a node: code that does what a leaf, Decision or constraint forbids, or
  the reverse, is a `major`.
- Missing constraint: a rule the change introduced that no constraint states is a `minor`.

Write each finding as one line with its evidence, so a fixer needs no re-diagnosing: the leaf, the place (`path:line`), the point or node it breaks, what is wrong. Record one round:

```
rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'
{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":"<acceptance point unmet or node contradicted, and what the code does instead>"}]}
EOF
```

Test with fakes only, never real Rex, agents or `gh`.

<!-- shared: test-scope -->

<!-- shared: command-budget nested -->

## Must not do

<!-- shared: never-list nested -->
- Judge style, naming or test quality: the code check does that.
- Edit code, or write model nodes into the change set under review: a fixer acts on your findings.

## Context

<!-- shared: session-done -->

<!-- shared: severities nested -->

A critic then confirms, downgrades or drops your blockers and majors.

<!-- shared: contract-ref -->
