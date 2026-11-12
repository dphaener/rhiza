# Run agents working on Rhiza

Read this when you run as a builder, verifier, check or fixer on Rhiza
itself. The build and review recipes are in `rhiza guide`; Rhiza's commands
and layout are in [AGENTS.md](../AGENTS.md). This doc holds the traps that
failed earlier attempts.

## Commit trailers

`Change:`, `Actor:`, `Leaf:` (or `Review:`) and `Co-Authored-By:` share the
one last paragraph, with no blank line between them. Git parses only the
last paragraph as trailers; with a blank line before `Co-Authored-By:` the
run sees no `Leaf:` trailer and fails the attempt. Check before
`rhiza session done`:

```bash
git log -1 --format=%B | git interpret-trailers --parse
```

If a commit already has broken trailers, add an `--allow-empty` commit with
the right ones. Never rewrite history.

## A run change set that already landed

A plan run's record change set lands right after the interview, before
review. Any run step that writes to the model (a question, an interview
source, an arbiter's answer) must handle that: write into the open change
set that holds the node, else a record write. A write in one unlanded change
set cannot refer to a node in another unlanded one. Test the landed path
(see the `*-landed.txtar` testscripts).

## Findings HEAD already fixes

A fixer that makes no commit fails with "no new commit". If the findings
sent to you are already fixed at HEAD (compare the quoted code with HEAD, and
`git log`), do not invent a change: make one empty commit with the leaf
trailers, naming the commit that fixed it, and run `rhiza session done`:

```bash
git commit --allow-empty
```

## A PR that conflicts with main

Merge `origin/main` into the branch and push normally. Never rebase and
force-push: it is denied, and the merge leaves ship's own rebase step
nothing to move. A conflicting PR gets no CI, so check
`git merge-tree --write-tree origin/main HEAD` before waiting on it.

## A hand push to a live run's PR

A run's `wait-pr` step passes only when the PR head equals `pushed` in
`.git/rhiza/ship/<cs>.json`. After pushing by hand, set `pushed` there to
the new HEAD, or the run stalls until its timeout. Where you can, use
`rhiza ship rebase` with `RHIZA_RUN` and `RHIZA_CHANGE` set: it records its
own push.
