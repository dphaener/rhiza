# Build

A build run takes leaves one at a time: a fresh builder, then a fresh verifier. `RHIZA_CHANGE`, `RHIZA_LEAF` (the leaf's handle) and
`RHIZA_ATTEMPT` (1, 2, ...) are set. Each role ends with `rhiza session done` as its own command after its recording
command succeeded, never chained with `;` or `&&`.

## Start

```
rhiza slice --change "$RHIZA_CHANGE" --json
rhiza show "$RHIZA_LEAF" --json
```

The leaf's `acceptance` is the bar; build only that leaf.

## Builder

Write code and tests to the acceptance. End with exactly one new commit
whose trailers name the leaf, the change set and you:

```
Stream CSV exports

Change: 01M5PWGEF85Y185JZ2ZGCNH40N
Actor: claude-code-1a2b3c4d
Leaf: export-csv
Co-Authored-By: Agent Name <agent@example.com>
```

Keep every trailer, `Co-Authored-By:` included, in the one final paragraph, no blank line.

Stage your untracked actor record in `.rhiza/log` with your work:

```
git add $(grep -l "Record agent actor $RHIZA_ACTOR " .rhiza/log/*/*.json | xargs -n1 dirname | sort -u)
```

Stage only that record, never another actor's.

On a retry (`RHIZA_ATTEMPT` above 1), read the earlier attempts' findings in
`rhiza run status --json` first, then add one fix commit with the same `Leaf:` trailer.
Never rewrite earlier commits.

## Tests

Builders and verifiers run only the tests covering what the change touches, with the
repo's scoped commands, never the repo's full
suite: the PR's CI runs it before merge.

While CI is unstable, or a red is not explained by your change,
do not push or rerun CI to see. Reproduce the
failing test locally, fix it, and push only after it and the
touched packages pass locally.
Reproduce a timing flake by injecting a delay at the suspected race, never by busy loops or parallel stress on the host.

## Command budget

Never sleep and never poll. No command runs longer than about four
minutes; a longer build or test runs once in the background and is
waited for once.

## Outside work

- Outside work the leaf needs: raise a question linked `blocks` to the leaf, and stop.
- Noticed-only work: a discovered-from draft change set, never built in this run.
  `rhiza change new "<intent>" --json`, then
  `rhiza link <new> discovered-from <this> --change <new> --json`, then
  `rhiza change use <this>`.

## Verifier

You work in a detached checkout at the leaf's commit. Re-run the leaf's acceptance as
written, and read the diff. Never edit code. Test runs with fakes only, never real Rex,
agents or `gh`. End with exactly one of:

- Evidence written as yourself, verifying the leaf, then
  `rhiza apply <file> --change "$RHIZA_CHANGE" --json`:

  ```json
  {"writes":[
    {"op":"add","kind":"evidence","handle":"<leaf>-verified",
     "fields":["statement=<what you ran and saw>","result=PASS"]},
    {"op":"link","from":"<leaf>-verified","type":"verifies","to":"<leaf>"}
  ]}
  ```

- Findings, one per failed point, on stdin:

  ```
  rhiza run findings --leaf <leaf> - <<'EOF'
  {"findings":[{"reason":"<what failed and how you ran it>"}]}
  EOF
  ```

Then `rhiza session done`. An attempt with neither Evidence nor findings fails, and Evidence
written by anyone else does not count.
