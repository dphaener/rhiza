# Plan check

## Task

Check one built change set against its contract: `RHIZA_CHANGE`, `RHIZA_CYCLE`,
`RHIZA_CHECK` and `RHIZA_ROLE` are set. You run in a throwaway checkout of the branch.
Judge each leaf's `acceptance`, the scope and the
code against the slice's constraints and Decisions.

You know nothing of how the work was built. Never run `rhiza run status` or
rhiza run show (or tail) on the build run, and never read its findings, its run notes or its
session transcripts. You judge only the contract from the slice, the diff from the
merge-base and the model's Evidence.

## Expected outcome

One round recorded with `rhiza run findings --check`: findings, each one line with its
evidence (`"findings":[]` when clean), and evidence lines for every leaf (a check with no
round fails).

## Read first

```
rhiza slice --change "$RHIZA_CHANGE" --json
git diff $(git merge-base origin/main HEAD)
```

## Must do

Work leaf by leaf through each `acceptance`: find the code or test that meets each
point as written, and run the test it names. Then read the diff for what no leaf asked
for, and against each constraint and Decision. Grade as follows.

- Under-built: a leaf with nothing in the diff for it is a `blocker`; an acceptance point
  that is unmet or only partly met is a `major`.
- Scope creep: behaviour outside the intent and the leaves is a `major`.
- Unaccounted files: a changed file no leaf or the intent accounts for is a
  `major` if it changes behaviour and a `minor` otherwise.
- Code contradicts a node: code that does what a leaf, Decision or constraint forbids, or
  the reverse, is a `major`.
- Missing constraint: a rule the change introduced that no constraint states is a `minor`.

Write each finding as one line with its evidence (the leaf, `path:line`, the point or
node it breaks), so a fixer needs no re-diagnosing. Write one evidence line per
acceptance point of each leaf, findings or not: what you ran or read, with the result.
Record one round:

```
rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'
{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":"<what is unmet or contradicted, and what the code does>"}],
"evidence":[{"leaf":"<leaf>","line":"<what you ran or read for one point, with the result>"}]}
EOF
```

In a conductor run, a round that leaves a leaf uncovered (no finding naming it, no
evidence line) is refused and not recorded.

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
