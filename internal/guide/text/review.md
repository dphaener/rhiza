# Review

Sessions get `RHIZA_CHANGE`, `RHIZA_CYCLE` and `RHIZA_ROLE`; checks and critics also
`RHIZA_CHECK`, leaf fixers `RHIZA_LEAF`. Start with
`rhiza slice --change "$RHIZA_CHANGE" --json`. Checks and critics: a throwaway checkout.
Review from the merge-base with the default branch:
`git diff $(git merge-base origin/main HEAD)`. Each role ends with `rhiza session done` as its own command after its recording
command succeeded, never chained with `;` or `&&`.

## Checks

Every check except `qa` tests runs with fakes only, never real Rex, agents or `gh`.
Checks, critics and fixers run only the tests covering what the change touches, by
scoped commands, never the repo's full
suite: the PR's CI runs it before merge.

While CI is unstable or a red is not explained by your change,
do not push or rerun CI to see. Reproduce the
failing test locally, fix it, push only after it and the
touched packages pass locally.
Reproduce flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host.

- `compliance`: each leaf's `acceptance` as written.
- `quality`: correctness, tests, conventions.
- `qa`: hands-on, built binary; fakes alone do not count. Then
  Evidence written as yourself on each leaf it exercised, via
  `rhiza apply <file> --change "$RHIZA_CHANGE" --json`, `<cycle>` being $RHIZA_CYCLE:

  ```json
  {"writes":[{"op":"add","kind":"evidence","handle":"<leaf>-qa-<cycle>",
  "fields":["statement=<what you ran and saw>","result=PASS"]},
  {"op":"link","from":"<leaf>-qa-<cycle>","type":"verifies","to":"<leaf>"}]}
  ```

- `scope`: work stays inside the intent and the leaves.
- `model`: model matches code; missing constraints: a discovered-from draft.

Record one round, `{"findings":[]}` when clean (a check with none fails):

```
rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'
{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":"<why>"}]}
EOF
```

Schema: `rhiza describe run.findings`.

## Severities

- `blocker`: behaviour broken, a build or test fails, or a leaf missing.
- `major`: an acceptance point unmet, or a real defect.
- `minor`: style, naming, docs or a non-blocking follow-up.
- `nit`: taste.

Per check: `compliance` leaf missing blocker, point unmet major; `quality` failing build or
test blocker, defect major, convention minor; `qa` broken behaviour blocker, rough edge
minor; `scope` outside the intent major; `model` code contradicted major, missing
constraint minor. Plan review: `missing-acceptance`, `contradiction`, `outside-intent`
major, the rest minor. `minor` and `nit` are never fixed.

## Critic

Confirm, downgrade or drop each blocker and major of the check's last round (`rhiza run status --json`). Never add a finding or raise a severity.

```
rhiza run verdict --check "$RHIZA_CHECK" - <<'EOF'
{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":"<why>"}]}
EOF
```

`confirmed` keeps the severity, `downgraded` lowers it (`"severity":"minor"` or `"nit"`), `dropped` removes the finding; all stay in the run journal, no Decision.

## Fixer

Fix only the named confirmed findings in one new commit trailed `Change:`, `Actor:` and `Leaf: <leaf>` (leaf fixer) or
`Review: <run>` ($RHIZA_RUN). Keep all trailers, `Co-Authored-By:`
included, in one final paragraph, no blank line. Never rewrite commits.

First reproduce your findings and run `rhiza ship generate` and `rhiza render`; commit with a `Checks:` block of exit codes.

## Command budget

Never sleep and never poll. No command runs over about four
minutes; longer work runs once in the background and is
waited for once.

## Outside work

Work outside the intent: a question or a discovered-from draft change set:
`rhiza change new "<intent>" --json`,
`rhiza link <new> discovered-from <this> --change <new> --json`, then
`rhiza change use <this>`.
