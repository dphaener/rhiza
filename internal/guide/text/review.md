# Review

A review run checks a built change set. Each session gets `RHIZA_CHANGE`, `RHIZA_CYCLE`
and `RHIZA_ROLE` (check, critic, fixer or verifier); checks and critics also `RHIZA_CHECK`,
a leaf's fixer and verifier `RHIZA_LEAF`. Start with
`rhiza slice --change "$RHIZA_CHANGE" --json`. Checks and critics work in a
throwaway checkout at the branch head; their rhiza writes reach the run's log.
Review the diff from the merge-base with the default branch
(`git diff $(git merge-base origin/main HEAD)`). Each role ends with `rhiza session done`.

## Checks

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

Record one round, `{"findings":[]}` when clean (a check with no round fails), with
`rhiza run findings --check "$RHIZA_CHECK" <file>`:

```
{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":"<why>"}]}
```

`leaf` is optional. Severities: `blocker` (wrong or broken; cannot ship), `major` (must be
fixed), `minor` and `nit` (reported, never fixed). Schema: `rhiza describe run.findings`.

## Critic

Read the check's latest round in `rhiza run status --json`. Confirm or drop each blocker
and major against the code. Never add a finding. Record with
`rhiza run verdict --check "$RHIZA_CHECK" <file>`:

```
{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":"<why>"}]}
```

`verdict` is `confirmed` or `dropped`. Dropped findings stay in the run journal; no Decision
is written for them.

## Fixer

Your prompt names the confirmed findings to fix. Fix them, nothing else, in one new
commit with the `Change:` and `Actor:` trailers plus `Leaf: <leaf>` for a leaf fixer, or
`Review: <run>` ($RHIZA_RUN) for the change set's fixer. Never rewrite earlier commits.
A verifier re-run follows `rhiza guide build`.

## Outside work

Work outside the intent, found by scope or any check, is a question, or a
discovered-from draft change set: `rhiza change new "<intent>" --json`, then
`rhiza link <new> discovered-from <this> --change <new> --json`, then
`rhiza change use <this>`.
