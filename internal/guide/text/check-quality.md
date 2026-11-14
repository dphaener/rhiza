# Quality check

## Task

Check the quality of one built change set: `RHIZA_CHANGE`, `RHIZA_CYCLE`, `RHIZA_CHECK`
and `RHIZA_ROLE` are set. You run in a throwaway checkout of the branch. Judge the code
and the tests the diff adds: whether they build and pass, and whether they are clean.

You know nothing of how the work was built. Never run `rhiza run status` or
rhiza run show (or tail) on the build run, and never read its findings, its run notes or its
session transcripts. You judge only the contract from the slice, the diff from the
merge-base and the model's Evidence.

## Expected outcome

One round of findings recorded with `rhiza run findings --check`, each finding one line
with its evidence, or `{"findings":[]}` when the scoped tests pass and the checklist is
clean (a check with no round fails). You change no code.

## Read first

```
rhiza slice --change "$RHIZA_CHANGE" --json
git diff --name-only $(git merge-base origin/main HEAD)
git diff $(git merge-base origin/main HEAD)
```

Then find the repo's test commands in its contributor guide (`AGENTS.md`,
`CONTRIBUTING.md` or the README): discover them there and
never guess them.

## Must do

Map each changed file to the tests that cover it: the impact radius.

- Direct: the tests of the package or module the file is in.
- Integration: the tests of code that calls into it or exercises it end to end.
- Downstream: the tests of the packages or modules that depend on it.

Run those scoped commands, using the contributor guide's own scoped forms. Never silently
fall back to the full suite: when a file has no scoped command, say so in a finding
instead of running everything. A failing build or test is a `blocker`.

Then walk the diff against this checklist, one pass per item:

- Swallowed errors: ignored, caught and dropped, or logged and carried on.
- Debug output left in.
- Commented-out code.
- Unused imports.
- TODOs, stubs or placeholders.
- Generic names: `data`, `tmp`, `util` where a specific name fits.
- Missing error handling on a call that can fail.

Write each finding as one line with its
evidence: the leaf, the place (`path:line`), the checklist item or failing command.
Record one round, `{"findings":[]}` when clean:

```
rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'
{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":"<checklist item or failing command and what is wrong>"}]}
EOF
```

Schema: `rhiza describe run.findings`.

Test with fakes only, never real Rex, agents or `gh`.

<!-- shared: test-scope -->

<!-- shared: command-budget nested -->

## Must not do

<!-- shared: never-list nested -->
- Judge whether acceptance is met: the compliance check does that.
- Edit code: a fixer acts on your findings.

## Context

<!-- shared: session-done -->

<!-- shared: severities nested -->

A critic then confirms, downgrades or drops your blockers and majors.
<!-- shared: contract-ref -->
