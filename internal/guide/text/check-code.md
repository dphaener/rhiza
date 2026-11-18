# Code check

## Task

Check the code of one built change set: `RHIZA_CHANGE`, `RHIZA_CYCLE`, `RHIZA_CHECK`
and `RHIZA_ROLE` are set. You run in a throwaway checkout of the branch. Judge whether
the code and tests the diff adds build, pass and are clean.

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

Then find the repo's test commands in its contributor guide (`AGENTS.md`,
`CONTRIBUTING.md` or the README): discover them there and
never guess them.

## Must do

Map each changed file to the tests that cover it: the impact radius.

- Direct: the tests of the file's package or module.
- Integration: the tests of code that calls it or runs it end to end.
- Downstream: the tests of what depends on it.

Run them with the contributor guide's scoped forms. Never silently
fall back to the full suite: when a file has no scoped command, say so in a finding.
A failing build or test is a `blocker`.

Then walk the diff, one pass per checklist item:

- Swallowed errors: ignored, dropped, or logged and carried on.
- Debug output left in.
- Commented-out code.
- Unused imports.
- TODOs, stubs or placeholders.
- Generic names: `data`, `tmp`, `util`.
- Missing error handling on a call that can fail.

Write each finding as one line with its evidence: the leaf, `path:line`, the item or
failing command. Write one evidence line per acceptance point of each leaf, findings or
not: what you ran or read, with the result. Record one round:

```
rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'
{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":"<checklist item or failing command, what is wrong>"}],
"evidence":[{"leaf":"<leaf>","line":"<what you ran or read for one point, with the result>"}]}
EOF
```

In a conductor run, a round that leaves a leaf uncovered (no finding naming it, no
evidence line) is refused and not recorded. Schema: `rhiza describe run.findings`.

Test with fakes only, never real Rex, agents or `gh`.

<!-- shared: test-scope -->

<!-- shared: command-budget nested -->

## Must not do

<!-- shared: never-list nested -->
- Judge whether acceptance is met: the plan check does that.
- Edit code: a fixer acts on your findings.

## Context

<!-- shared: session-done -->

<!-- shared: severities nested -->

A critic then confirms, downgrades or drops your blockers and majors.
<!-- shared: contract-ref -->
