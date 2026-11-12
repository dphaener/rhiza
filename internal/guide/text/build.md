# Build

Per leaf, a builder then a verifier; `RHIZA_CHANGE`, `RHIZA_LEAF` and `RHIZA_ATTEMPT` are set.
Each role ends with `rhiza session done` as its own command after its recording
command succeeded, never chained with `;` or `&&`.

## Start

```
rhiza slice --change "$RHIZA_CHANGE" --json
rhiza show "$RHIZA_LEAF" --json
```

The leaf's `acceptance` is the bar; build only that leaf.

## Builder

Write code and tests, then run its acceptance
commands, `rhiza ship generate` and `rhiza render`; commit what they change.
End with exactly one new commit, each command's exit code in a
`Checks:` block above the trailers:

```
Stream CSV exports

Checks:
- <acceptance command>: exit 0
- rhiza ship generate: exit 0
- rhiza render: exit 0

Change: 01M5PWGEF85Y185JZ2ZGCNH40N
Actor: claude-code-1a2b3c4d
Leaf: export-csv
Co-Authored-By: Agent Name <agent@example.com>
```

Keep every trailer, `Co-Authored-By:` included, in the one final paragraph, no blank line.

Stage your untracked actor record in `.rhiza/log`, never another actor's:

```
git add $(grep -l "Record agent actor $RHIZA_ACTOR " .rhiza/log/*/*.json | xargs -n1 dirname | sort -u)
```

On a retry (`RHIZA_ATTEMPT` > 1), read findings in
`rhiza run status --json`, then add one fix commit with the same `Leaf:` trailer.
Never rewrite history.

## Tests

Builders and verifiers run only the tests covering what the change touches, with the
repo's scoped commands, never the repo's full
suite: the PR's CI runs it before merge.

While CI is unstable, or a red is not explained by your change,
do not push or rerun CI to see. Reproduce the
failing test locally, fix it, and push only after it and the
touched packages pass locally.
Reproduce flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host.

## Command budget

Never sleep and never poll. No command runs over about four
minutes; longer work runs once in the background and is
waited for once.

## Outside work

- Needed: raise a question linked `blocks` to the leaf; stop.
- Noticed only: a discovered-from draft change set, never built in this run:
  `rhiza change new "<intent>" --json`, then
  `rhiza link <new> discovered-from <this> --change <new> --json`, then `rhiza change use <this>`.

## Verifier

Work in a detached checkout at the leaf's commit. Re-run every command in its
`Checks:` block and every acceptance command; read the diff. A different exit code, or
no `Checks:` block, is a finding naming the command. Never edit code. Tests run with fakes only, never real Rex,
agents or `gh`. End with exactly one of:

- Evidence written as yourself, then
  `rhiza apply <file> --change "$RHIZA_CHANGE" --json`:

  ```json
  {"writes":[{"op":"add","kind":"evidence","handle":"<leaf>-verified",
  "fields":["statement=<what you ran and saw>","result=PASS"]},
  {"op":"link","from":"<leaf>-verified","type":"verifies","to":"<leaf>"}]}
  ```

- Findings, one per failed point, on stdin:

  ```
  rhiza run findings --leaf <leaf> - <<'EOF'
  {"findings":[{"reason":"<what failed and how you ran it>"}]}
  EOF
  ```

An attempt with neither Evidence nor findings fails, and Evidence
written by anyone else does not count.
