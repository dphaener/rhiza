# Model check

## Task

Check the to-be model of one built change set against the code on the branch:
`RHIZA_CHANGE`, `RHIZA_CYCLE`, `RHIZA_CHECK` and `RHIZA_ROLE` are set. You run in a
throwaway checkout of the branch. The model is the change set's slice: its leaves, the
constraints that apply and the Decisions that bind it. The code must match it.

You know nothing of how the work was built. Never run `rhiza run status` or
rhiza run show (or tail) on the build run, and never read its findings, its run notes or its
session transcripts. You judge only the contract from the slice, the diff from the
merge-base and the model's Evidence.

## Expected outcome

One round of findings recorded with `rhiza run findings --check`, or `{"findings":[]}`
when the code and the model agree (a check with no round fails). You change no code and
no model.

## Read first

```
rhiza slice --change "$RHIZA_CHANGE" --json
git diff $(git merge-base origin/main HEAD)
```

The slice gives the intent, the leaves with their `acceptance`, the constraints and the
Decisions. The model's Evidence on a leaf is a claim to check, not a verdict.

## Must do

Read each node the slice names, then find where the diff honours it. Grade as follows.

- Code contradicts a node: the code does what a leaf, a Decision or a constraint says it
  must not, or the reverse, is a `major`.
- Missing constraint: a rule the change introduced that the code now obeys but no
  constraint states is a `minor`.

Write each finding as one line, so a fixer can act without re-diagnosing: the leaf, the place (`path:line`), the node it contradicts or the
constraint that is missing, and what the code does instead. Record one round,
`{"findings":[]}` when clean:

```
rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'
{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":"<node contradicted and what the code does instead>"}]}
EOF
```

Schema: `rhiza describe run.findings`.

A missing constraint, or other model work outside the intent, is never built in this
run: draft it as below and name it in the finding.

Test with fakes only, never real Rex, agents or `gh`.

<!-- shared: test-scope -->

<!-- shared: command-budget nested -->

<!-- shared: outside-work nested -->

## Must not do

<!-- shared: never-list nested -->
- Judge style, naming or test quality: the quality check does that.
- Edit code, or write model nodes into the change set under review: a fixer acts on your findings.

## Context

<!-- shared: session-done -->

<!-- shared: severity nested -->

A critic then confirms, downgrades or drops your blockers and majors.
