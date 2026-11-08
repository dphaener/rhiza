# Review

A review run checks a built change set. Each session gets `RHIZA_CHANGE`, `RHIZA_CYCLE`
and `RHIZA_ROLE` (check, critic, fixer or verifier); checks and critics also `RHIZA_CHECK`,
a leaf's fixer and verifier `RHIZA_LEAF`. Start with
`rhiza slice --change "$RHIZA_CHANGE" --json`. Checks and critics work in a
throwaway checkout at the branch head.
Review the diff from the merge-base with the default branch
(`git diff $(git merge-base origin/main HEAD)`). Each role ends with `rhiza session done`.

## Checks

Every check except `qa` tests runs with fakes only, never real Rex, agents or `gh`.
Checks, critics and fixers run only the tests covering what the change touches, with the
repo's scoped commands, never the repo's full
suite: the PR's CI runs it before merge.

While CI is unstable, or a red is not explained by your change,
do not push or rerun CI to see. Reproduce the
failing test locally, fix it, and push only after it and the
touched packages pass locally.

- `compliance`: each leaf's `acceptance`, met as written.
- `quality`: correctness, tests, the repo's conventions.
- `qa`: hands-on, with the built binary; fakes alone do not count. Then
  Evidence written as yourself on each leaf it exercised, via
  `rhiza apply <file> --change "$RHIZA_CHANGE" --json`, `<cycle>` being $RHIZA_CYCLE:

  ```json
  {"writes":[
    {"op":"add","kind":"evidence","handle":"<leaf>-qa-<cycle>",
     "fields":["statement=<what you ran and saw>","result=PASS"]},
    {"op":"link","from":"<leaf>-qa-<cycle>","type":"verifies","to":"<leaf>"}
  ]}
  ```

- `scope`: work stays inside the intent and the leaves.
- `model`: the model matches the code; missing constraints go in a discovered-from
  draft change set.

Record one round, `{"findings":[]}` when clean (a check with none fails), on stdin:

```
rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'
{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":"<why>"}]}
EOF
```

`leaf` is optional. Severities: `blocker` (broken; cannot ship), `major` (must be
fixed), `minor` and `nit` (reported, never fixed). Schema: `rhiza describe run.findings`.

## Critic

Read the check's latest round in `rhiza run status --json`. Confirm or drop each blocker
and major against the code. Never add a finding. Record on stdin:

```
rhiza run verdict --check "$RHIZA_CHECK" - <<'EOF'
{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":"<why>"}]}
EOF
```

`verdict` is `confirmed` or `dropped`. Dropped findings stay in the run journal; no Decision.

## Fixer

Your prompt names the confirmed findings to fix. Fix only those, in one new
commit with the `Change:` and `Actor:` trailers plus `Leaf: <leaf>` for a leaf fixer, or
`Review: <run>` ($RHIZA_RUN) for the change set's fixer. Keep every trailer, `Co-Authored-By:`
included, in the commit's one final paragraph, with no blank line between them.
Never rewrite earlier commits.
A verifier re-run follows `rhiza guide build`.

## Outside work

Work outside the intent is a question, or a discovered-from draft change set:
`rhiza change new "<intent>" --json`, then
`rhiza link <new> discovered-from <this> --change <new> --json`, then
`rhiza change use <this>`.
