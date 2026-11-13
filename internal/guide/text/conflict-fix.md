# Conflict fixer

## Task

Fix a conflict that stops a ship run: `RHIZA_CHANGE` and `RHIZA_RUN` are set. You work on
the run's worktree. There are two parts, and the run's prompt tells you which one you
are in:

- A code conflict: origin/main was merged into the branch and the merge is left in
  progress in this worktree.
- A model conflict: another change set landed a write to a node this change set also
  writes, so `rhiza land` refused this one.

Resolve that conflict only.

## Expected outcome

The conflict resolved so both sides' behaviour holds, in one commit trailed `Change:`,
`Actor:` and `Review: <run>`, the merge committed and never left in progress. On a real
clash, nothing is committed and one blocking question names both change sets.

## Read first

For a code conflict:

```
rhiza ship check --json
```

Its `merge_conflict` holds the conflicted files, origin/main's commits with their change
sets, and this branch's change set. Run `rhiza show` on both change sets for their intent.

For a model conflict:

```
rhiza ship check --model-conflict --json
```

It holds the node, the op, the aspects and both change sets' IDs and intents.

## Must do

Code conflict (the merge in progress):

1. Resolve every conflicted file so both sides' behaviour holds.
2. Regenerate golden files and the `rhiza render` output.
3. Run only the touched packages' tests and scripts, with the repo's scoped commands.
4. Run `rhiza ship generate` and `rhiza render`, and fold what they change into the merge.
5. Stage `.rhiza/rendered` and `.rhiza/log` too, so nothing is left untracked:
   `git add -A -- .rhiza/rendered .rhiza/log`
6. Commit the merge with `git commit`, each command's exit code in a `Checks:` block
   above the trailers.

```
Merge origin/main into the export branch

Checks:
- <scoped command for the touched package>: exit 0
- rhiza ship generate: exit 0
- rhiza render: exit 0

Change: 01M5PWGEF85Y185JZ2ZGCNH40N
Actor: claude-code-1a2b3c4d
Review: 01M5PX3N6Y0D2W8VQK4TSE1HMB
Co-Authored-By: Agent Name <agent@example.com>
```

Model conflict (the land was refused):

1. For each stale update of this change set to the conflicting node, move it onto the
   node's landed version with `rhiza set <node> --rebase`, setting every field the
   update writes, so both sides' edits hold.
2. Run `rhiza check`.
3. Stage the log and the rendered files: `git add -A -- .rhiza/rendered .rhiza/log`
4. Commit with the same trailers plus a `Model-Conflict: <other change set>` trailer.

Keep every trailer, `Co-Authored-By:` included, in the one final paragraph, no blank line.

A real clash is one where both sides change the same behaviour or the same statement
incompatibly, so no resolution honours both. Then commit nothing, leave the merge as it
is and do not end your session: ask one blocking question here in the session
(`AskUserQuestion`) that names both change sets, the clashing files or the node, and the
incompatible changes. Wait for the answer, then resolve following it.

<!-- shared: test-scope -->

<!-- shared: command-budget nested -->

## Must not do

<!-- shared: never-list nested -->
- End your session with the merge still in progress.
- Push, rebase or force anything.
- Run the repo's full suite: the PR's CI runs it after the push.
- Ask more than one blocking question.

## Context

<!-- shared: session-done -->

The ship run reaches you after origin/main moved under the branch, or a land was refused.
