# Review

Sessions get `RHIZA_CHANGE`, `RHIZA_CYCLE`
and `RHIZA_ROLE` (check, critic, fixer or verifier); checks and critics also `RHIZA_CHECK`,
a leaf's fixer and verifier `RHIZA_LEAF`. Start with
`rhiza slice --change "$RHIZA_CHANGE" --json`. Checks and critics: a throwaway checkout.
Review from the merge-base with the default branch:
`git diff $(git merge-base origin/main HEAD)`. Each role ends with `rhiza session done` as its own command once its record
is written, never chained by `;` or `&&`.

## Checks

Every check except `qa` tests runs with fakes only, never real Rex, agents or `gh`.
Checks, critics and fixers run only the tests covering what the change touches, with
scoped commands, never the repo's full
suite: the PR's CI runs it before merge.

While CI is unstable, or a red is not explained by your change,
do not push or rerun CI to see. Reproduce the
failing test locally, fix it, and push only after it and the
touched packages pass locally.
Reproduce a flake by injecting a delay at the suspected race, never by busy loops or parallel stress on the host.

- `compliance`: each leaf's `acceptance`, met as written.
- `quality`: correctness, tests, conventions.
- `qa`: hands-on with the built binary; fakes alone do not count. Then
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
- `model`: it matches the code; missing constraints: a discovered-from
  draft change set.

Record one round on stdin, `{"findings":[]}` when clean (none fails):

```
rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'
{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":"<why>"}]}
EOF
```

Severities: `blocker`, `major` (must be fixed), `minor`, `nit` (never fixed). Schema: `rhiza describe run.findings`.

## Critic

Read the check's last round (`rhiza run status --json`). Confirm or drop each blocker and major. Never add a finding. Record on stdin:

```
rhiza run verdict --check "$RHIZA_CHECK" - <<'EOF'
{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":"<why>"}]}
EOF
```

`confirmed` or `dropped`; dropped ones stay in the run journal, no Decision.

## Fixer

Fix only the confirmed findings named, in one new commit trailed `Change:`, `Actor:` and `Leaf: <leaf>` (leaf fixer) or
`Review: <run>` ($RHIZA_RUN, change set fixer). Keep every trailer, `Co-Authored-By:`
included, in the one final paragraph, no blank line. Never rewrite commits.

## Command budget

Never sleep and never poll. No command runs longer than about four
minutes; a longer build or test runs once in the background and is
waited for once.

## Outside work

Work outside the intent is a question or a discovered-from draft change set:
`rhiza change new "<intent>" --json`, then
`rhiza link <new> discovered-from <this> --change <new> --json`, then
`rhiza change use <this>`.
