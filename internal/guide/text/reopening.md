# Reopening

`rhiza land` is not always the last word. A change set that landed on its branch, but not yet
on main, can be reopened to take more writes, then land again.

## When

The work landed on its branch and needs more: a review asked for a fix, or something was
missed. Reopen it instead of opening a second change set for the same intent.

## The command

```
rhiza reopen [<change>] [--base <rev>] --json
```

- It appends a reopen entry: the change set goes from landed back to accepted.
- It becomes the worktree's current change set (`RHIZA_CHANGE` still wins).
- Write as usual (`rhiza apply`, `rhiza add`, ...), then `rhiza land` again.
- With no `<change>` it takes the current one, else the one named.
- The base is the one `rhiza check` uses: the merge-base with `origin/HEAD`, `origin/main` or
  `main`. `--base <rev>` picks another.

## The refusal

Reopen exits 3 when any land entry of the change set is already in the tree at the base.
Landed work on main is final. Open a new change set and link it to the old one instead:

```
rhiza change new "<intent>" --json
rhiza link <new> discovered-from <old> --change <new> --json
```

It also exits 3 with no git or no base to compare against.

Only landed change sets reopen. Rejected, withdrawn and superseded stay final.

## Never delete or edit log entries

The log is append-only. To fix landed work, never delete or edit a log entry, and never
rewrite a land entry away: reopen the change set, or open a new one as above.
