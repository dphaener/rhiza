# Build

A build run takes a change set's leaves one at a time. Each leaf gets a fresh builder
session, then a fresh verifier session. `RHIZA_CHANGE`, `RHIZA_LEAF` (the leaf's handle) and
`RHIZA_ATTEMPT` (1, 2, ...) are set. Each role ends with `rhiza session done`.

## Start

```
rhiza slice --change "$RHIZA_CHANGE" --json
rhiza show "$RHIZA_LEAF" --json
```

The leaf's `acceptance` is the bar. Build that leaf only.

## Builder

Write the code and tests to the leaf's acceptance. End with exactly one new commit
whose trailers name the leaf as well as the change set and you:

```
Stream CSV exports

Change: 01M5PWGEF85Y185JZ2ZGCNH40N
Actor: claude-code-1a2b3c4d
Leaf: export-csv
Co-Authored-By: Agent Name <agent@example.com>
```

Keep every trailer, `Co-Authored-By:` included, in the one final paragraph, with no
blank line between them.

Your actor record sits untracked in `.rhiza/log`. Stage it with your work, so the
`Actor:` trailer names a tracked actor:

```
git add $(grep -l "Record agent actor $RHIZA_ACTOR " .rhiza/log/*/*.json | xargs -n1 dirname | sort -u)
```

Stage only that record, never another actor's.

On a retry (`RHIZA_ATTEMPT` above 1), read the earlier attempts' findings in
`rhiza run status --json` first, then add one fix commit with the same `Leaf:` trailer.
Never rewrite the earlier commits.

## Tests

Builders and verifiers run only the tests covering what the change touches, with the
repo's own scoped commands (its contributor guide names them), never the repo's full
suite: the PR's CI runs it before merge.

While CI is unstable, or a red is not explained by your change,
do not push or rerun CI to see. Reproduce the
failing test locally, fix it, and push only after it and the
touched packages pass locally.
Reproduce a timing flake by injecting a delay at the suspected race, never by busy loops or parallel stress on the host.

## Outside work

- Needed outside work, without which the leaf cannot pass: raise a question (the recipe in
  `rhiza guide`) linked `blocks` to the leaf, and stop.
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

- Findings, one per failed point, recorded from stdin with a heredoc, no temp file:

  ```
  rhiza run findings --leaf <leaf> - <<'EOF'
  {"findings":[{"reason":"<what failed and how you ran it>"}]}
  EOF
  ```

Then `rhiza session done`. An attempt with neither Evidence nor findings fails, and Evidence
written by anyone else does not count.
