# Change sets

Every write goes into a change set: a group of operations with an intent, moved through a
lifecycle by entries in the log.

## Kinds

- **intention**: a change to what the product should be. It is reviewed: proposed, accepted,
  then landed. The default for `genesis change new`.
- **record**: writing down what already is (existing code, a past decision). It lands as is,
  with no review: `genesis change new "<intent>" --kind record`, then `genesis land`.

## Lifecycle

```
draft --propose--> proposed --accept--> accepted --land--> landed
                   proposed --reject--> rejected
draft | proposed | accepted --withdraw--> withdrawn
```

- `genesis change new "<intent>" --json` opens one in `draft` and makes it current.
- Writes (`genesis apply`, `genesis add`, ...) go in while it is draft, proposed or accepted.
- `genesis propose`: draft to proposed.
- `genesis accept --rule <rule> --reason "<why>"`: proposed to accepted. Only where policy
  lets you accept; the default rule is `owner`.
- `genesis reject --reason "<why>"`: proposed to rejected.
- `genesis withdraw`: abandon one that has not landed.
- `genesis land`: an accepted intention, or a record in any open state, becomes as-is.
- Landed, rejected, withdrawn and superseded are terminal: nothing more can be written. A
  terminal move on the current change set clears it.

## Which change set

1. `--change <id>` on the command;
2. else `GENESIS_CHANGE` in the environment;
3. else the worktree's current one (`genesis change use <id>`, cleared with
   `genesis change use --none`).

A write with none fails with exit 2 and `next` naming `genesis change new`.

## --record

`--record` on any write opens a record change set, applies the write and lands it in one step.
It ignores the current change set and sees only landed state. Use it to write down what
already exists; never for new intentions.

## Working on one

An agent works on an accepted change set and writes back only revisions within its intent.
Work found outside the intent becomes a new change set linked `discovered-from` the one it was
found in (`genesis link <new> discovered-from <this>`), or a question that `blocks` it. The
work ends with `genesis land`, committed last.

## Building on another

- `genesis link <this> depends-on <that> --change <this>`. A loop is refused (exit 3);
  `check` fails on one in the log.
- Writes and what-if see landed state, every open change set this depends on (directly or
  not), and its own operations.
- `genesis land` exits 3 until every dependency has landed (`next` names it), or on a
  what-if conflict.
- A dependency dropped, or revised over what you wrote: `check` and `ready` name it. Fix by
  writing again, revising, or retiring the link. Rebase is manual.

## Inspecting

- `genesis change show [<id>] --json`: intent, state, operations, the accept or reject.
- `genesis change ls --json`: every change set with kind, state and intent.
- Commits made for a change set carry a `Change: <id>` trailer.
