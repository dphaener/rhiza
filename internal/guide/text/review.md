# Review

Review sessions get `RHIZA_CHANGE`, `RHIZA_CYCLE`
and `RHIZA_ROLE` (check, critic, fixer or verifier); checks and critics also `RHIZA_CHECK`,
a leaf's fixer and verifier `RHIZA_LEAF`. Start with
`rhiza slice --change "$RHIZA_CHANGE" --json`. Checks and critics use a throwaway checkout at the head.
Review the diff from the merge-base with the default branch. Each role ends with `rhiza session done` as its own command, never chained with `;` or `&&`.

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

Record a round, `{"findings":[]}` when clean, on stdin:

```
rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'
{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":"<why>"}]}
EOF
```

Severities: `blocker`, `major`, `minor`, `nit`. Schema: `rhiza describe run.findings`.

## Critic

Read the check's latest round (`rhiza run status --json`). Confirm or drop each blocker and major. Never add a finding. Record on stdin:

```
rhiza run verdict --check "$RHIZA_CHECK" - <<'EOF'
{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":"<why>"}]}
EOF
```

`verdict` is `confirmed` or `dropped`. Dropped findings stay in the run journal, no Decision.

## Fixer

Fix only the confirmed findings named, in one new commit with `Change:` and `Actor:` trailers plus `Leaf: <leaf>` for a leaf fixer, or
`Review: <run>` ($RHIZA_RUN) for the change set's fixer. Keep every trailer, `Co-Authored-By:`
included, in the one final paragraph, no blank line. Never rewrite commits.

Before the commit run the commands that reproduce your findings, `rhiza ship generate` and `rhiza render`, folding in their changes. List each with its exit code in a `Checks:` block.

## Outside work

Work outside the intent is a question or a discovered-from draft change set:
`rhiza change new "<intent>" --json`, then
`rhiza link <new> discovered-from <this> --change <new> --json`, then
`rhiza change use <this>`.
